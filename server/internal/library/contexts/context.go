// Package contexts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package contexts

import (
	"client-app/internal/consts"
	"client-app/internal/model"
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改
func Init(r *ghttp.Request, customCtx *model.Context) {
	r.SetCtxVar(consts.ContextHTTPKey, customCtx)
}

// Get 获得上下文变量，如果没有设置，那么返回nil
func Get(ctx context.Context) *model.Context {
	value := ctx.Value(consts.ContextHTTPKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.Context); ok {
		return localCtx
	}
	return nil
}

// SetResponse 设置组件响应 用于访问日志使用
func SetResponse(ctx context.Context, response *model.Response) {
	c := Get(ctx)
	if c == nil {
		g.Log().Warning(ctx, "contexts.SetResponse, c == nil ")
		return
	}
	c.Response = response
}

// SetDataMap 设置额外数据
func SetDataMap(ctx context.Context, vs g.Map) {
	c := Get(ctx)
	if c == nil {
		g.Log().Warning(ctx, "contexts.SetData, c == nil ")
		return
	}

	for k, v := range vs {
		Get(ctx).Data[k] = v
	}
}
