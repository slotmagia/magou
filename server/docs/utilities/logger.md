# 日志系统设计与实现

## 概述

本项目实现了一个完整的分层日志系统，基于GoFrame框架的glog包构建，支持多环境配置、结构化日志、日志轮转、性能监控等功能。日志系统设计遵循开发文档的架构规范，提供了统一的日志记录接口和中间件。

## 系统架构

### 1. 日志架构层次

```
┌─────────────────────────────────────────┐
│              应用层                      │
│  Controller/Logic层调用日志接口         │
├─────────────────────────────────────────┤
│              服务层                      │
│  Service接口定义日志方法                │
├─────────────────────────────────────────┤
│              中间件层                    │
│  LogRequest中间件、业务日志中间件       │
├─────────────────────────────────────────┤
│              工具层                      │
│  Logger工具类封装GoFrame的glog          │
├─────────────────────────────────────────┤
│              GoFrame层                   │
│  glog提供底层日志功能                   │
└─────────────────────────────────────────┘
```

### 2. 日志分类

- **HTTP请求日志**: 记录所有HTTP请求和响应信息
- **业务操作日志**: 记录关键业务操作的执行过程
- **错误日志**: 记录系统错误和异常信息
- **安全日志**: 记录安全相关事件
- **审计日志**: 记录资源操作的审计信息
- **性能指标日志**: 记录系统性能指标

## 配置系统

### 1. 配置文件结构

```yaml
# manifest/config/config.yaml
logger:
  # 日志级别: debug, info, warn, error, fatal
  level: "debug"
  # 日志格式: text, json
  format: "text"
  # 是否输出到控制台
  stdoutPrint: true
  # 日志文件路径
  path: "logs"
  # 日志文件名
  file: "app.log"
  # 是否异步写入
  async: true
  # 切割大小 (100MB)
  rotateSize: 104857600
  # 切割过期时间 (7天)
  rotateExpire: "168h"
  # 切割数量
  rotateCount: 10
  # 是否备份
  rotateBackup: true
  # 是否压缩
  rotateCompress: true
  # 上下文键名
  ctxKeys: ["request_id", "trace_id", "user_id", "method", "path"]
  # 请求头键名
  headerKeys: ["X-Request-ID", "X-Trace-ID", "X-User-ID"]
```

### 2. 环境相关配置

不同环境下的日志配置会自动调整：

**开发环境 (develop)**:
- 级别: debug
- 控制台输出: true
- 彩色输出: true
- 显示调用栈: true

**测试环境 (testing)**:
- 级别: debug
- 控制台输出: true
- 彩色输出: false
- 显示调用栈: true

**预发布环境 (staging)**:
- 级别: info
- 控制台输出: false
- 显示调用栈: false

**生产环境 (product)**:
- 级别: warn
- 控制台输出: false
- 显示调用栈: false

## 核心组件

### 1. Logger工具类

位于`utility/logger/logger.go`，提供核心日志功能：

```go
// 创建日志实例
logger := logger.NewFromConfig(ctx)

// 记录不同级别的日志
logger.Debug(ctx, "调试信息")
logger.Info(ctx, "信息日志")
logger.Warn(ctx, "警告日志")
logger.Error(ctx, "错误日志")
logger.Fatal(ctx, "致命错误")

// 格式化日志
logger.Infof(ctx, "用户 %s 执行了 %s 操作", username, action)
```

### 2. 日志中间件

位于`internal/logic/middleware/log.go`，提供HTTP请求日志记录：

```go
// 在HTTP服务中注册日志中间件
s.BindMiddleware("/*any", []ghttp.HandlerFunc{
    service.Middleware().LogRequest, // 日志中间件
    // 其他中间件...
}...)
```

### 3. 服务接口

位于`internal/service/middleware.go`，定义日志相关的服务接口：

```go
type IMiddleware interface {
    // 请求日志中间件
    LogRequest(r *ghttp.Request)
    
    // 业务日志
    LogBusiness(ctx context.Context, action string, data interface{}, result string)
    
    // 错误日志
    LogError(ctx context.Context, err error, message string, args ...interface{})
    
    // 安全日志
    LogSecurity(ctx context.Context, event string, level string, details ...interface{})
    
    // 审计日志
    LogAudit(ctx context.Context, action string, resource string, result string, details ...interface{})
    
    // 性能指标日志
    LogMetrics(ctx context.Context, metric string, value interface{}, tags ...string)
}
```

