package api

import (
	"client-app/internal/model/entity"
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"
	"client-app/internal/service"
	"client-app/utility/encrypt"
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// sTenant 租户业务逻辑实现
type sTenant struct{}

// NewTenant 创建租户服务实例
func NewTenant() *sTenant {
	return &sTenant{}
}

// 在init函数中注册服务
func init() {
	service.RegisterTenant(NewTenant())
}

// GetTenantList 获取租户列表
func (s *sTenant) GetTenantList(ctx context.Context, in *sysin.TenantListInp) (*sysout.TenantListModel, error) {
	// 构建查询条件
	m := g.DB().Model("sys_tenants").Where("deleted_at IS NULL")

	// 按条件筛选
	if in.Name != "" {
		m = m.WhereLike("name", "%"+in.Name+"%")
	}
	if in.Code != "" {
		m = m.WhereLike("code", "%"+in.Code+"%")
	}
	if in.Domain != "" {
		m = m.WhereLike("domain", "%"+in.Domain+"%")
	}
	if in.Status > 0 {
		m = m.Where("status", in.Status)
	}

	// 获取总数
	totalCount, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取租户总数失败")
	}

	// 分页查询
	m = m.Order("id DESC").
		Limit(in.PageSize).
		Offset((in.Page - 1) * in.PageSize)

	var list []*entity.Tenant
	if err := m.Scan(&list); err != nil {
		return nil, gerror.Wrap(err, "查询租户列表失败")
	}

	// 获取管理员用户名
	var adminUserIds []uint64
	for _, tenant := range list {
		if tenant.AdminUserId > 0 {
			adminUserIds = append(adminUserIds, tenant.AdminUserId)
		}
	}

	var adminUsers []*entity.User
	if len(adminUserIds) > 0 {
		err = g.DB().Model("sys_users").WhereIn("id", adminUserIds).Scan(&adminUsers)
		if err != nil {
			g.Log().Warningf(ctx, "查询管理员用户失败: %v", err)
		}
	}

	// 构建管理员用户映射
	adminUserMap := make(map[uint64]string)
	for _, user := range adminUsers {
		adminUserMap[uint64(user.Id)] = user.Username
	}

	// 转换为输出模型
	var tenantList []*sysout.TenantModel
	for _, tenant := range list {
		tenantModel := &sysout.TenantModel{
			Id:           tenant.Id,
			Name:         tenant.Name,
			Code:         tenant.Code,
			Domain:       tenant.Domain,
			Status:       tenant.Status,
			StatusName:   sysout.GetTenantStatusName(tenant.Status),
			MaxUsers:     tenant.MaxUsers,
			StorageLimit: tenant.StorageLimit,
			ExpireAt:     tenant.ExpireAt,
			AdminUserId:  tenant.AdminUserId,
			AdminName:    adminUserMap[tenant.AdminUserId],
			Remark:       tenant.Remark,
			CreatedAt:    tenant.CreatedAt,
			UpdatedAt:    tenant.UpdatedAt,
		}

		// 解析配置
		if tenant.Config != "" {
			var config interface{}
			if err := json.Unmarshal([]byte(tenant.Config), &config); err == nil {
				tenantModel.Config = config
			}
		}

		tenantList = append(tenantList, tenantModel)
	}

	return &sysout.TenantListModel{
		List:     tenantList,
		Total:    int64(totalCount),
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}

// CreateTenant 创建租户
func (s *sTenant) CreateTenant(ctx context.Context, in *sysin.CreateTenantInp) (*sysout.TenantModel, error) {
	// 验证租户编码唯一性
	count, err := g.DB().Model("sys_tenants").Where("code", in.Code).Where("deleted_at IS NULL").Count()
	if err != nil {
		return nil, gerror.Wrap(err, "验证租户编码失败")
	}
	if count > 0 {
		return nil, gerror.New("租户编码已存在")
	}

	// 验证管理员用户名唯一性（全局唯一）
	userCount, err := g.DB().Model("sys_users").Where("username", in.AdminName).Where("deleted_at IS NULL").Count()
	if err != nil {
		return nil, gerror.Wrap(err, "验证管理员用户名失败")
	}
	if userCount > 0 {
		return nil, gerror.New("管理员用户名已存在")
	}

	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 创建租户记录
		tenantData := g.Map{
			"name":          in.Name,
			"code":          in.Code,
			"domain":        in.Domain,
			"status":        1, // 默认正常状态
			"max_users":     in.MaxUsers,
			"storage_limit": in.StorageLimit,
			"expire_at":     in.ExpireAt,
			"remark":        in.Remark,
			"created_by":    gconv.Uint64(ctx.Value("userId")),
			"updated_by":    gconv.Uint64(ctx.Value("userId")),
			"created_at":    gtime.Now(),
			"updated_at":    gtime.Now(),
		}

		tenantResult, err := tx.Model("sys_tenants").Data(tenantData).Insert()
		if err != nil {
			return gerror.Wrap(err, "创建租户失败")
		}

		tenantId, err := tenantResult.LastInsertId()
		if err != nil {
			return gerror.Wrap(err, "获取租户ID失败")
		}

		// 2. 创建管理员用户
		salt := encrypt.GenerateSalt()
		hashedPassword := encrypt.HashPassword(in.AdminPassword, salt)

		adminUserData := g.Map{
			"tenant_id":  tenantId,
			"username":   in.AdminName,
			"password":   hashedPassword,
			"salt":       salt,
			"email":      in.AdminEmail,
			"real_name":  "租户管理员",
			"nickname":   "管理员",
			"status":     1,
			"created_by": 0,
			"updated_by": 0,
			"created_at": gtime.Now(),
			"updated_at": gtime.Now(),
		}

		adminResult, err := tx.Model("sys_users").Data(adminUserData).Insert()
		if err != nil {
			return gerror.Wrap(err, "创建管理员用户失败")
		}

		adminUserId, err := adminResult.LastInsertId()
		if err != nil {
			return gerror.Wrap(err, "获取管理员用户ID失败")
		}

		// 3. 更新租户的管理员用户ID
		_, err = tx.Model("sys_tenants").Where("id", tenantId).Data(g.Map{
			"admin_user_id": adminUserId,
			"updated_at":    gtime.Now(),
		}).Update()
		if err != nil {
			return gerror.Wrap(err, "更新租户管理员ID失败")
		}

		// 4. 创建租户管理员角色
		roleData := g.Map{
			"tenant_id":  tenantId,
			"name":       "租户管理员",
			"code":       "tenant_admin",
			"data_scope": 1, // 全部数据权限
			"status":     1,
			"sort":       1,
			"remark":     "租户管理员角色，拥有租户内所有权限",
			"created_by": adminUserId,
			"updated_by": adminUserId,
			"created_at": gtime.Now(),
			"updated_at": gtime.Now(),
		}

		roleResult, err := tx.Model("sys_roles").Data(roleData).Insert()
		if err != nil {
			return gerror.Wrap(err, "创建租户管理员角色失败")
		}

		roleId, err := roleResult.LastInsertId()
		if err != nil {
			return gerror.Wrap(err, "获取角色ID失败")
		}

		// 5. 分配角色给管理员用户
		userRoleData := g.Map{
			"tenant_id":  tenantId,
			"user_id":    adminUserId,
			"role_id":    roleId,
			"is_primary": 1,
			"created_at": gtime.Now(),
			"updated_at": gtime.Now(),
		}

		_, err = tx.Model("sys_user_roles").Data(userRoleData).Insert()
		if err != nil {
			return gerror.Wrap(err, "分配角色失败")
		}

		// 6. 为租户管理员角色分配默认菜单权限
		err = s.assignDefaultMenusToRole(ctx, tx, roleId, tenantId)
		if err != nil {
			return gerror.Wrap(err, "分配默认菜单权限失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 查询创建的租户信息并返回
	var tenant *entity.Tenant
	err = g.DB().Model("sys_tenants").Where("code", in.Code).Scan(&tenant)
	if err != nil {
		return nil, gerror.Wrap(err, "查询租户信息失败")
	}

	return &sysout.TenantModel{
		Id:           tenant.Id,
		Name:         tenant.Name,
		Code:         tenant.Code,
		Domain:       tenant.Domain,
		Status:       tenant.Status,
		StatusName:   sysout.GetTenantStatusName(tenant.Status),
		MaxUsers:     tenant.MaxUsers,
		StorageLimit: tenant.StorageLimit,
		ExpireAt:     tenant.ExpireAt,
		AdminUserId:  tenant.AdminUserId,
		AdminName:    in.AdminName,
		Remark:       tenant.Remark,
		CreatedAt:    tenant.CreatedAt,
		UpdatedAt:    tenant.UpdatedAt,
	}, nil
}

// assignDefaultMenusToRole 为角色分配默认菜单权限
func (s *sTenant) assignDefaultMenusToRole(ctx context.Context, tx gdb.TX, roleId int64, tenantId int64) error {
	// 获取系统默认菜单（租户管理员应该拥有的菜单）
	var defaultMenus []*entity.Menu
	err := tx.Model("sys_menus").Where("status = ? AND deleted_at IS NULL", entity.MenuStatusNormal).
		Where("type IN (?)", []int{entity.MenuTypeDir, entity.MenuTypeMenu}).
		Order("sort ASC, id ASC").Scan(&defaultMenus)
	if err != nil {
		return err
	}

	// 为角色分配菜单权限
	for _, menu := range defaultMenus {
		roleMenuData := g.Map{
			"tenant_id":  tenantId,
			"role_id":    roleId,
			"menu_id":    menu.Id,
			"created_at": gtime.Now(),
			"updated_at": gtime.Now(),
		}

		_, err = tx.Model("sys_role_menus").Data(roleMenuData).Insert()
		if err != nil {
			// 如果重复插入，忽略错误
			g.Log().Warningf(ctx, "菜单权限已存在，跳过: role_id=%d, menu_id=%d", roleId, menu.Id)
		}
	}

	return nil
}

// UpdateTenant 更新租户
func (s *sTenant) UpdateTenant(ctx context.Context, in *sysin.UpdateTenantInp) (*sysout.TenantModel, error) {
	// 检查租户是否存在
	var tenant *entity.Tenant
	err := g.DB().Model("sys_tenants").Where("id", in.Id).Where("deleted_at IS NULL").Scan(&tenant)
	if err != nil {
		return nil, gerror.Wrap(err, "查询租户失败")
	}
	if tenant == nil {
		return nil, gerror.New("租户不存在")
	}

	// 检查是否为系统租户
	if tenant.Code == "system" {
		return nil, gerror.New("系统租户不能修改")
	}

	// 更新租户信息
	updateData := g.Map{
		"name":          in.Name,
		"domain":        in.Domain,
		"max_users":     in.MaxUsers,
		"storage_limit": in.StorageLimit,
		"expire_at":     in.ExpireAt,
		"remark":        in.Remark,
		"updated_by":    gconv.Uint64(ctx.Value("userId")),
		"updated_at":    gtime.Now(),
	}

	_, err = g.DB().Model("sys_tenants").Where("id", in.Id).Data(updateData).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新租户失败")
	}

	// 查询更新后的租户信息
	err = g.DB().Model("sys_tenants").Where("id", in.Id).Scan(&tenant)
	if err != nil {
		return nil, gerror.Wrap(err, "查询更新后的租户信息失败")
	}

	return &sysout.TenantModel{
		Id:           tenant.Id,
		Name:         tenant.Name,
		Code:         tenant.Code,
		Domain:       tenant.Domain,
		Status:       tenant.Status,
		StatusName:   sysout.GetTenantStatusName(tenant.Status),
		MaxUsers:     tenant.MaxUsers,
		StorageLimit: tenant.StorageLimit,
		ExpireAt:     tenant.ExpireAt,
		AdminUserId:  tenant.AdminUserId,
		Remark:       tenant.Remark,
		CreatedAt:    tenant.CreatedAt,
		UpdatedAt:    tenant.UpdatedAt,
	}, nil
}

// DeleteTenant 删除租户
func (s *sTenant) DeleteTenant(ctx context.Context, in *sysin.DeleteTenantInp) error {
	// 检查租户是否存在
	var tenant *entity.Tenant
	err := g.DB().Model("sys_tenants").Where("id", in.Id).Where("deleted_at IS NULL").Scan(&tenant)
	if err != nil {
		return gerror.Wrap(err, "查询租户失败")
	}
	if tenant == nil {
		return gerror.New("租户不存在")
	}

	// 检查是否为系统租户
	if tenant.Code == "system" {
		return gerror.New("系统租户不能删除")
	}

	// 开启事务删除
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		now := gtime.Now()

		// 软删除租户
		_, err := tx.Model("sys_tenants").Where("id", in.Id).Data(g.Map{
			"deleted_at": now,
			"updated_by": gconv.Uint64(ctx.Value("userId")),
			"updated_at": now,
		}).Update()
		if err != nil {
			return gerror.Wrap(err, "删除租户失败")
		}

		// 软删除租户下的所有用户
		_, err = tx.Model("sys_users").Where("tenant_id", in.Id).Data(g.Map{
			"deleted_at": now,
			"updated_by": gconv.Uint64(ctx.Value("userId")),
			"updated_at": now,
		}).Update()
		if err != nil {
			return gerror.Wrap(err, "删除租户用户失败")
		}

		// 软删除租户下的所有角色
		_, err = tx.Model("sys_roles").Where("tenant_id", in.Id).Data(g.Map{
			"deleted_at": now,
			"updated_by": gconv.Uint64(ctx.Value("userId")),
			"updated_at": now,
		}).Update()
		if err != nil {
			return gerror.Wrap(err, "删除租户角色失败")
		}

		// 删除租户的角色菜单关联（菜单本身不删除，因为是全局共享的）
		_, err = tx.Model("sys_role_menus").Where("tenant_id", in.Id).Delete()
		if err != nil {
			return gerror.Wrap(err, "删除租户角色菜单关联失败")
		}

		return nil
	})

	return err
}

