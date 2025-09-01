// Package hook
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package hook

import (
	"client-app/internal/library/contexts"
	"client-app/utility/simple"
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"strings"
)

// 忽略的请求方式
var ignoredRequestMethods = []string{"HEAD", "PRI"}

// accessLog 访问日志
func (s *sHook) accessLog(r *ghttp.Request) {
	if s.isIgnoredRequest(r) {
		return
	}

	var ctx = r.Context()
	if contexts.Get(ctx) == nil {
		return
	}

	contexts.SetDataMap(ctx, g.Map{
		"request.takeUpTime": gtime.Now().Sub(gtime.New(r.EnterTime)).Milliseconds(),
		// ...
	})

	// 记录访问日志
	simple.SafeGo(ctx, func(ctx context.Context) {
		g.Log().Infof(ctx, "访问日志记录 - 请求路径: %s, 方法: %s, 耗时: %dms, 客户端IP: %s",
			r.URL.Path, r.Method, contexts.Get(ctx).Data["request.takeUpTime"], r.GetClientIp())
	})
}

// isIgnoredRequest 是否忽略请求
func (s *sHook) isIgnoredRequest(r *ghttp.Request) bool {
	if r.IsFileRequest() {
		return true
	}

	if gstr.InArray(ignoredRequestMethods, strings.ToUpper(r.Method)) {
		return true
	}
	return false
}