## 使用方法

### 1. HTTP请求日志

HTTP请求日志通过中间件自动记录，包含以下信息：

```
2024-01-01 10:00:00.000 [trace_id=def456 request_id=abc123 method=POST path=/api/role/create user_id=1001] HTTP请求完成 POST /api/role/create 200 156ms 1024b
```

**重要说明**：日志系统与现有的GoFrame链路追踪系统完全兼容：
- `trace_id`: 使用GoFrame的`gctx.CtxId(ctx)`，与Jaeger链路追踪保持一致
- `request_id`: HTTP请求的唯一标识，可通过`X-Request-ID`请求头传递
- 日志中间件在`Ctx`中间件之后执行，确保正确获取链路追踪信息

### 2. 业务操作日志

在控制器或逻辑层记录业务操作：

```go
// 在控制器中记录业务日志
func (c *cRole) CreateRole(ctx context.Context, req *role.CreateRoleReq) (res *role.CreateRoleRes, err error) {
    // 记录操作开始
    service.Middleware().LogBusiness(ctx, "CreateRole", req, "开始")
    
    out, err := service.Role().CreateRole(ctx, &req.CreateRoleInp)
    if err != nil {
        service.Middleware().LogError(ctx, err, "创建角色失败")
        return nil, err
    }
    
    // 记录操作成功
    service.Middleware().LogBusiness(ctx, "CreateRole", out, "成功")
    service.Middleware().LogAudit(ctx, "CREATE", "ROLE", "SUCCESS", "创建角色", out.Id)
    
    return &role.CreateRoleRes{RoleModel: out}, nil
}
```

### 3. 错误日志记录

```go
// 记录错误日志
if err := someOperation(); err != nil {
    service.Middleware().LogError(ctx, err, "操作失败: %s", operation)
    return err
}
```

### 4. 安全日志记录

```go
// 记录安全事件
service.Middleware().LogSecurity(ctx, "UNAUTHORIZED_ACCESS", "high", "用户尝试访问未授权资源", userID, resource)
```

### 5. 审计日志记录

```go
// 记录资源操作审计
service.Middleware().LogAudit(ctx, "DELETE", "ROLE", "SUCCESS", "删除角色", roleID)
```

### 6. 性能指标日志

```go
// 记录性能指标
service.Middleware().LogMetrics(ctx, "database_query_time", duration.Milliseconds(), "table=users", "operation=select")
```

## 日志格式

### 1. 文本格式 (默认)

```
2024-01-01 10:00:00.000 [trace_id=def456 request_id=abc123 user_id=1001] INFO 业务操作 CreateRole 结果: 成功
```

### 2. JSON格式

```json
{
  "timestamp": "2024-01-01 10:00:00.000",
  "level": "info",
  "env": "develop",
  "app": "hotgo",
  "version": "1.0.0",
  "message": "业务操作 CreateRole 结果: 成功",
  "trace_id": "def456",
  "request_id": "abc123",
  "user_id": "1001",
  "method": "POST",
  "path": "/api/role/create",
  "caller": "role.go:25"
}
```

**注意**：`trace_id`始终使用GoFrame框架提供的链路追踪ID，确保与现有系统的兼容性。

## 日志轮转

### 1. 按大小轮转

- 默认100MB切割
- 保留10个备份文件
- 支持gzip压缩

### 2. 按时间轮转

- 默认7天过期
- 自动清理过期文件

### 3. 文件结构

```
logs/
├── develop/                    # 开发环境日志
│   ├── app.log                # 当前日志文件
│   ├── app.log.2024010101.gz  # 压缩备份文件
│   └── app.log.2024010102.gz
├── testing/                   # 测试环境日志
└── product/                   # 生产环境日志
```

## 性能优化

### 1. 异步写入

```go
// 配置中启用异步写入
logger:
  async: true
```

### 2. 敏感信息过滤

中间件自动过滤敏感字段：

```go
// 自动过滤的敏感字段
sensitiveFields := []string{
    "password", "token", "secret", "key", "authorization",
    "credit", "card", "account", "payment", "billing",
}
```

### 3. 日志级别判断

```go
// 避免不必要的字符串操作
if logger.IsDebugEnabled() {
    expensiveDebugInfo := generateDebugInfo()
    logger.Debug(ctx, "调试信息:", expensiveDebugInfo)
}
```

## 监控告警

### 1. 慢请求监控