// GetTenantDetail 获取租户详情
func (s *sTenant) GetTenantDetail(ctx context.Context, in *sysin.TenantDetailInp) (*sysout.TenantDetailModel, error) {
	// 查询租户基本信息
	var tenant *entity.Tenant
	err := g.DB().Model("sys_tenants").Where("id", in.Id).Where("deleted_at IS NULL").Scan(&tenant)
	if err != nil {
		return nil, gerror.Wrap(err, "查询租户失败")
	}
	if tenant == nil {
		return nil, gerror.New("租户不存在")
	}

	// 查询管理员用户名
	var adminName string
	if tenant.AdminUserId > 0 {
		val, err := g.DB().Model("sys_users").Where("id", tenant.AdminUserId).Value("username")
		if err != nil {
			g.Log().Warningf(ctx, "查询管理员用户名失败: %v", err)
		} else if val != nil {
			adminName = val.String()
		}
	}

	tenantModel := &sysout.TenantModel{
		Id:           tenant.Id,
		Name:         tenant.Name,
		Code:         tenant.Code,
		Domain:       tenant.Domain,
		Status:       tenant.Status,
		StatusName:   sysout.GetTenantStatusName(tenant.Status),
		MaxUsers:     tenant.MaxUsers,
		StorageLimit: tenant.StorageLimit,
		ExpireAt:     tenant.ExpireAt,
		AdminUserId:  tenant.AdminUserId,
		AdminName:    adminName,
		Remark:       tenant.Remark,
		CreatedAt:    tenant.CreatedAt,
		UpdatedAt:    tenant.UpdatedAt,
	}

	// 解析配置
	if tenant.Config != "" {
		var config interface{}
		if err := json.Unmarshal([]byte(tenant.Config), &config); err == nil {
			tenantModel.Config = config
		}
	}

	// 获取统计信息
	stats, err := s.GetTenantStats(ctx, &sysin.TenantStatsInp{Id: in.Id})
	if err != nil {
		g.Log().Warningf(ctx, "获取租户统计信息失败: %v", err)
		stats = &sysout.TenantStatsModel{}
	}

	return &sysout.TenantDetailModel{
		TenantModel: tenantModel,
		Stats:       stats,
	}, nil
}

