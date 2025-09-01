package api

import (
	"client-app/internal/consts"
	"client-app/internal/model/entity"
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"
	"client-app/utility/simple"
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// ValidateUserWithTenant 验证用户密码（支持多租户）
func (s *sUser) ValidateUserWithTenant(ctx context.Context, username, password string, tenantId uint64) (user *sysout.UserModel, err error) {
	// 获取用户信息（加入租户过滤条件）
	var userEntity *entity.User
	err = g.DB().Model("sys_users").
		Where("username = ? AND tenant_id = ? AND deleted_at IS NULL", username, tenantId).
		Scan(&userEntity)
	if err != nil {
		return nil, gerror.Newf("查询用户失败: %v", err)
	}

	if userEntity == nil {
		return nil, gerror.New("用户名或密码错误")
	}

	// 验证密码
	if err := simple.CheckPassword(password, userEntity.Salt, userEntity.Password); err != nil {
		return nil, gerror.New("用户名或密码错误")
	}

	// 验证租户状态
	var tenant *entity.Tenant
	err = g.DB().Model("sys_tenants").Where("id = ? AND deleted_at IS NULL", tenantId).Scan(&tenant)
	if err != nil {
		return nil, gerror.Newf("查询租户失败: %v", err)
	}
	if tenant == nil {
		return nil, gerror.New("租户不存在")
	}
	if !tenant.IsNormal() {
		return nil, gerror.New("租户已被禁用或锁定")
	}
	if tenant.IsExpired() {
		return nil, gerror.New("租户已过期")
	}

	return sysout.ConvertToUserModel(userEntity), nil
}

// LoginWithTenant 用户登录（支持多租户）
func (s *sUser) LoginWithTenant(ctx context.Context, in *sysin.UserLoginInp) (res *sysout.LoginTokenModel, err error) {
	// 验证验证码
	if err = s.VerifyCaptcha(ctx, in.CaptchaId, in.Captcha); err != nil {
		return nil, err
	}

	// 根据租户编码获取租户信息
	tenant, err := s.getTenantByCode(ctx, in.TenantCode)
	if err != nil {
		return nil, gerror.Newf("租户验证失败: %v", err)
	}

	// 验证用户密码（包含租户信息）
	user, err := s.ValidateUserWithTenant(ctx, in.Username, in.Password, tenant.Id)
	if err != nil {
		return nil, err
	}

	// 检查用户状态
	if user.Status != entity.UserStatusNormal {
		switch user.Status {
		case entity.UserStatusLocked:
			return nil, gerror.New(consts.GetAuthErrorMessage(consts.ErrUserLocked))
		case entity.UserStatusDisabled:
			return nil, gerror.New(consts.GetAuthErrorMessage(consts.ErrUserDisabled))
		default:
			return nil, gerror.New("用户状态异常，无法登录")
		}
	}

	// 获取用户角色信息（租户过滤）
	userRole, err := s.getUserPrimaryRoleWithTenant(ctx, user.Id, tenant.Id)
	if err != nil {
		return nil, gerror.Newf("获取用户角色失败: %v", err)
	}

	// 生成JWT Token（包含租户信息）
	payload := &simple.JWTPayload{
		UserId:     user.Id,
		TenantId:   int64(tenant.Id),
		TenantCode: tenant.Code,
		Username:   user.Username,
		RoleId:     userRole.RoleId,
		RoleKey:    userRole.RoleCode,
		DeptId:     user.DeptId,
		App:        consts.AppApi,
	}

	secretKey := simple.GetJWTSecretKey(ctx)
	accessToken, err := simple.GenerateJWTToken(payload, secretKey)
	if err != nil {
		return nil, gerror.Newf("生成访问令牌失败: %v", err)
	}

	// 生成刷新令牌
	refreshToken, err := s.generateRefreshToken(ctx, user.Id)
	if err != nil {
		return nil, gerror.Newf("生成刷新令牌失败: %v", err)
	}

	// 更新用户登录信息
	if err = s.updateLoginInfo(ctx, user.Id); err != nil {
		g.Log().Warningf(ctx, "更新用户登录信息失败: %v", err)
	}

	// 获取用户权限（租户过滤）
	permissions, menuIds, err := s.getUserPermissionsWithTenant(ctx, user.Id, tenant.Id)
	if err != nil {
		g.Log().Warningf(ctx, "获取用户权限失败: %v", err)
		permissions = []string{}
		menuIds = []int64{}
	}

	// 构建响应
	res = &sysout.LoginTokenModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    24 * 3600, // 24小时
		UserInfo:     user,
		Permissions:  permissions,
		MenuIds:      menuIds,
	}

	return res, nil
}

