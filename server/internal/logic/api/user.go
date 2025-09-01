package api

import (
	"client-app/internal/consts"
	"client-app/internal/model/entity"
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"
	"client-app/internal/service"
	"client-app/utility/captcha"
	"client-app/utility/simple"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sUser struct{}

func init() {
	service.RegisterUser(NewUser())
}

func NewUser() *sUser {
	return &sUser{}
}

// Login 用户登录（多租户版本）
func (s *sUser) Login(ctx context.Context, in *sysin.UserLoginInp) (res *sysout.LoginTokenModel, err error) {
	// 调用多租户登录方法
	return s.LoginWithTenant(ctx, in)
}

// Logout 用户退出
func (s *sUser) Logout(ctx context.Context, in *sysin.UserLogoutInp) error {
	// 从token中解析用户信息
	secretKey := simple.GetJWTSecretKey(ctx)
	payload, err := simple.ParseJWTToken(in.Token, secretKey)
	if err != nil {
		// 即使token无效也继续执行，确保退出操作成功
		g.Log().Warningf(ctx, "解析退出token失败: %v", err)
	}

	// 将token加入黑名单（可选，这里简化处理）
	if payload != nil {
		cacheKey := fmt.Sprintf("blacklist_token_%s", in.Token)
		gcache.Set(ctx, cacheKey, 1, time.Duration(payload.Exp-time.Now().Unix())*time.Second)
	}

	return nil
}

// GetProfile 获取用户资料
func (s *sUser) GetProfile(ctx context.Context, userId int64) (res *sysout.UserModel, err error) {
	var user *entity.User
	err = g.DB().Model("sys_users").Where("id = ? AND deleted_at IS NULL", userId).Scan(&user)
	if err != nil {
		return nil, gerror.Newf("查询用户信息失败: %v", err)
	}

	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	// 转换为输出模型
	res = sysout.ConvertToUserModel(user)
	return res, nil
}

// RefreshToken 刷新访问令牌
func (s *sUser) RefreshToken(ctx context.Context, refreshToken string) (res *service.TokenInfo, err error) {
	// 验证刷新令牌
	userId, err := s.validateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	var user *entity.User
	err = g.DB().Model("sys_users").Where("id = ? AND deleted_at IS NULL", userId).Scan(&user)
	if err != nil {
		return nil, gerror.Newf("查询用户信息失败: %v", err)
	}

	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	// 检查用户状态
	if user.Status != entity.UserStatusNormal {
		return nil, gerror.New("用户状态异常，无法刷新令牌")
	}

	// 获取用户角色信息
	userRole, err := s.getUserPrimaryRole(ctx, user.Id)
	if err != nil {
		return nil, gerror.Newf("获取用户角色失败: %v", err)
	}

	// 生成新的访问令牌
	payload := &simple.JWTPayload{
		UserId:   user.Id,
		Username: user.Username,
		RoleId:   userRole.RoleId,
		RoleKey:  userRole.RoleCode,
		DeptId:   user.DeptId,
		App:      consts.AppApi,
	}

	secretKey := simple.GetJWTSecretKey(ctx)
	newAccessToken, err := simple.GenerateJWTToken(payload, secretKey)
	if err != nil {
		return nil, gerror.Newf("生成访问令牌失败: %v", err)
	}

	// 生成新的刷新令牌
	newRefreshToken, err := s.generateRefreshToken(ctx, user.Id)
	if err != nil {
		return nil, gerror.Newf("生成刷新令牌失败: %v", err)
	}

	// 删除旧的刷新令牌
	cacheKey := fmt.Sprintf("refresh_token_%s", refreshToken)
	gcache.Remove(ctx, cacheKey)

	res = &service.TokenInfo{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    24 * 3600, // 24小时
	}

	return res, nil
}

// ChangePassword 修改密码
func (s *sUser) ChangePassword(ctx context.Context, userId int64, oldPassword, newPassword string) error {
	// 获取用户信息
	var user *entity.User
	err := g.DB().Model("sys_users").Where("id = ?", userId).Scan(&user)
	if err != nil {
		return gerror.Newf("查询用户信息失败: %v", err)
	}

	if user == nil {
		return gerror.New("用户不存在")
	}

	// 验证原密码
	if err := simple.CheckPassword(oldPassword, user.Salt, user.Password); err != nil {
		return gerror.New("原密码错误")
	}

	// 生成新密码hash
	newSalt := gmd5.MustEncryptString(gconv.String(time.Now().UnixNano()))
	decryptedPassword, err := simple.DecryptText(newPassword)
	if err != nil {
		return gerror.Newf("密码解密失败: %v", err)
	}
	newPasswordHash := gmd5.MustEncryptString(decryptedPassword + newSalt)

	// 更新密码
	_, err = g.DB().Model("sys_users").Where("id = ?", userId).Update(g.Map{
		"password":   newPasswordHash,
		"salt":       newSalt,
		"updated_at": gtime.Now(),
	})
	if err != nil {
		return gerror.Newf("更新密码失败: %v", err)
	}

	return nil
}

// GenerateCaptcha 生成验证码
func (s *sUser) GenerateCaptcha(ctx context.Context) (captchaId, captchaImage string, err error) {
	// 生成验证码ID
	captchaId = gmd5.MustEncryptString(gconv.String(time.Now().UnixNano()))

	// 生成4位数字验证码
	code := captcha.GenerateCode(4)

	// 将验证码存储到缓存中，有效期5分钟
	cacheKey := fmt.Sprintf("captcha_%s", captchaId)
	gcache.Set(ctx, cacheKey, code, 5*time.Minute)

	// 生成验证码图片
	captchaImage, err = captcha.GenerateImage(code)
	if err != nil {
		return "", "", gerror.Newf("生成验证码图片失败: %v", err)
	}

	return captchaId, captchaImage, nil
}

// VerifyCaptcha 验证验证码
func (s *sUser) VerifyCaptcha(ctx context.Context, captchaId, captcha string) error {
	// 从缓存中获取验证码
	cacheKey := fmt.Sprintf("captcha_%s", captchaId)
	cachedCode, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		return gerror.New("验证码已过期或不存在")
	}

	// 验证验证码
	if gconv.String(cachedCode) != captcha {
		return gerror.New("验证码错误")
	}

	// 验证成功后删除验证码
	gcache.Remove(ctx, cacheKey)
	return nil
}