// UpdateTenantStatus 更新租户状态
func (s *sTenant) UpdateTenantStatus(ctx context.Context, in *sysin.TenantStatusInp) error {
	// 检查租户是否存在
	var tenant *entity.Tenant
	err := g.DB().Model("sys_tenants").Where("id", in.Id).Where("deleted_at IS NULL").Scan(&tenant)
	if err != nil {
		return gerror.Wrap(err, "查询租户失败")
	}
	if tenant == nil {
		return gerror.New("租户不存在")
	}

	// 检查是否为系统租户
	if tenant.Code == "system" {
		return gerror.New("系统租户状态不能修改")
	}

	// 更新状态
	_, err = g.DB().Model("sys_tenants").Where("id", in.Id).Data(g.Map{
		"status":     in.Status,
		"updated_by": gconv.Uint64(ctx.Value("userId")),
		"updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return gerror.Wrap(err, "更新租户状态失败")
	}

	return nil
}

// GetTenantStats 获取租户统计信息
func (s *sTenant) GetTenantStats(ctx context.Context, in *sysin.TenantStatsInp) (*sysout.TenantStatsModel, error) {
	// 用户数量
	userCount, err := g.DB().Model("sys_users").Where("tenant_id", in.Id).Where("deleted_at IS NULL").Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计用户数量失败")
	}

	// 角色数量
	roleCount, err := g.DB().Model("sys_roles").Where("tenant_id", in.Id).Where("deleted_at IS NULL").Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计角色数量失败")
	}

	// 菜单数量（租户可访问的菜单数量，通过角色菜单关联统计）
	menuCount, err := g.DB().Model("sys_role_menus rm").
		LeftJoin("sys_menus m", "rm.menu_id = m.id").
		Where("rm.tenant_id = ? AND m.deleted_at IS NULL", in.Id).
		Count("DISTINCT rm.menu_id")
	if err != nil {
		return nil, gerror.Wrap(err, "统计菜单数量失败")
	}

	// 最后活跃时间（最近登录的用户）
	var lastActiveTime *gtime.Time
	val, err := g.DB().Model("sys_users").Where("tenant_id", in.Id).Where("deleted_at IS NULL").
		Where("login_at IS NOT NULL").Order("login_at DESC").Limit(1).Value("login_at")
	if err != nil {
		g.Log().Warningf(ctx, "查询最后活跃时间失败: %v", err)
	} else if val != nil {
		lastActiveTime = val.GTime()
	}

	// 存储使用量（这里是示例，实际需要根据业务逻辑计算）
	var storageUsed int64 = 0

	return &sysout.TenantStatsModel{
		UserCount:       userCount,
		RoleCount:       roleCount,
		MenuCount:       menuCount,
		StorageUsed:     storageUsed,
		StorageUsedText: sysout.FormatStorageSize(storageUsed),
		LastActiveTime:  lastActiveTime,
	}, nil
}

