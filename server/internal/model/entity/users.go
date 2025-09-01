package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// 用户状态常量
const (
	UserStatusNormal   = 1 // 正常
	UserStatusLocked   = 2 // 锁定
	UserStatusDisabled = 3 // 禁用
)

// 用户性别常量
const (
	UserGenderUnknown = 0 // 未知
	UserGenderMale    = 1 // 男
	UserGenderFemale  = 2 // 女
)

// 用户角色关联常量
const (
	IsNotPrimaryRole = 0 // 非主要角色
	IsPrimaryRole    = 1 // 主要角色
)

// 角色编码常量
const (
	RoleCodeSuperAdmin  = "super_admin"  // 超级管理员
	RoleCodeSystemAdmin = "system_admin" // 系统管理员
	RoleCodeNormalUser  = "normal_user"  // 普通用户
)

// User 用户实体
type User struct {
	Id                   int64       `json:"id"                    description:"主键ID"`
	Username             string      `json:"username"              description:"用户名"`
	Email                string      `json:"email"                 description:"邮箱地址"`
	Phone                string      `json:"phone"                 description:"手机号码"`
	Password             string      `json:"-"                     description:"密码（不返回给前端）"`
	Salt                 string      `json:"-"                     description:"密码盐值（不返回给前端）"`
	RealName             string      `json:"realName"              description:"真实姓名"`
	Nickname             string      `json:"nickname"              description:"昵称"`
	Avatar               string      `json:"avatar"                description:"头像URL"`
	Gender               int         `json:"gender"                description:"性别"`
	Birthday             *gtime.Time `json:"birthday"              description:"生日"`
	DeptId               int64       `json:"deptId"                description:"部门ID"`
	Position             string      `json:"position"              description:"职位"`
	Status               int         `json:"status"                description:"状态"`
	LoginIp              string      `json:"loginIp"               description:"最后登录IP"`
	LoginAt              *gtime.Time `json:"loginAt"               description:"最后登录时间"`
	LoginCount           int         `json:"loginCount"            description:"登录次数"`
	PasswordResetToken   string      `json:"-"                     description:"密码重置令牌"`
	PasswordResetExpires *gtime.Time `json:"-"                     description:"密码重置过期时间"`
	EmailVerifiedAt      *gtime.Time `json:"emailVerifiedAt"       description:"邮箱验证时间"`
	PhoneVerifiedAt      *gtime.Time `json:"phoneVerifiedAt"       description:"手机验证时间"`
	TwoFactorEnabled     int         `json:"twoFactorEnabled"      description:"是否启用双因子认证"`
	TwoFactorSecret      string      `json:"-"                     description:"双因子认证密钥"`
	Remark               string      `json:"remark"                description:"备注说明"`
	CreatedBy            int64       `json:"createdBy"             description:"创建人ID"`
	UpdatedBy            int64       `json:"updatedBy"             description:"修改人ID"`
	CreatedAt            *gtime.Time `json:"createdAt"             description:"创建时间"`
	UpdatedAt            *gtime.Time `json:"updatedAt"             description:"更新时间"`
	DeletedAt            *gtime.Time `json:"deletedAt,omitempty"   description:"删除时间"`
}

// UserRole 用户角色关联实体
type UserRole struct {
	Id           int64       `json:"id"           description:"主键ID"`
	UserId       int64       `json:"userId"       description:"用户ID"`
	RoleId       int64       `json:"roleId"       description:"角色ID"`
	IsPrimaryVal int         `json:"isPrimary"    description:"是否主要角色"`
	AssignedBy   int64       `json:"assignedBy"   description:"分配人ID"`
	ExpiresAt    *gtime.Time `json:"expiresAt"    description:"过期时间"`
	CreatedAt    *gtime.Time `json:"createdAt"    description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    description:"更新时间"`
}

// TwoFactorStatus 双因子认证状态常量
const (
	TwoFactorDisabled = 0 // 禁用
	TwoFactorEnabled  = 1 // 启用
)

