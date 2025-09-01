package sysout

import (
	"client-app/internal/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
)

// UserModel 用户基础信息模型
type UserModel struct {
	Id               int64       `json:"id"               description:"主键ID"`
	Username         string      `json:"username"         description:"用户名"`
	Email            string      `json:"email"            description:"邮箱地址"`
	Phone            string      `json:"phone"            description:"手机号码"`
	RealName         string      `json:"realName"         description:"真实姓名"`
	Nickname         string      `json:"nickname"         description:"昵称"`
	Avatar           string      `json:"avatar"           description:"头像URL"`
	Gender           int         `json:"gender"           description:"性别"`
	GenderName       string      `json:"genderName"       description:"性别名称"`
	Birthday         *gtime.Time `json:"birthday"         description:"生日"`
	DeptId           int64       `json:"deptId"           description:"部门ID"`
	DeptName         string      `json:"deptName"         description:"部门名称"`
	Position         string      `json:"position"         description:"职位"`
	Status           int         `json:"status"           description:"状态"`
	StatusName       string      `json:"statusName"       description:"状态名称"`
	LoginIp          string      `json:"loginIp"          description:"最后登录IP"`
	LoginAt          *gtime.Time `json:"loginAt"          description:"最后登录时间"`
	LoginCount       int         `json:"loginCount"       description:"登录次数"`
	EmailVerifiedAt  *gtime.Time `json:"emailVerifiedAt"  description:"邮箱验证时间"`
	PhoneVerifiedAt  *gtime.Time `json:"phoneVerifiedAt"  description:"手机验证时间"`
	TwoFactorEnabled int         `json:"twoFactorEnabled" description:"是否启用双因子认证"`
	SecurityLevel    int         `json:"securityLevel"    description:"安全等级"`
	Remark           string      `json:"remark"           description:"备注说明"`
	CreatedBy        int64       `json:"createdBy"        description:"创建人ID"`
	UpdatedBy        int64       `json:"updatedBy"        description:"修改人ID"`
	CreatedAt        *gtime.Time `json:"createdAt"        description:"创建时间"`
	UpdatedAt        *gtime.Time `json:"updatedAt"        description:"更新时间"`
}

// UserListModel 用户列表响应模型
type UserListModel struct {
	List     []*UserModel `json:"list"     description:"用户列表"`
	Total    int          `json:"total"    description:"总数量"`
	Page     int          `json:"page"     description:"当前页码"`
	PageSize int          `json:"pageSize" description:"每页数量"`
}

// UserDetailModel 用户详情响应模型
type UserDetailModel struct {
	UserModel
	RoleIds     []int64      `json:"roleIds"     description:"角色ID列表"`
	RoleCodes   []string     `json:"roleCodes"   description:"角色编码列表"`
	RoleNames   []string     `json:"roleNames"   description:"角色名称列表"`
	Roles       []*RoleModel `json:"roles"       description:"角色详情列表"`
	PrimaryRole *RoleModel   `json:"primaryRole" description:"主要角色"`
	Permissions []string     `json:"permissions" description:"权限列表"`
	MenuIds     []int64      `json:"menuIds"     description:"菜单ID列表"`
	DataScope   int          `json:"dataScope"   description:"数据权限范围"`
}

// UserStatsModel 用户统计信息模型
type UserStatsModel struct {
	TotalCount       int64            `json:"totalCount"       description:"总用户数"`
	ActiveCount      int64            `json:"activeCount"      description:"活跃用户数"`
	LockedCount      int64            `json:"lockedCount"      description:"锁定用户数"`
	DisabledCount    int64            `json:"disabledCount"    description:"禁用用户数"`
	OnlineCount      int64            `json:"onlineCount"      description:"在线用户数"`
	NewUserToday     int64            `json:"newUserToday"     description:"今日新增用户"`
	LoginToday       int64            `json:"loginToday"       description:"今日登录用户"`
	EmailVerified    int64            `json:"emailVerified"    description:"邮箱已验证用户"`
	PhoneVerified    int64            `json:"phoneVerified"    description:"手机已验证用户"`
	TwoFactorEnabled int64            `json:"twoFactorEnabled" description:"启用双因子认证用户"`
	GenderStats      map[string]int64 `json:"genderStats"      description:"性别统计"`
	StatusStats      map[string]int64 `json:"statusStats"      description:"状态统计"`
}