// GetUserByUsername 根据用户名获取用户
func (s *sUser) GetUserByUsername(ctx context.Context, username string) (user *sysout.UserModel, err error) {
	var entity *entity.User
	err = g.DB().Model("sys_users").Where("username = ? AND deleted_at IS NULL", username).Scan(&entity)
	if err != nil {
		return nil, gerror.Newf("查询用户失败: %v", err)
	}

	if entity == nil {
		return nil, gerror.New("用户不存在")
	}

	return sysout.ConvertToUserModel(entity), nil
}

// ValidateUser 验证用户密码
func (s *sUser) ValidateUser(ctx context.Context, username, password string) (user *sysout.UserModel, err error) {
	// 获取用户信息
	var userEntity *entity.User
	err = g.DB().Model("sys_users").Where("username = ? AND deleted_at IS NULL", username).Scan(&userEntity)
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

	return sysout.ConvertToUserModel(userEntity), nil
}

// UserRoleWithCode 用户角色信息（包含角色编码）
type UserRoleWithCode struct {
	entity.UserRole
	RoleCode string `json:"roleCode" description:"角色编码"`
}

// getUserPrimaryRole 获取用户主要角色
func (s *sUser) getUserPrimaryRole(ctx context.Context, userId int64) (*UserRoleWithCode, error) {
	var userRoleWithCode *UserRoleWithCode
	err := g.DB().Model("sys_user_roles ur").
		LeftJoin("sys_roles r", "ur.role_id = r.id").
		Where("ur.user_id = ? AND ur.is_primary = ? AND r.deleted_at IS NULL", userId, entity.IsPrimaryRole).
		Fields("ur.*, r.code as role_code").
		Scan(&userRoleWithCode)
	if err != nil {
		return nil, err
	}

	if userRoleWithCode == nil {
		return nil, gerror.New("用户没有分配主要角色")
	}

	return userRoleWithCode, nil
}

// generateRefreshToken 生成刷新令牌
func (s *sUser) generateRefreshToken(ctx context.Context, userId int64) (string, error) {
	// 生成随机字符串作为刷新令牌
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// 将刷新令牌存储到缓存中，有效期7天
	cacheKey := fmt.Sprintf("refresh_token_%s", token)
	gcache.Set(ctx, cacheKey, userId, 7*24*time.Hour)

	return token, nil
}

// validateRefreshToken 验证刷新令牌
func (s *sUser) validateRefreshToken(ctx context.Context, token string) (int64, error) {
	cacheKey := fmt.Sprintf("refresh_token_%s", token)
	userIdVal, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		return 0, gerror.New("刷新令牌无效或已过期")
	}

	userId := gconv.Int64(userIdVal)
	if userId <= 0 {
		return 0, gerror.New("刷新令牌无效")
	}

	return userId, nil
}

// updateLoginInfo 更新用户登录信息
func (s *sUser) updateLoginInfo(ctx context.Context, userId int64) error {
	// 这里可以获取客户端IP等信息
	_, err := g.DB().Model("sys_users").Where("id = ?", userId).Update(g.Map{
		"login_at":    gtime.Now(),
		"login_count": gdb.Raw("login_count + 1"),
		"updated_at":  gtime.Now(),
	})
	return err
}

// getUserPermissions 获取用户权限
func (s *sUser) getUserPermissions(ctx context.Context, userId int64) (permissions []string, menuIds []int64, err error) {
	// 通过角色获取权限
	var results []struct {
		Permission string `json:"permission"`
		MenuId     int64  `json:"menu_id"`
	}

	err = g.DB().Model("sys_user_roles ur").
		LeftJoin("sys_role_menus rm", "ur.role_id = rm.role_id").
		LeftJoin("sys_menus m", "rm.menu_id = m.id").
		Where("ur.user_id = ? AND m.deleted_at IS NULL AND m.status = ?", userId, 1).
		Fields("m.permission, m.id as menu_id").
		Scan(&results)
	if err != nil {
		return nil, nil, err
	}

	permissions = make([]string, 0)
	menuIds = make([]int64, 0)
	permissionSet := make(map[string]bool)
	menuIdSet := make(map[int64]bool)

	for _, result := range results {
		if result.Permission != "" && !permissionSet[result.Permission] {
			permissions = append(permissions, result.Permission)
			permissionSet[result.Permission] = true
		}
		if result.MenuId > 0 && !menuIdSet[result.MenuId] {
			menuIds = append(menuIds, result.MenuId)
			menuIdSet[result.MenuId] = true
		}
	}

	return permissions, menuIds, nil
}
