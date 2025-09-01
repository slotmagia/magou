// Package router
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package router

import (
	"client-app/internal/consts"
	"client-app/internal/controller/api"
	"client-app/internal/service"
	"client-app/utility/simple"
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Api 前台路由
func Api(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group(simple.RouterPrefix(ctx, consts.AppApi), func(group *ghttp.RouterGroup) {
		// 不需要认证的公开接口
		group.Bind(
			api.User, // 用户认证接口
		)

		// API 签名验证
		//group.Middleware(service.Middleware().ApiVerify)
		group.Bind()

		// 需要认证的受保护接口
		group.Middleware(service.Middleware().ApiAuth)
		group.Bind(
			api.Role, // 角色管理接口
			api.Menu,
			api.NewTenant(),
		)
	})
}