// IsActive 判断用户是否活跃
func (u *User) IsActive() bool {
	return u.Status == UserStatusNormal && u.DeletedAt == nil
}

// IsLocked 判断用户是否被锁定
func (u *User) IsLocked() bool {
	return u.Status == UserStatusLocked
}

// IsDisabled 判断用户是否被禁用
func (u *User) IsDisabled() bool {
	return u.Status == UserStatusDisabled
}

// IsDeleted 判断用户是否被删除
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// IsEmailVerified 判断邮箱是否已验证
func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

// IsPhoneVerified 判断手机是否已验证
func (u *User) IsPhoneVerified() bool {
	return u.PhoneVerifiedAt != nil
}

// IsTwoFactorEnabled 判断是否启用了双因子认证
func (u *User) IsTwoFactorEnabled() bool {
	return u.TwoFactorEnabled == TwoFactorEnabled
}

// GetStatusName 获取状态名称
func (u *User) GetStatusName() string {
	switch u.Status {
	case UserStatusNormal:
		return "正常"
	case UserStatusLocked:
		return "锁定"
	case UserStatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}

// GetGenderName 获取性别名称
func (u *User) GetGenderName() string {
	switch u.Gender {
	case UserGenderMale:
		return "男"
	case UserGenderFemale:
		return "女"
	default:
		return "未知"
	}
}

// IsPasswordResetExpired 判断密码重置令牌是否过期
func (u *User) IsPasswordResetExpired() bool {
	if u.PasswordResetExpires == nil {
		return true
	}
	return u.PasswordResetExpires.Before(gtime.Now())
}

// CanLogin 判断用户是否可以登录
func (u *User) CanLogin() bool {
	return u.IsActive() && !u.IsLocked() && !u.IsDeleted()
}

// IsPrimary 判断用户角色关联是否为主要角色
func (ur *UserRole) IsPrimary() bool {
	return ur.IsPrimaryVal == IsPrimaryRole
}

// IsExpired 判断用户角色是否过期
func (ur *UserRole) IsExpired() bool {
	if ur.ExpiresAt == nil {
		return false // 永不过期
	}
	return ur.ExpiresAt.Before(gtime.Now())
}

// IsValid 判断用户角色关联是否有效
func (ur *UserRole) IsValid() bool {
	return !ur.IsExpired()
}

// UserWithRoles 带角色信息的用户
type UserWithRoles struct {
	User
	Roles       []*Role  `json:"roles"       description:"用户角色列表"`
	RoleIds     []int64  `json:"roleIds"     description:"角色ID列表"`
	RoleCodes   []string `json:"roleCodes"   description:"角色编码列表"`
	RoleNames   []string `json:"roleNames"   description:"角色名称列表"`
	PrimaryRole *Role    `json:"primaryRole" description:"主要角色"`
	Permissions []string `json:"permissions" description:"权限列表"`
	MenuIds     []int64  `json:"menuIds"     description:"菜单ID列表"`
}

// UserProfile 用户详细信息
type UserProfile struct {
	User
	DeptName        string     `json:"deptName"         description:"部门名称"`
	RoleInfo        []*Role    `json:"roleInfo"         description:"角色详情"`
	LastLoginInfo   *LoginInfo `json:"lastLoginInfo"    description:"最后登录信息"`
	PermissionCount int        `json:"permissionCount"  description:"权限数量"`
	SecurityLevel   int        `json:"securityLevel"    description:"安全等级"`
}

// LoginInfo 登录信息
type LoginInfo struct {
	Ip        string      `json:"ip"        description:"登录IP"`
	Location  string      `json:"location"  description:"登录地点"`
	UserAgent string      `json:"userAgent" description:"用户代理"`
	LoginAt   *gtime.Time `json:"loginAt"   description:"登录时间"`
}

// UserStats 用户统计信息
type UserStats struct {
	TotalCount    int64 `json:"totalCount"      description:"总用户数"`
	ActiveCount   int64 `json:"activeCount"     description:"活跃用户数"`
	LockedCount   int64 `json:"lockedCount"     description:"锁定用户数"`
	DisabledCount int64 `json:"disabledCount"   description:"禁用用户数"`
	OnlineCount   int64 `json:"onlineCount"     description:"在线用户数"`
	NewUserToday  int64 `json:"newUserToday"    description:"今日新增用户"`
	LoginToday    int64 `json:"loginToday"      description:"今日登录用户"`
}

// ValidateUserStatus 验证用户状态是否有效
func ValidateUserStatus(status int) bool {
	return status >= UserStatusNormal && status <= UserStatusDisabled
}

// ValidateUserGender 验证用户性别是否有效
func ValidateUserGender(gender int) bool {
	return gender >= UserGenderUnknown && gender <= UserGenderFemale
}

// GetAllUserStatuses 获取所有用户状态选项
func GetAllUserStatuses() map[int]string {
	return map[int]string{
		UserStatusNormal:   "正常",
		UserStatusLocked:   "锁定",
		UserStatusDisabled: "禁用",
	}
}

// GetAllUserGenders 获取所有性别选项
func GetAllUserGenders() map[int]string {
	return map[int]string{
		UserGenderUnknown: "未知",
		UserGenderMale:    "男",
		UserGenderFemale:  "女",
	}
}

// HasRole 判断用户是否拥有指定角色
func (uwr *UserWithRoles) HasRole(roleCode string) bool {
	for _, code := range uwr.RoleCodes {
		if code == roleCode {
			return true
		}
	}
	return false
}

// HasPermission 判断用户是否拥有指定权限
func (uwr *UserWithRoles) HasPermission(permission string) bool {
	for _, perm := range uwr.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// HasAnyRole 判断用户是否拥有任一指定角色
func (uwr *UserWithRoles) HasAnyRole(roleCodes []string) bool {
	for _, roleCode := range roleCodes {
		if uwr.HasRole(roleCode) {
			return true
		}
	}
	return false
}

// HasAllRoles 判断用户是否拥有所有指定角色
func (uwr *UserWithRoles) HasAllRoles(roleCodes []string) bool {
	for _, roleCode := range roleCodes {
		if !uwr.HasRole(roleCode) {
			return false
		}
	}
	return true
}

// IsSuperAdmin 判断用户是否为超级管理员
func (uwr *UserWithRoles) IsSuperAdmin() bool {
	return uwr.HasRole(RoleCodeSuperAdmin)
}

// IsSystemAdmin 判断用户是否为系统管理员
func (uwr *UserWithRoles) IsSystemAdmin() bool {
	return uwr.HasRole(RoleCodeSystemAdmin)
}

// GetMaxDataScope 获取用户的最大数据权限范围
func (uwr *UserWithRoles) GetMaxDataScope() int {
	maxScope := DataScopeSelf // 默认最小权限
	for _, role := range uwr.Roles {
		if role.DataScope < maxScope { // 数字越小权限越大
			maxScope = role.DataScope
		}
	}
	return maxScope
}

// CanAccessAllData 判断用户是否可以访问全部数据
func (uwr *UserWithRoles) CanAccessAllData() bool {
	return uwr.GetMaxDataScope() == DataScopeAll
}

// CanAccessDeptData 判断用户是否可以访问部门数据
func (uwr *UserWithRoles) CanAccessDeptData() bool {
	scope := uwr.GetMaxDataScope()
	return scope == DataScopeAll || scope == DataScopeDept || scope == DataScopeDeptAndSub
}

// GetSecurityLevel 计算用户安全等级
func (up *UserProfile) GetSecurityLevel() int {
	level := 1 // 基础等级

	if up.IsEmailVerified() {
		level++
	}
	if up.IsPhoneVerified() {
		level++
	}
	if up.IsTwoFactorEnabled() {
		level += 2
	}
	if len(up.RoleInfo) > 0 {
		level++
	}

	return level
}
