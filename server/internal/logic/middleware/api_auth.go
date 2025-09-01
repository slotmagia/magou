package middleware

import (
	"client-app/internal/consts"
	"client-app/internal/library/contexts"
	"client-app/internal/library/response"
	"client-app/internal/model"
	"client-app/internal/model/entity"
	"client-app/internal/service"
	"client-app/utility/simple"
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
)

// ApiAuth API鉴权中间件
func (s *sMiddleware) ApiAuth(r *ghttp.Request) {
	var (
		ctx  = r.Context()
		path = gstr.Replace(r.URL.Path, simple.RouterPrefix(ctx, consts.AppApi), "", 1)
	)

	// 不需要验证登录的路由地址
	if s.IsExceptLogin(ctx, consts.AppApi, path) {
		r.Middleware.Next()
		return
	}

	// 从请求头获取Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		s.authFailed(r, consts.ErrTokenMissing, "访问令牌不能为空")
		return
	}

	// 提取Token
	token, err := simple.ExtractTokenFromHeader(authHeader)
	if err != nil {
		s.authFailed(r, consts.ErrTokenInvalid, err.Error())
		return
	}

	// 解析JWT Token
	secretKey := simple.GetJWTSecretKey(ctx)
	payload, err := simple.ParseJWTToken(token, secretKey)
	if err != nil {
		if strings.Contains(err.Error(), "过期") {
			s.authFailed(r, consts.ErrTokenExpired, consts.GetAuthErrorMessage(consts.ErrTokenExpired))
		} else {
			s.authFailed(r, consts.ErrTokenInvalid, consts.GetAuthErrorMessage(consts.ErrTokenInvalid))
		}
		return
	}

	// 验证用户是否存在且状态正常
	user, err := s.getUserFromPayload(ctx, payload)
	if err != nil {
		s.authFailed(r, consts.ErrUserNotFound, err.Error())
		return
	}

	// 验证用户状态
	if err := s.validateUserStatus(user); err != nil {
		s.authFailed(r, consts.ErrUserDisabled, err.Error())
		return
	}

	// 设置用户身份到上下文
	identity := s.buildIdentity(user, payload)
	s.setUserToContext(r, identity)

	// 不需要验证权限的路由地址
	if s.IsExceptAuth(ctx, consts.AppApi, path) {
		r.Middleware.Next()
		return
	}

	// 验证API访问权限
	if err := s.checkAPIPermission(ctx, user.Id, path, r.Method); err != nil {
		s.authFailed(r, consts.ErrPermissionDenied, err.Error())
		return
	}

	// 继续执行下一个中间件
	r.Middleware.Next()
}

// getUserFromPayload 从JWT payload获取用户信息
func (s *sMiddleware) getUserFromPayload(ctx context.Context, payload *simple.JWTPayload) (*entity.User, error) {
	var user *entity.User

	// 从数据库查询用户信息
	err := g.DB().Model("sys_users").Where("id = ? AND deleted_at IS NULL", payload.UserId).Scan(&user)
	if err != nil {
		return nil, gerror.Newf("查询用户信息失败: %v", err)
	}

	if user == nil {
		return nil, gerror.New("用户不存在或已被删除")
	}

	return user, nil
}

// validateUserStatus 验证用户状态
func (s *sMiddleware) validateUserStatus(user *entity.User) error {
	if user.Status == entity.UserStatusDisabled {
		return gerror.New(consts.GetAuthErrorMessage(consts.ErrUserDisabled))
	}
	if user.Status == entity.UserStatusLocked {
		return gerror.New(consts.GetAuthErrorMessage(consts.ErrUserLocked))
	}
	return nil
}

// buildIdentity 构建用户身份信息
func (s *sMiddleware) buildIdentity(user *entity.User, payload *simple.JWTPayload) *model.Identity {
	return &model.Identity{
		Id:       user.Id,
		Pid:      0, // 如果有上级关系，这里需要从数据库查询
		DeptId:   user.DeptId,
		DeptType: "", // 如果有部门类型，这里需要从部门表查询
		RoleId:   payload.RoleId,
		RoleKey:  payload.RoleKey,
		Username: user.Username,
		RealName: user.RealName,
		Avatar:   user.Avatar,
		Email:    user.Email,
		Mobile:   user.Phone,
		App:      payload.App,
		LoginAt:  gtime.Now(),
	}
}

// setUserToContext 设置用户信息到上下文
func (s *sMiddleware) setUserToContext(r *ghttp.Request, identity *model.Identity) {
	// 获取现有的上下文
	customCtx := contexts.Get(r.Context())
	if customCtx == nil {
		// 如果没有上下文，创建一个新的
		customCtx = &model.Context{
			Data:   make(g.Map),
			Module: getModule(r.URL.Path),
		}
		contexts.Init(r, customCtx)
	}

	// 设置用户身份信息
	customCtx.User = identity
}

// checkAPIPermission 检查API访问权限
func (s *sMiddleware) checkAPIPermission(ctx context.Context, userId int64, path string, method string) error {
	// 构造权限标识，通常是 path:method 的格式
	permission := s.buildPermissionKey(path, method)

	// 检查用户是否有该权限
	hasPermission, err := service.Role().CheckUserPermission(ctx, userId, permission)
	if err != nil {
		return gerror.Newf("权限检查失败: %v", err)
	}

	if !hasPermission {
		return gerror.New("您没有访问该接口的权限")
	}

	return nil
}

// buildPermissionKey 构建权限标识
func (s *sMiddleware) buildPermissionKey(path string, method string) string {
	// 清理路径，移除开头的斜杠
	cleanPath := strings.TrimPrefix(path, "/")

	// 将路径中的斜杠替换为冒号，构成权限标识
	// 例如: /role/list -> role:list
	permission := strings.ReplaceAll(cleanPath, "/", ":")

	// 添加HTTP方法作为后缀（可选）
	// permission = fmt.Sprintf("%s:%s", permission, strings.ToLower(method))

	return permission
}

// authFailed 认证失败处理
func (s *sMiddleware) authFailed(r *ghttp.Request, errCode string, message string) {
	// 记录认证失败日志
	g.Log().Warningf(r.Context(), "API认证失败: %s, Path: %s, IP: %s, UA: %s",
		message, r.URL.Path, r.GetClientIp(), r.Header.Get("User-Agent"))

	// 返回统一的错误响应
	response.JsonExit(r, 401, message)
}

// GetCurrentUser 获取当前登录用户信息
func (s *sMiddleware) GetCurrentUser(ctx context.Context) *model.Identity {
	customCtx := contexts.Get(ctx)
	if customCtx != nil && customCtx.User != nil {
		return customCtx.User
	}
	return nil
}
