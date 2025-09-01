package api

import (
	"client-app/internal/api/v1/user"
	"client-app/internal/service"
	"context"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	User = &cUser{}
)

type cUser struct{}

// Login 用户登录
func (c *cUser) Login(ctx context.Context, req *user.UserLoginReq) (res *user.UserLoginRes, err error) {
	// 调用服务层处理登录逻辑
	out, err := service.User().Login(ctx, &req.UserLoginInp)
	if err != nil {
		return nil, err
	}

	res = new(user.UserLoginRes)
	if g.IsEmpty(out) {
		return res, nil
	}
	res.LoginTokenModel = out
	return res, nil
}

// Logout 用户退出
func (c *cUser) Logout(ctx context.Context, req *user.UserLogoutReq) (res *user.UserLogoutRes, err error) {
	// 调用服务层处理退出逻辑
	err = service.User().Logout(ctx, &req.UserLogoutInp)
	if err != nil {
		return nil, err
	}

	res = &user.UserLogoutRes{
		Success: true,
		Message: "退出成功",
	}
	return res, nil
}

// Profile 获取用户信息
func (c *cUser) Profile(ctx context.Context, req *user.UserProfileReq) (res *user.UserProfileRes, err error) {
	// 获取当前用户信息
	currentUser := service.Middleware().GetCurrentUser(ctx)
	if currentUser == nil {
		return nil, gerror.New("用户未登录")
	}

	// 调用服务层获取用户详细信息
	out, err := service.User().GetProfile(ctx, currentUser.Id)
	if err != nil {
		return nil, err
	}

	res = new(user.UserProfileRes)
	if g.IsEmpty(out) {
		return res, nil
	}
	res.UserModel = out
	return res, nil
}

// RefreshToken 刷新访问令牌
func (c *cUser) RefreshToken(ctx context.Context, req *user.UserRefreshTokenReq) (res *user.UserRefreshTokenRes, err error) {
	// 调用服务层刷新令牌
	out, err := service.User().RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	res = &user.UserRefreshTokenRes{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		ExpiresIn:    out.ExpiresIn,
	}
	return res, nil
}

// ChangePassword 修改密码
func (c *cUser) ChangePassword(ctx context.Context, req *user.UserChangePasswordReq) (res *user.UserChangePasswordRes, err error) {
	// 获取当前用户信息
	currentUser := service.Middleware().GetCurrentUser(ctx)
	if currentUser == nil {
		return nil, gerror.New("用户未登录")
	}

	// 调用服务层修改密码
	err = service.User().ChangePassword(ctx, currentUser.Id, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}

	res = &user.UserChangePasswordRes{
		Success: true,
		Message: "密码修改成功",
	}
	return res, nil
}

// GetCaptcha 获取验证码
func (c *cUser) GetCaptcha(ctx context.Context, req *user.GetCaptchaReq) (res *user.GetCaptchaRes, err error) {

	// 调用服务层生成验证码
	captchaId, captchaImage, err := service.User().GenerateCaptcha(ctx)
	if err != nil {
		return nil, err
	}

	res = &user.GetCaptchaRes{
		CaptchaId:    captchaId,
		CaptchaImage: captchaImage,
	}
	return res, nil
}
