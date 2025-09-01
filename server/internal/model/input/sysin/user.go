package sysin

import (
	"client-app/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserListInp 用户列表查询参数
type UserListInp struct {
	Page      int    `json:"page"         v:"min:1"              d:"1"          description:"页码"`
	PageSize  int    `json:"pageSize"     v:"min:1|max:100"      d:"20"         description:"每页数量"`
	Username  string `json:"username"     v:"length:0,50"                       description:"用户名"`
	RealName  string `json:"realName"     v:"length:0,50"                       description:"真实姓名"`
	Email     string `json:"email"        v:"length:0,100"                      description:"邮箱地址"`
	Phone     string `json:"phone"        v:"length:0,20"                       description:"手机号码"`
	Status    int    `json:"status"       v:"in:-1,1,2,3"        d:"-1"        description:"状态：-1=全部 1=正常 2=锁定 3=禁用"`
	DeptId    int64  `json:"deptId"       v:"min:0"                             description:"部门ID"`
	RoleId    int64  `json:"roleId"       v:"min:0"                             description:"角色ID"`
	Gender    int    `json:"gender"       v:"in:-1,0,1,2"        d:"-1"        description:"性别：-1=全部 0=未知 1=男 2=女"`
	StartDate string `json:"startDate"    v:"date"                              description:"开始日期"`
	EndDate   string `json:"endDate"      v:"date"                              description:"结束日期"`
	OrderBy   string `json:"orderBy"      v:"in:id,username,created_at" d:"id"  description:"排序字段"`
	OrderType string `json:"orderType"    v:"in:asc,desc"       d:"desc"       description:"排序方式"`
}

// Filter 参数过滤和验证
func (inp *UserListInp) Filter(ctx context.Context) error {
	return g.Validator().Rules(inp.getValidationRules()).Data(inp).Run(ctx)
}

func (inp *UserListInp) getValidationRules() string {
	return `
		page@页码: min:1
		pageSize@每页数量: min:1|max:100
		username@用户名: length:0,50
		realName@真实姓名: length:0,50
		email@邮箱地址: length:0,100
		phone@手机号码: length:0,20
		status@状态: in:-1,1,2,3
		deptId@部门ID: min:0
		roleId@角色ID: min:0
		gender@性别: in:-1,0,1,2
		startDate@开始日期: date
		endDate@结束日期: date
		orderBy@排序字段: in:id,username,created_at
		orderType@排序方式: in:asc,desc
	`
}

// UserDetailInp 用户详情查询参数
type UserDetailInp struct {
	Id int64 `json:"id" v:"required|min:1" description:"用户ID"`
}

// Filter 参数过滤和验证
func (inp *UserDetailInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// CreateUserInp 创建用户参数
type CreateUserInp struct {
	Username         string      `json:"username"         v:"required|length:3,50"              description:"用户名"`
	Email            string      `json:"email"            v:"email|length:0,100"                description:"邮箱地址"`
	Phone            string      `json:"phone"            v:"phone|length:0,20"                 description:"手机号码"`
	Password         string      `json:"password"         v:"required|length:6,32"              description:"密码"`
	ConfirmPassword  string      `json:"confirmPassword"  v:"required|same:password"            description:"确认密码"`
	RealName         string      `json:"realName"         v:"required|length:2,50"              description:"真实姓名"`
	Nickname         string      `json:"nickname"         v:"length:0,50"                       description:"昵称"`
	Avatar           string      `json:"avatar"           v:"url|length:0,255"                  description:"头像URL"`
	Gender           int         `json:"gender"           v:"in:0,1,2"             d:"0"        description:"性别：0=未知 1=男 2=女"`
	Birthday         *gtime.Time `json:"birthday"         v:"date"                              description:"生日"`
	DeptId           int64       `json:"deptId"           v:"min:0"                             description:"部门ID"`
	Position         string      `json:"position"         v:"length:0,100"                      description:"职位"`
	Status           int         `json:"status"           v:"in:1,2,3"             d:"1"        description:"状态：1=正常 2=锁定 3=禁用"`
	RoleIds          []int64     `json:"roleIds"          v:"required"                          description:"角色ID列表"`
	TwoFactorEnabled int         `json:"twoFactorEnabled" v:"in:0,1"               d:"0"        description:"是否启用双因子认证"`
	Remark           string      `json:"remark"           v:"length:0,500"                      description:"备注说明"`
}

// Filter 参数过滤和验证
func (inp *CreateUserInp) Filter(ctx context.Context) error {
	// 自定义验证规则
	rules := map[string]string{
		"username":         "required|length:3,50",
		"email":            "email|length:0,100",
		"phone":            "phone|length:0,20",
		"password":         "required|length:6,32",
		"confirmPassword":  "required|same:password",
		"realName":         "required|length:2,50",
		"nickname":         "length:0,50",
		"avatar":           "url|length:0,255",
		"gender":           "in:0,1,2",
		"birthday":         "date",
		"deptId":           "min:0",
		"position":         "length:0,100",
		"status":           "in:1,2,3",
		"roleIds":          "required",
		"twoFactorEnabled": "in:0,1",
		"remark":           "length:0,500",
	}

	return g.Validator().Rules(rules).Data(inp).Run(ctx)
}

// UpdateUserInp 更新用户参数
type UpdateUserInp struct {
	Id               int64       `json:"id"               v:"required|min:1"                    description:"用户ID"`
	Username         string      `json:"username"         v:"required|length:3,50"              description:"用户名"`
	Email            string      `json:"email"            v:"email|length:0,100"                description:"邮箱地址"`
	Phone            string      `json:"phone"            v:"phone|length:0,20"                 description:"手机号码"`
	RealName         string      `json:"realName"         v:"required|length:2,50"              description:"真实姓名"`
	Nickname         string      `json:"nickname"         v:"length:0,50"                       description:"昵称"`
	Avatar           string      `json:"avatar"           v:"url|length:0,255"                  description:"头像URL"`
	Gender           int         `json:"gender"           v:"in:0,1,2"             d:"0"        description:"性别：0=未知 1=男 2=女"`
	Birthday         *gtime.Time `json:"birthday"         v:"date"                              description:"生日"`
	DeptId           int64       `json:"deptId"           v:"min:0"                             description:"部门ID"`
	Position         string      `json:"position"         v:"length:0,100"                      description:"职位"`
	Status           int         `json:"status"           v:"in:1,2,3"                         description:"状态：1=正常 2=锁定 3=禁用"`
	RoleIds          []int64     `json:"roleIds"          v:"required"                          description:"角色ID列表"`
	TwoFactorEnabled int         `json:"twoFactorEnabled" v:"in:0,1"               d:"0"        description:"是否启用双因子认证"`
	Remark           string      `json:"remark"           v:"length:0,500"                      description:"备注说明"`
}

// Filter 参数过滤和验证
func (inp *UpdateUserInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// DeleteUserInp 删除用户参数
type DeleteUserInp struct {
	Id int64 `json:"id" v:"required|min:1" description:"用户ID"`
}

// Filter 参数过滤和验证
func (inp *DeleteUserInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// BatchDeleteUserInp 批量删除用户参数
type BatchDeleteUserInp struct {
	Ids []int64 `json:"ids" v:"required|length:1,100" description:"用户ID列表"`
}

// Filter 参数过滤和验证
func (inp *BatchDeleteUserInp) Filter(ctx context.Context) error {
	if len(inp.Ids) == 0 {
		return gerror.New("用户ID列表不能为空")
	}
	if len(inp.Ids) > 100 {
		return gerror.New("批量删除数量不能超过100个")
	}
	return nil
}

// UpdateUserStatusInp 更新用户状态参数
type UpdateUserStatusInp struct {
	Id     int64 `json:"id"     v:"required|min:1"  description:"用户ID"`
	Status int   `json:"status" v:"required|in:1,2,3" description:"状态：1=正常 2=锁定 3=禁用"`
}

// Filter 参数过滤和验证
func (inp *UpdateUserStatusInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// ResetPasswordInp 重置密码参数
type ResetPasswordInp struct {
	Id              int64  `json:"id"              v:"required|min:1"      description:"用户ID"`
	NewPassword     string `json:"newPassword"     v:"required|length:6,32" description:"新密码"`
	ConfirmPassword string `json:"confirmPassword" v:"required|same:newPassword" description:"确认密码"`
}

// Filter 参数过滤和验证
func (inp *ResetPasswordInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// ChangePasswordInp 修改密码参数
type ChangePasswordInp struct {
	Id              int64  `json:"id"              v:"required|min:1"        description:"用户ID"`
	OldPassword     string `json:"oldPassword"     v:"required|length:6,32"  description:"原密码"`
	NewPassword     string `json:"newPassword"     v:"required|length:6,32"  description:"新密码"`
	ConfirmPassword string `json:"confirmPassword" v:"required|same:newPassword" description:"确认密码"`
}

// Filter 参数过滤和验证
func (inp *ChangePasswordInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// UserProfileInp 用户资料查询参数
type UserProfileInp struct {
	Id int64 `json:"id" v:"required|min:1" description:"用户ID"`
}

// Filter 参数过滤和验证
func (inp *UserProfileInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// UpdateProfileInp 更新用户资料参数
type UpdateProfileInp struct {
	Id       int64       `json:"id"       v:"required|min:1"       description:"用户ID"`
	RealName string      `json:"realName" v:"required|length:2,50" description:"真实姓名"`
	Nickname string      `json:"nickname" v:"length:0,50"          description:"昵称"`
	Avatar   string      `json:"avatar"   v:"url|length:0,255"     description:"头像URL"`
	Gender   int         `json:"gender"   v:"in:0,1,2"             description:"性别：0=未知 1=男 2=女"`
	Birthday *gtime.Time `json:"birthday" v:"date"                 description:"生日"`
	Phone    string      `json:"phone"    v:"phone|length:0,20"    description:"手机号码"`
	Email    string      `json:"email"    v:"email|length:0,100"   description:"邮箱地址"`
}

// Filter 参数过滤和验证
func (inp *UpdateProfileInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// UserStatsInp 用户统计查询参数
type UserStatsInp struct {
	StartDate string `json:"startDate" v:"date" description:"开始日期"`
	EndDate   string `json:"endDate"   v:"date" description:"结束日期"`
}

// Filter 参数过滤和验证
func (inp *UserStatsInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// AssignUserRoleInp 分配用户角色参数
type AssignUserRoleInp struct {
	UserId  int64   `json:"userId"  v:"required|min:1" description:"用户ID"`
	RoleIds []int64 `json:"roleIds" v:"required"       description:"角色ID列表"`
}

// Filter 参数过滤和验证
func (inp *AssignUserRoleInp) Filter(ctx context.Context) error {
	if err := g.Validator().Data(inp).Run(ctx); err != nil {
		return err
	}
	if len(inp.RoleIds) == 0 {
		return gerror.New("角色ID列表不能为空")
	}
	return nil
}

// RemoveUserRoleInp 移除用户角色参数
type RemoveUserRoleInp struct {
	UserId  int64   `json:"userId"  v:"required|min:1" description:"用户ID"`
	RoleIds []int64 `json:"roleIds" v:"required"       description:"角色ID列表"`
}

// Filter 参数过滤和验证
func (inp *RemoveUserRoleInp) Filter(ctx context.Context) error {
	if err := g.Validator().Data(inp).Run(ctx); err != nil {
		return err
	}
	if len(inp.RoleIds) == 0 {
		return gerror.New("角色ID列表不能为空")
	}
	return nil
}

// UserRoleInp 查询用户角色参数
type UserRoleInp struct {
	UserId int64 `json:"userId" v:"required|min:1" description:"用户ID"`
}

// Filter 参数过滤和验证
func (inp *UserRoleInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// SetPrimaryRoleInp 设置主要角色参数
type SetPrimaryRoleInp struct {
	UserId int64 `json:"userId" v:"required|min:1" description:"用户ID"`
	RoleId int64 `json:"roleId" v:"required|min:1" description:"角色ID"`
}

// Filter 参数过滤和验证
func (inp *SetPrimaryRoleInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// UserPermissionInp 查询用户权限参数
type UserPermissionInp struct {
	UserId int64 `json:"userId" v:"required|min:1" description:"用户ID"`
}

// Filter 参数过滤和验证
func (inp *UserPermissionInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// CheckPermissionInp 检查权限参数
type CheckPermissionInp struct {
	UserId     int64  `json:"userId"     v:"required|min:1" description:"用户ID"`
	Permission string `json:"permission" v:"required"       description:"权限标识"`
}

// Filter 参数过滤和验证
func (inp *CheckPermissionInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// BatchCheckPermissionInp 批量检查权限参数
type BatchCheckPermissionInp struct {
	UserIds    []int64 `json:"userIds"    v:"required" description:"用户ID列表"`
	Permission string  `json:"permission" v:"required" description:"权限标识"`
}

// Filter 参数过滤和验证
func (inp *BatchCheckPermissionInp) Filter(ctx context.Context) error {
	if err := g.Validator().Data(inp).Run(ctx); err != nil {
		return err
	}
	if len(inp.UserIds) == 0 {
		return gerror.New("用户ID列表不能为空")
	}
	return nil
}

// UserLoginInp 用户登录参数
type UserLoginInp struct {
	TenantCode string `json:"tenantCode" v:"required|length:1,50"  description:"租户编码"`
	Username   string `json:"username"   v:"required|length:3,50"  description:"用户名"`
	Password   string `json:"password"   v:"required|length:6,32"  description:"密码"`
	Captcha    string `json:"captcha"    v:"required|length:4,6"   description:"验证码"`
	CaptchaId  string `json:"captchaId"  v:"required"              description:"验证码ID"`
	RememberMe bool   `json:"rememberMe" d:"false"                 description:"记住我"`
}

// Filter 参数过滤和验证
func (inp *UserLoginInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// UserLogoutInp 用户登出参数
type UserLogoutInp struct {
	Token string `json:"token" v:"required" description:"访问令牌"`
}

// Filter 参数过滤和验证
func (inp *UserLogoutInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// RefreshTokenInp 刷新令牌参数
type RefreshTokenInp struct {
	RefreshToken string `json:"refreshToken" v:"required" description:"刷新令牌"`
}

// Filter 参数过滤和验证
func (inp *RefreshTokenInp) Filter(ctx context.Context) error {
	return g.Validator().Data(inp).Run(ctx)
}

// ValidateUserStatus 验证用户状态
func ValidateUserStatus(status int) error {
	if !entity.ValidateUserStatus(status) {
		return gerror.New("无效的用户状态")
	}
	return nil
}

// ValidateUserGender 验证用户性别
func ValidateUserGender(gender int) error {
	if !entity.ValidateUserGender(gender) {
		return gerror.New("无效的用户性别")
	}
	return nil
}

// ValidateUsername 验证用户名
func ValidateUsername(ctx context.Context, username string) error {
	if err := g.Validator().Rules("required|length:3,50").Messages("用户名|用户名是必填项|用户名长度必须在3-50字符之间").Data(username).Run(ctx); err != nil {
		return err
	}

	// 检查用户名格式（只允许字母、数字、下划线）
	if err := g.Validator().Rules("regex:^[a-zA-Z0-9_]+$").Messages("用户名格式不正确，只能包含字母、数字和下划线").Data(username).Run(ctx); err != nil {
		return gerror.New("用户名只能包含字母、数字和下划线")
	}

	return nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(ctx context.Context, password string) error {
	if err := g.Validator().Rules("required|length:6,32").Messages("密码|密码是必填项|密码长度必须在6-32字符之间").Data(password).Run(ctx); err != nil {
		return err
	}

	// 检查密码复杂度（至少包含数字和字母）
	hasNumber := false
	hasLetter := false

	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
		}
	}

	if !hasNumber || !hasLetter {
		return gerror.New("密码必须包含字母和数字")
	}

	return nil
}
