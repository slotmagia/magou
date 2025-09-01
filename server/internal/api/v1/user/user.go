package user

import (
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"
	"github.com/gogf/gf/v2/frame/g"
)

// UserLoginReq 用户登录请求
type UserLoginReq struct {
	g.Meta `path:"/login" method:"post" summary:"用户登录" tags:"用户认证"`
	sysin.UserLoginInp
}

// UserLoginRes 用户登录响应
type UserLoginRes struct {
	*sysout.LoginTokenModel
}

// UserLogoutReq 用户退出请求
type UserLogoutReq struct {
	g.Meta `path:"/logout" method:"post" summary:"用户退出" tags:"用户认证"`
	sysin.UserLogoutInp
}

// UserLogoutRes 用户退出响应
type UserLogoutRes struct {
	Success bool   `json:"success" description:"是否成功"`
	Message string `json:"message" description:"提示信息"`
}

// UserProfileReq 获取用户信息请求
type UserProfileReq struct {
	g.Meta `path:"/profile" method:"get" summary:"获取用户信息" tags:"用户认证"`
}

// UserProfileRes 获取用户信息响应
type UserProfileRes struct {
	*sysout.UserModel
}

// UserRefreshTokenReq 刷新Token请求
type UserRefreshTokenReq struct {
	g.Meta       `path:"/refresh-token" method:"post" summary:"刷新访问令牌" tags:"用户认证"`
	RefreshToken string `json:"refreshToken" v:"required" description:"刷新令牌"`
}

// UserRefreshTokenRes 刷新Token响应
type UserRefreshTokenRes struct {
	AccessToken  string `json:"accessToken"  description:"新的访问令牌"`
	RefreshToken string `json:"refreshToken" description:"新的刷新令牌"`
	ExpiresIn    int64  `json:"expiresIn"    description:"过期时间（秒）"`
}

// UserChangePasswordReq 修改密码请求
type UserChangePasswordReq struct {
	g.Meta          `path:"/change-password" method:"post" summary:"修改密码" tags:"用户认证"`
	OldPassword     string `json:"oldPassword" v:"required|length:6,32" description:"原密码"`
	NewPassword     string `json:"newPassword" v:"required|length:6,32" description:"新密码"`
	ConfirmPassword string `json:"confirmPassword" v:"required|same:NewPassword" description:"确认新密码"`
}

// UserChangePasswordRes 修改密码响应
type UserChangePasswordRes struct {
	Success bool   `json:"success" description:"是否成功"`
	Message string `json:"message" description:"提示信息"`
}

// GetCaptchaReq 获取验证码请求
type GetCaptchaReq struct {
	g.Meta `path:"/captcha" method:"get" summary:"获取验证码" tags:"用户认证"`
}

// GetCaptchaRes 获取验证码响应
type GetCaptchaRes struct {
	CaptchaId    string `json:"captchaId"    description:"验证码ID"`
	CaptchaImage string `json:"captchaImage" description:"验证码图片（base64）"`
}
