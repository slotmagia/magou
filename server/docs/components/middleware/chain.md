# 中间件链

## 中间件概述

中间件是处理 HTTP 请求的过滤器，按照特定顺序执行。在我们的项目中，中间件链定义在`internal/cmd/http.go`文件中。

## 中间件执行顺序

```go
// 注册全局中间件
s.BindMiddleware("/*any", []ghttp.HandlerFunc{
    service.Middleware().Ctx,  // 初始化请求上下文，一般需要第一个进行加载，后续中间件存在依赖关系
    service.Middleware().LogRequest, // 请求日志中间件，在Ctx之后执行以获取正确的traceID
    service.Middleware().CORS, // 跨域中间件，自动处理跨域问题
    //service.Middleware().Blacklist,       // IP黑名单中间件，如果请求IP被后台拉黑，所有请求将被拒绝
    //service.Middleware().DemoLimit,       // 演示系統操作限制，当开启演示模式时，所有POST请求将被拒绝
    service.Middleware().PreFilter,       // 请求输入预处理，api使用gf规范路由并且XxxReq结构体实现了validate.Filter接口即可隐式预处理
    service.Middleware().ResponseHandler, // HTTP响应预处理，在业务处理完成后，对响应结果进行格式化和错误过滤，将处理后的数据发送给请求方
}...)
```

## 执行顺序说明

1. **Ctx**: 初始化请求上下文，必须最先执行，设置链路追踪ID
2. **LogRequest**: 请求日志中间件，记录HTTP请求信息，获取GoFrame的traceID
3. **CORS**: 处理跨域请求，应在业务处理前
4. **PreFilter**: 请求参数预处理，验证和过滤请求数据
5. **ResponseHandler**: 响应处理，统一格式化响应

## 重要说明

- **LogRequest**中间件必须在**Ctx**中间件之后执行，以确保能正确获取链路追踪ID
- **LogRequest**中间件与现有的链路追踪系统完全兼容，使用`gctx.CtxId(ctx)`获取traceID
- 所有日志都会自动包含与Jaeger链路追踪系统一致的traceID

## 路由特定中间件

除了全局中间件外，特定路由组还可以添加额外的中间件：

```go
// API 签名验证
//group.Middleware(service.Middleware().ApiVerify)

// API 认证中间件
group.Middleware(service.Middleware().ApiAuth)
```

## 中间件开发规范

1. **结构**: 中间件实现在`internal/logic/middleware`目录
2. **接口定义**: 接口在`internal/service/middleware.go`中定义
3. **注册方式**: 在`init()`函数中通过`service.RegisterMiddleware()`注册

## 中间件编写示例

```go
// 中间件实现
func (s *sMiddleware) Ctx(r *ghttp.Request) {
    // 初始化请求上下文
    customCtx := &model.Context{
        // 初始化上下文数据
    }
    contexts.Init(r, customCtx)
    r.Middleware.Next()
}

// 中间件注册
func init() {
    service.RegisterMiddleware(New())
}
```

## 中间件注意事项

1. 中间件应保持单一职责原则
2. 中间件应处理好异常情况，避免中断请求链
3. 性能敏感的中间件应放在链的前端
4. 所有中间件应正确调用`r.Middleware.Next()`向下传递请求