```go
// 自动记录超过2秒的慢请求
if duration > 2*time.Second {
    logger.Warn(ctx, "慢请求告警", method, path, duration)
}
```

### 2. 错误率监控

```go
// 记录4xx和5xx错误
if status >= 400 {
    level := "warn"
    if status >= 500 {
        level = "error"
    }
    logger.Log(level, ctx, "错误响应", status, method, path)
}
```

### 3. 系统指标监控

```go
// 记录系统性能指标
service.Middleware().LogMetrics(ctx, "memory_usage", memStats.Alloc)
service.Middleware().LogMetrics(ctx, "goroutine_count", runtime.NumGoroutine())
```

## 链路追踪集成

日志系统完全集成了GoFrame的链路追踪机制：

### 1. Jaeger配置

```yaml
# manifest/config/config.yaml
jaeger:
  # 是否启用链路追踪
  switch: true
  # Jaeger agent地址
  endpoint: "http://localhost:14268/api/traces"
```

### 2. traceID的生成和使用

```go
// 在Ctx中间件中初始化链路追踪
if g.Cfg().MustGet(r.Context(), "jaeger.switch").Bool() {
    ctx, span := gtrace.NewSpan(r.Context(), "middleware.ctx")
    span.SetAttributes(attribute.KeyValue{
        Key:   "traceID",
        Value: attribute.StringValue(gctx.CtxId(ctx)),
    })
    span.End()
    r.SetCtx(ctx)
}

// 在日志中间件中获取traceID
traceID := gctx.CtxId(ctx) // 与Jaeger保持一致
```

### 3. 中间件执行顺序

确保正确的中间件执行顺序以获取traceID：

```go
s.BindMiddleware("/*any", []ghttp.HandlerFunc{
    service.Middleware().Ctx,        // 1. 初始化上下文和链路追踪
    service.Middleware().LogRequest, // 2. 日志中间件获取traceID
    service.Middleware().CORS,       // 3. 其他中间件...
}...)
```

## 访问日志

访问日志通过Hook机制实现，记录在`internal/logic/hook/access_log.go`：

```go
// 记录访问日志
func (s *sHook) accessLog(r *ghttp.Request) {
    if s.isIgnoredRequest(r) {
        return
    }
    
    // 记录访问信息
    g.Log().Infof(ctx, "访问日志记录 - 请求路径: %s, 方法: %s, 耗时: %dms, 客户端IP: %s",
        r.URL.Path, r.Method, takeUpTime, r.GetClientIp())
}
```

## 最佳实践

### 1. 上下文传递

```go
// 始终传递上下文
func handleRequest(ctx context.Context) error {
    logger.Info(ctx, "开始处理请求")
    // 业务逻辑...
    logger.Info(ctx, "请求处理完成")
    return nil
}
```

### 2. 结构化日志

```go
// 使用上下文添加结构化信息
ctx = context.WithValue(ctx, "user_id", userID)
ctx = context.WithValue(ctx, "action", "create_role")
logger.Info(ctx, "执行业务操作")
```

### 3. 错误处理

```go
// 完整的错误处理和日志记录
if err := doSomething(); err != nil {
    service.Middleware().LogError(ctx, err, "操作失败")
    return gerror.Wrap(err, "业务操作失败")
}
```

### 4. 安全考虑

```go
// 避免记录敏感信息
// 错误示例
logger.Info(ctx, "用户登录", "password", password) // 不要这样做

// 正确示例
logger.Info(ctx, "用户登录", "username", username) // 只记录用户名
```

## 常见问题

### 1. 如何在不同环境使用不同配置？

系统会根据`system.mode`配置自动调整日志级别和输出方式。

### 2. 如何自定义日志格式？

```go
// 创建自定义配置
config := logger.DefaultConfig()
config.Format = "json"
config.TimeFormat = "2006-01-02T15:04:05.000Z"

logger := logger.NewLogger(config)
```

### 3. 如何添加自定义日志字段？

```go
// 在上下文中添加自定义字段
ctx = context.WithValue(ctx, "custom_field", "custom_value")
logger.Info(ctx, "自定义日志")
```

### 4. 如何优化大并发下的日志性能？

1. 启用异步写入
2. 适当调整日志级别
3. 避免记录过多调试信息
4. 使用条件判断减少不必要的日志

## 完整使用示例

以下是一个完整的示例，展示如何在业务代码中使用与现有traceID系统一致的日志：