// getUserPrimaryRoleWithTenant 获取用户主要角色（租户过滤）
func (s *sUser) getUserPrimaryRoleWithTenant(ctx context.Context, userId int64, tenantId uint64) (*UserRoleWithCode, error) {
	var userRoleWithCode *UserRoleWithCode
	err := g.DB().Model("sys_user_roles ur").
		LeftJoin("sys_roles r", "ur.role_id = r.id").
		Where("ur.user_id = ? AND ur.tenant_id = ? AND ur.is_primary = ? AND r.deleted_at IS NULL",
			userId, tenantId, entity.IsPrimaryRole).
		Fields("ur.*, r.role_code").
		Scan(&userRoleWithCode)
	if err != nil {
		return nil, err
	}

	if userRoleWithCode == nil {
		return nil, gerror.New("用户没有分配主要角色")
	}

	return userRoleWithCode, nil
}

// getUserPermissionsWithTenant 获取用户权限（租户过滤）
func (s *sUser) getUserPermissionsWithTenant(ctx context.Context, userId int64, tenantId uint64) (permissions []string, menuIds []int64, err error) {
	// 获取用户的所有角色
	var roleIds []int64
	roleIdsVal, err := g.DB().Model("sys_user_roles ur").
		LeftJoin("sys_roles r", "ur.role_id = r.id").
		Where("ur.user_id = ? AND ur.tenant_id = ? AND r.status = ? AND r.deleted_at IS NULL",
			userId, tenantId, entity.RoleStatusNormal).
		Array("ur.role_id")
	if err != nil {
		return nil, nil, err
	}

	for _, val := range roleIdsVal {
		roleIds = append(roleIds, val.Int64())
	}

	if len(roleIds) == 0 {
		return []string{}, []int64{}, nil
	}

	// 获取角色关联的菜单权限
	var menus []*entity.Menu
	err = g.DB().Model("sys_role_menus rm").
		LeftJoin("sys_menus m", "rm.menu_id = m.id").
		Where("rm.role_id IN (?) AND rm.tenant_id = ? AND m.status = ? AND m.deleted_at IS NULL",
			roleIds, tenantId, entity.MenuStatusNormal).
		Fields("m.*").
		Scan(&menus)
	if err != nil {
		return nil, nil, err
	}

	// 提取权限标识和菜单ID
	permissionSet := make(map[string]bool)
	menuIdSet := make(map[int64]bool)

	for _, menu := range menus {
		if menu.Permission != "" {
			permissionSet[menu.Permission] = true
		}
		menuIdSet[menu.Id] = true
	}

	// 转换为切片
	permissions = make([]string, 0, len(permissionSet))
	for permission := range permissionSet {
		permissions = append(permissions, permission)
	}

	menuIds = make([]int64, 0, len(menuIdSet))
	for menuId := range menuIdSet {
		menuIds = append(menuIds, menuId)
	}

	return permissions, menuIds, nil
}

// GetUserByUsernameWithTenant 根据用户名获取用户（租户过滤）
func (s *sUser) GetUserByUsernameWithTenant(ctx context.Context, username string, tenantId uint64) (user *sysout.UserModel, err error) {
	var entity *entity.User
	err = g.DB().Model("sys_users").
		Where("username = ? AND tenant_id = ? AND deleted_at IS NULL", username, tenantId).
		Scan(&entity)
	if err != nil {
		return nil, gerror.Newf("查询用户失败: %v", err)
	}

	if entity == nil {
		return nil, gerror.New("用户不存在")
	}

	return sysout.ConvertToUserModel(entity), nil
}

// getTenantByCode 根据租户编码获取租户信息
func (s *sUser) getTenantByCode(ctx context.Context, tenantCode string) (*entity.Tenant, error) {
	var tenant *entity.Tenant
	err := g.DB().Model("sys_tenants").Where("code = ? AND deleted_at IS NULL", tenantCode).Scan(&tenant)
	if err != nil {
		return nil, gerror.Wrap(err, "查询租户失败")
	}
	if tenant == nil {
		return nil, gerror.New("租户不存在")
	}

	// 检查租户状态
	if !tenant.IsNormal() {
		return nil, gerror.New("租户已被禁用或锁定")
	}

	// 检查租户是否过期
	if tenant.IsExpired() {
		return nil, gerror.New("租户已过期")
	}

	return tenant, nil
}
