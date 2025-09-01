# 入口点说明

## 主入口文件

主入口文件位于项目根目录的`main.go`，它是程序执行的起点。

```go
// main.go
package main

import (
    "client-app/internal/cmd"
    "client-app/internal/global"
    _ "client-app/internal/logic"
    "github.com/gogf/gf/v2/os/gctx"
)

func main() {
    var ctx = gctx.GetInitCtx()
    global.Init(ctx)
    cmd.Main.Run(ctx)
}
```

## 入口流程说明

1. **初始化上下文**: 通过`gctx.GetInitCtx()`获取初始化上下文
2. **全局初始化**: 调用`global.Init(ctx)`进行全局初始化，包括时区设置和链路追踪初始化
3. **引入包**: 通过`_ "client-app/internal/logic"`导入逻辑层，触发其中的`init()`函数
4. **启动命令**: 执行`cmd.Main.Run(ctx)`启动服务

## 命令行入口

服务启动依赖于`internal/cmd`包中定义的命令。主要命令为 HTTP 服务：

```go
// internal/cmd/http.go
var (
    Http = &gcmd.Command{
        Name:  "http",
        Usage: "http",
        Brief: "HTTP服务，也可以称为主服务，包含http、websocket、tcpserver多个可对外服务",
        Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
            // 初始化http服务
            s := g.Server()

            // 初始化中间件和路由
            // ...

            // 启动服务
            s.Run()
            return
        },
    }
)
```

## HTTP 服务初始化流程

HTTP 服务初始化主要包括以下步骤：

1. **创建服务实例**: `s := g.Server()`
2. **绑定请求钩子**:
   ```go
   // 初始化请求前回调
   s.BindHookHandler("/*any", ghttp.HookBeforeServe, service.Hook().BeforeServe)
   // 请求响应结束后回调
   s.BindHookHandler("/*any", ghttp.HookAfterOutput, service.Hook().AfterOutput)
   ```
3. **注册全局中间件**:
   ```go
   s.BindMiddleware("/*any", []ghttp.HandlerFunc{
       service.Middleware().Ctx,
       service.Middleware().CORS,
       service.Middleware().PreFilter,
       service.Middleware().ResponseHandler,
   }...)
   ```
4. **注册路由**:
   ```go
   s.Group("/", func(group *ghttp.RouterGroup) {
       // 注册Api路由
       router.Api(ctx, group)
   })
   ```
5. **启动服务**: `s.Run()`

## 优雅关闭机制

服务支持优雅关闭，通过信号监听实现：

```go
// 信号监听
signalListen(ctx, signalHandlerForOverall)
go func() {
    <-serverCloseSignal
    _ = s.Shutdown() // 关闭http服务
    g.Log().Debug(ctx, "http successfully closed ..")
    serverWg.Done()
}()
```

当收到操作系统的中断信号(如 SIGINT)时，会触发服务器的优雅关闭流程。

## 全局初始化

全局初始化在`internal/global/init.go`中实现：

```go
func Init(ctx context.Context) {
    // 设置gf运行模式
    SetGFMode(ctx)

    // 默认上海时区
    if err := gtime.SetTimeZone("Asia/Shanghai"); err != nil {
        g.Log().Fatalf(ctx, "时区设置异常 err：%+v", err)
        return
    }

    // 初始化链路追踪
    InitTrace(ctx)
}
```

主要包括：

1. **运行模式设置**: 根据配置设置 GoFrame 的运行模式(开发/测试/生产)
2. **时区设置**: 设置默认时区为上海时区
3. **链路追踪初始化**: 如果配置了 Jaeger 则初始化链路追踪

## 扩展服务启动流程

如需添加新的服务启动命令，需要在`internal/cmd`包中定义新的命令，并在`cmd.go`中注册：

```go
var (
    Main = gcmd.Command{
        Name:  "main",
        Usage: "main",
        Brief: "支付通道服务",
        Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
            return nil
        },
    }
)

func init() {
    Main.AddCommand(Http)
    // 添加新命令
    // Main.AddCommand(NewCommand)
}
```