```go
// 控制器示例
func (c *cRole) CreateRole(ctx context.Context, req *role.CreateRoleReq) (res *role.CreateRoleRes, err error) {
    // traceID已经由Ctx中间件自动设置，LogRequest中间件会正确获取
    // 这里的ctx包含完整的链路追踪信息
    
    service.Middleware().LogBusiness(ctx, "CreateRole", req, "开始")
    
    // 业务逻辑...
    out, err := service.Role().CreateRole(ctx, &req.CreateRoleInp)
    if err != nil {
        // 错误日志会自动包含traceID，便于问题追踪
        service.Middleware().LogError(ctx, err, "创建角色失败")
        return nil, err
    }
    
    // 成功日志和审计日志
    service.Middleware().LogBusiness(ctx, "CreateRole", out, "成功")
    service.Middleware().LogAudit(ctx, "CREATE", "ROLE", "SUCCESS", "创建角色", out.Id)
    
    return &role.CreateRoleRes{RoleModel: out}, nil
}
```

实际日志输出：
```
2024-01-01 10:00:00.000 [trace_id=abc123-def456-ghi789 request_id=req-001 method=POST path=/api/role/create user_id=1001] INFO 业务操作 CreateRole 结果: 开始
2024-01-01 10:00:00.156 [trace_id=abc123-def456-ghi789 request_id=req-001 method=POST path=/api/role/create user_id=1001] INFO 业务操作 CreateRole 结果: 成功
2024-01-01 10:00:00.157 [trace_id=abc123-def456-ghi789 request_id=req-001 method=POST path=/api/role/create user_id=1001] INFO 审计日志 操作: CREATE 资源: ROLE 结果: SUCCESS 详情: [创建角色 123]
```

## 与微服务架构的集成

在微服务架构中，traceID将在服务间自动传递：

```go
// 调用下游服务时，traceID会自动传递
client := g.Client()
client.SetHeader("X-Trace-ID", gctx.CtxId(ctx))
response := client.PostContent(ctx, downstreamURL, data)
```

## 验证traceID一致性

### 1. 开发环境验证

```go
// 在业务代码中验证traceID一致性
func (c *cRole) CreateRole(ctx context.Context, req *role.CreateRoleReq) (res *role.CreateRoleRes, err error) {
    // 获取当前请求的traceID
    traceID := gctx.CtxId(ctx)
    
    // 如果开启了调试模式，可以打印验证
    if logger.IsDebugEnabled() {
        g.Log().Debugf(ctx, "当前请求traceID: %s", traceID)
    }
    
    // 确保所有日志都包含相同的traceID
    service.Middleware().LogBusiness(ctx, "CreateRole", req, "开始")
    // ... 业务逻辑
}
```

### 2. 日志输出验证

所有相关日志应包含相同的traceID：

```
2024-01-01 10:00:00.000 [trace_id=abc123-def456-ghi789 request_id=req-001] HTTP请求开始 POST /api/role/create
2024-01-01 10:00:00.001 [trace_id=abc123-def456-ghi789 request_id=req-001] INFO 业务操作 CreateRole 结果: 开始  
2024-01-01 10:00:00.156 [trace_id=abc123-def456-ghi789 request_id=req-001] INFO 业务操作 CreateRole 结果: 成功
2024-01-01 10:00:00.157 [trace_id=abc123-def456-ghi789 request_id=req-001] HTTP请求完成 POST /api/role/create 200 156ms
```

### 3. Jaeger UI验证

如果启用了Jaeger链路追踪，在Jaeger UI中可以看到：
- 日志中的traceID与Jaeger中的TraceID完全一致
- 可以通过traceID关联日志和链路追踪数据
- 便于问题排查和性能分析

## 总结

本日志系统提供了完整的日志记录解决方案，包含：

- **完美兼容**: 与现有GoFrame链路追踪系统无缝集成
- **traceID一致性**: 使用`gctx.CtxId(ctx)`确保与Jaeger链路追踪保持一致
- **多环境配置**: 根据运行环境自动调整日志行为
- **结构化日志**: 支持文本和JSON格式
- **完整覆盖**: HTTP请求、业务操作、错误、安全、审计等日志
- **性能优化**: 异步写入、敏感信息过滤、条件判断
- **运维友好**: 自动轮转、压缩、监控告警

通过统一的接口和中间件，开发者可以轻松在应用中添加日志记录功能，同时保持与现有系统的完全兼容性和代码的整洁性。
