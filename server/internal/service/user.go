// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"
	"context"
)

type (
	IUser interface {
		// Login 用户登录
		Login(ctx context.Context, in *sysin.UserLoginInp) (res *sysout.LoginTokenModel, err error)
		
		// Logout 用户退出
		Logout(ctx context.Context, in *sysin.UserLogoutInp) error
		
		// GetProfile 获取用户资料
		GetProfile(ctx context.Context, userId int64) (res *sysout.UserModel, err error)
		
		// RefreshToken 刷新访问令牌
		RefreshToken(ctx context.Context, refreshToken string) (res *TokenInfo, err error)
		
		// ChangePassword 修改密码
		ChangePassword(ctx context.Context, userId int64, oldPassword, newPassword string) error
		
		// GenerateCaptcha 生成验证码
		GenerateCaptcha(ctx context.Context) (captchaId, captchaImage string, err error)
		
		// VerifyCaptcha 验证验证码
		VerifyCaptcha(ctx context.Context, captchaId, captcha string) error
		
		// GetUserByUsername 根据用户名获取用户
		GetUserByUsername(ctx context.Context, username string) (user *sysout.UserModel, err error)
		
		// ValidateUser 验证用户密码
		ValidateUser(ctx context.Context, username, password string) (user *sysout.UserModel, err error)
	}
)

// TokenInfo 令牌信息
type TokenInfo struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
} 