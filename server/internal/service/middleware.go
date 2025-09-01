// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"client-app/internal/model"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IMiddleware interface {
		// ApiAuth API鉴权中间件
		ApiAuth(r *ghttp.Request)

		// TenantFilter 租户过滤中间件
		TenantFilter(r *ghttp.Request)
		
		// TenantAuth 租户权限验证中间件
		TenantAuth(r *ghttp.Request)

		// Ctx 初始化请求上下文
		Ctx(r *ghttp.Request)
		// CORS allows Cross-origin resource sharing.
		CORS(r *ghttp.Request)
		// DemoLimit 演示系统操作限制
		//DemoLimit(r *ghttp.Request)

		// IsExceptAuth 是否是不需要验证权限的路由地址
		IsExceptAuth(ctx context.Context, appName string, path string) bool
		// IsExceptLogin 是否是不需要登录的路由地址
		IsExceptLogin(ctx context.Context, appName string, path string) bool
		// Blacklist IP黑名单限制中间件
		//Blacklist(r *ghttp.Request)
		// Develop 开发工具白名单过滤
		//Develop(r *ghttp.Request)
		// PreFilter 请求输入预处理
		// api使用gf规范路由并且XxxReq结构体实现了validate.Filter接口即可
		PreFilter(r *ghttp.Request)
		// ResponseHandler HTTP响应预处理
		ResponseHandler(r *ghttp.Request)

		// GetCurrentUser 获取当前登录用户信息
		GetCurrentUser(ctx context.Context) *model.Identity

		// LogRequest 请求日志中间件
		LogRequest(r *ghttp.Request)

		// LogError 记录错误日志
		LogError(ctx context.Context, err error, message string, args ...interface{})

		// LogBusiness 记录业务日志
		LogBusiness(ctx context.Context, action string, data interface{}, result string)

		// LogSecurity 记录安全日志
		LogSecurity(ctx context.Context, event string, level string, details ...interface{})

		// LogMetrics 记录性能指标
		LogMetrics(ctx context.Context, metric string, value interface{}, tags ...string)

		// LogAudit 记录审计日志
		LogAudit(ctx context.Context, action string, resource string, result string, details ...interface{})
	}
)

var (
	localMiddleware IMiddleware
)

func Middleware() IMiddleware {
	if localMiddleware == nil {
		panic("implement not found for interface IMiddleware, forgot register?")
	}
	return localMiddleware
}

func RegisterMiddleware(i IMiddleware) {
	localMiddleware = i
}