// LoginTokenModel 登录令牌响应模型
type LoginTokenModel struct {
	AccessToken  string     `json:"accessToken"  description:"访问令牌"`
	RefreshToken string     `json:"refreshToken" description:"刷新令牌"`
	TokenType    string     `json:"tokenType"    description:"令牌类型"`
	ExpiresIn    int64      `json:"expiresIn"    description:"过期时间（秒）"`
	UserInfo     *UserModel `json:"userInfo"     description:"用户信息"`
	Permissions  []string   `json:"permissions"  description:"权限列表"`
	MenuIds      []int64    `json:"menuIds"      description:"菜单ID列表"`
}

// ConvertToUserModel 将用户实体转换为用户模型
func ConvertToUserModel(user *entity.User) *UserModel {
	if user == nil {
		return nil
	}

	return &UserModel{
		Id:               user.Id,
		Username:         user.Username,
		Email:            user.Email,
		Phone:            user.Phone,
		RealName:         user.RealName,
		Nickname:         user.Nickname,
		Avatar:           user.Avatar,
		Gender:           user.Gender,
		GenderName:       user.GetGenderName(),
		Birthday:         user.Birthday,
		DeptId:           user.DeptId,
		Position:         user.Position,
		Status:           user.Status,
		StatusName:       user.GetStatusName(),
		LoginIp:          user.LoginIp,
		LoginAt:          user.LoginAt,
		LoginCount:       user.LoginCount,
		EmailVerifiedAt:  user.EmailVerifiedAt,
		PhoneVerifiedAt:  user.PhoneVerifiedAt,
		TwoFactorEnabled: user.TwoFactorEnabled,
		SecurityLevel:    calculateSecurityLevel(user),
		Remark:           user.Remark,
		CreatedBy:        user.CreatedBy,
		UpdatedBy:        user.UpdatedBy,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}
}

// ConvertToUserDetailModel 将用户信息转换为用户详情模型
func ConvertToUserDetailModel(user *entity.User, roles []*RoleModel, primaryRole *RoleModel, permissions []string, menuIds []int64, dataScope int) *UserDetailModel {
	if user == nil {
		return nil
	}

	userModel := ConvertToUserModel(user)

	// 提取角色信息
	roleIds := make([]int64, len(roles))
	roleCodes := make([]string, len(roles))
	roleNames := make([]string, len(roles))

	for i, role := range roles {
		roleIds[i] = role.Id
		roleCodes[i] = role.Code
		roleNames[i] = role.Name
	}

	return &UserDetailModel{
		UserModel:   *userModel,
		RoleIds:     roleIds,
		RoleCodes:   roleCodes,
		RoleNames:   roleNames,
		Roles:       roles,
		PrimaryRole: primaryRole,
		Permissions: permissions,
		MenuIds:     menuIds,
		DataScope:   dataScope,
	}
}

// BuildUserStatsModel 构建用户统计模型
func BuildUserStatsModel(users []*entity.User) *UserStatsModel {
	stats := &UserStatsModel{
		TotalCount:  int64(len(users)),
		GenderStats: make(map[string]int64),
		StatusStats: make(map[string]int64),
	}

	for _, user := range users {
		// 状态统计
		switch user.Status {
		case entity.UserStatusNormal:
			stats.ActiveCount++
			stats.StatusStats["正常"]++
		case entity.UserStatusLocked:
			stats.LockedCount++
			stats.StatusStats["锁定"]++
		case entity.UserStatusDisabled:
			stats.DisabledCount++
			stats.StatusStats["禁用"]++
		}

		// 性别统计
		genderName := user.GetGenderName()
		stats.GenderStats[genderName]++

		// 验证状态统计
		if user.IsEmailVerified() {
			stats.EmailVerified++
		}
		if user.IsPhoneVerified() {
			stats.PhoneVerified++
		}
		if user.IsTwoFactorEnabled() {
			stats.TwoFactorEnabled++
		}
	}

	return stats
}

// calculateSecurityLevel 计算用户安全等级
func calculateSecurityLevel(user *entity.User) int {
	level := 1 // 基础等级

	if user.IsEmailVerified() {
		level++
	}
	if user.IsPhoneVerified() {
		level++
	}
	if user.IsTwoFactorEnabled() {
		level += 2
	}

	return level
}

// ConvertToLoginTokenModel 将登录信息转换为令牌模型
func ConvertToLoginTokenModel(accessToken, refreshToken string, expiresIn int64,
	user *entity.User, permissions []string, menuIds []int64) *LoginTokenModel {

	userModel := ConvertToUserModel(user)

	return &LoginTokenModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		UserInfo:     userModel,
		Permissions:  permissions,
		MenuIds:      menuIds,
	}
}
