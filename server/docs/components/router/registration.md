# 路由注册

## 路由系统概述

路由系统负责将 HTTP 请求映射到对应的控制器方法，是服务架构中的重要组成部分。本项目使用 GoFrame 框架的路由系统进行路由管理。

## 路由注册方式

路由注册主要在`internal/router`目录下完成，使用了 GoFrame 的路由分组和绑定机制：

```go
// internal/router/api.go
func Api(ctx context.Context, group *ghttp.RouterGroup) {
    group.Group(simple.RouterPrefix(ctx, consts.AppApi), func(group *ghttp.RouterGroup) {
        group.Bind(
            api.NewPayment(),
        )
        // API 签名验证
        //group.Middleware(service.Middleware().ApiVerify)
        group.Bind()
        group.Middleware(service.Middleware().ApiAuth)
        group.Bind()
    })
}
```

## 路由分组策略

项目采用了基于业务模块的路由分组策略：

1. **根分组**: 所有 API 路由都在根路由(`/`)下
2. **业务分组**: 在根分组下按业务模块进行分组，如`/api/payment`

路由前缀通过配置和工具函数动态生成：

```go
simple.RouterPrefix(ctx, consts.AppApi)
```

这样可以方便地通过配置修改 API 路径前缀。

## 控制器绑定

GoFrame 提供了自动路由注册机制，通过`Bind`方法将控制器结构体绑定到路由组：

```go
group.Bind(
    api.NewPayment(), // 绑定支付控制器
)
```

框架会自动根据控制器的方法名生成路由，例如：

- 控制器方法: `Payment.Submit`
- 生成路由: `/api/payment/submit`

## 中间件绑定

路由组可以绑定特定的中间件，作用于该组下的所有路由：

```go
// 绑定API认证中间件
group.Middleware(service.Middleware().ApiAuth)
```

这些中间件将在全局中间件之后执行。

## 路由扫描机制

项目支持在启动时扫描所有已注册的路由，并保存到全局变量中，便于调试和监控：

```go
// internal/global/httproutes.go
func GetRoutes() []*RouterItem {
    return routes
}
```

## 路由命名规范

路由命名遵循以下规范：

1. **URL 路径**: 使用小写字母，单词之间用短横线分隔

   - 正确: `/api/payment/query-order`
   - 错误: `/api/Payment/QueryOrder`

2. **HTTP 方法**: 默认使用 POST 方法，可通过标签指定其他方法

3. **版本控制**: 通过目录结构实现 API 版本控制
   - v1 版本: `/api/payment/submit`
   - v2 版本: `/api/v2/payment/submit`

## 添加新路由步骤

1. **创建控制器**: 在`internal/controller`目录下创建控制器结构体

2. **实现控制器方法**: 添加对应的控制器方法

3. **绑定路由**: 在`internal/router`中的对应函数中添加控制器绑定

```go
// 添加新控制器绑定
group.Bind(
    api.NewPayment(),
    api.NewOrder(), // 新添加的订单控制器
)
```

## 路由调试

在开发环境中，可以通过访问默认的路由信息接口查看所有已注册的路由：

```
GET /swagger
```

或者通过日志查看注册的路由信息：

```
[INFO] 2023-03-27 10:00:00.000 | http server started listening on [:8080]
[INFO] 2023-03-27 10:00:00.000 | registered http route: POST:/api/payment/submit
```

## 路由安全

1. **认证中间件**: 敏感路由必须添加认证中间件
2. **请求方法限制**: 根据 RESTful 原则限制请求方法
3. **路由频率限制**: 对敏感操作添加频率限制