// UpdateTenantConfig 更新租户配置
func (s *sTenant) UpdateTenantConfig(ctx context.Context, in *sysin.TenantConfigInp) error {
	// 检查租户是否存在
	count, err := g.DB().Model("sys_tenants").Where("id", in.Id).Where("deleted_at IS NULL").Count()
	if err != nil {
		return gerror.Wrap(err, "查询租户失败")
	}
	if count == 0 {
		return gerror.New("租户不存在")
	}

	// 序列化配置
	configJson, err := json.Marshal(in.Config)
	if err != nil {
		return gerror.Wrap(err, "序列化配置失败")
	}

	// 更新配置
	_, err = g.DB().Model("sys_tenants").Where("id", in.Id).Data(g.Map{
		"config":     string(configJson),
		"updated_by": gconv.Uint64(ctx.Value("userId")),
		"updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return gerror.Wrap(err, "更新租户配置失败")
	}

	return nil
}

// GetTenantOptions 获取租户选项列表
func (s *sTenant) GetTenantOptions(ctx context.Context) (*sysout.TenantOptionsModel, error) {
	var tenants []*entity.Tenant
	err := g.DB().Model("sys_tenants").Where("status", 1).Where("deleted_at IS NULL").
		Order("sort ASC, id ASC").Scan(&tenants)
	if err != nil {
		return nil, gerror.Wrap(err, "查询租户选项失败")
	}

	var options []*sysout.TenantOptionModel
	for _, tenant := range tenants {
		options = append(options, &sysout.TenantOptionModel{
			Value: tenant.Id,
			Label: tenant.Name,
			Code:  tenant.Code,
		})
	}

	return &sysout.TenantOptionsModel{
		List: options,
	}, nil
}

// GetTenantByCode 根据编码获取租户信息
func (s *sTenant) GetTenantByCode(ctx context.Context, code string) (*sysout.TenantModel, error) {
	var tenant *entity.Tenant
	err := g.DB().Model("sys_tenants").Where("code", code).Where("deleted_at IS NULL").Scan(&tenant)
	if err != nil {
		return nil, gerror.Wrap(err, "查询租户失败")
	}
	if tenant == nil {
		return nil, gerror.New("租户不存在")
	}

	return &sysout.TenantModel{
		Id:           tenant.Id,
		Name:         tenant.Name,
		Code:         tenant.Code,
		Domain:       tenant.Domain,
		Status:       tenant.Status,
		StatusName:   sysout.GetTenantStatusName(tenant.Status),
		MaxUsers:     tenant.MaxUsers,
		StorageLimit: tenant.StorageLimit,
		ExpireAt:     tenant.ExpireAt,
		AdminUserId:  tenant.AdminUserId,
		Remark:       tenant.Remark,
		CreatedAt:    tenant.CreatedAt,
		UpdatedAt:    tenant.UpdatedAt,
	}, nil
}

// ValidateTenantAccess 验证租户访问权限
func (s *sTenant) ValidateTenantAccess(ctx context.Context, tenantId uint64) error {
	// 检查租户是否存在且正常
	var tenant *entity.Tenant
	err := g.DB().Model("sys_tenants").Where("id", tenantId).Where("deleted_at IS NULL").Scan(&tenant)
	if err != nil {
		return gerror.Wrap(err, "查询租户失败")
	}
	if tenant == nil {
		return gerror.New("租户不存在")
	}

	// 检查租户状态
	if !tenant.IsNormal() {
		return gerror.New("租户已被禁用或锁定")
	}

	// 检查租户是否过期
	if tenant.IsExpired() {
		return gerror.New("租户已过期")
	}

	return nil
}
