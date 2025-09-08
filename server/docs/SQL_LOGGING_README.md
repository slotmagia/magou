# SQL日志记录功能使用指南

## 功能概述

本功能为您的GoFrame应用程序提供了增强的SQL日志记录功能，可以在控制台中显示SQL查询的生成函数，并支持点击跳转到对应的代码位置。

## 主要特性

1. **自动记录SQL查询**：自动记录所有数据库查询的SQL语句、参数和执行时间
2. **调用者信息**：显示SQL查询的调用者信息，包括文件名、行号和函数名
3. **IDE支持**：支持IDE点击跳转，格式为 `file:line:function`
4. **性能监控**：根据执行时间显示不同颜色（绿色=快速，黄色=中等，红色=慢查询）
5. **链路追踪**：集成GoFrame的链路追踪ID
6. **美化输出**：格式化的SQL语句和参数显示

## 安装和配置

### 1. 启用SQL日志功能

在您的应用程序初始化时调用：

```go
import "client-app/utility/logger"

// 启用SQL演示功能
logger.EnableSQLDemo()
```

### 2. 在中间件中初始化

在 `internal/logic/middleware/init.go` 中已经集成了初始化代码：

```go
func Init(ctx context.Context) error {
    // 启用SQL追踪
    logger.EnableSQLTracing()
    
    g.Log().Info(ctx, "所有中间件初始化完成")
    return nil
}
```

## 使用方法

### 方法1：在现有服务中集成

在您的服务方法中添加SQL日志记录：

```go
func (s *MenuService) GetMenuList(ctx context.Context, status int) ([]map[string]interface{}, error) {
    start := time.Now()
    
    // 执行数据库查询
    var menus []map[string]interface{}
    err := g.DB().Model("sys_menus").
        Where("status", status).
        Order("sort_order ASC").
        Scan(&menus)
    
    // 计算执行时间
    duration := time.Since(start)
    
    // 记录SQL日志（带调用者信息）
    logger.LogSQLWithCaller(ctx, "SELECT * FROM sys_menus WHERE status = ?", []interface{}{status}, duration)
    
    return menus, err
}
```

### 方法2：使用SQL包装器

使用提供的SQL包装器自动记录日志：

```go
import "client-app/utility/database"

// 使用包装器
db := database.DB()
var menus []map[string]interface{}
err := db.Model("sys_menus").
    Where("status", 1).
    Scan(&menus)
```

### 方法3：直接使用SQL追踪器

```go
import "client-app/utility/logger"

// 直接记录SQL日志
logger.TraceSQL(ctx, "SELECT * FROM sys_menus WHERE id = ?", []interface{}{1}, 25*time.Millisecond)
```

## 日志输出格式

### 控制台输出示例

```
[17:37:53.789] {7883d51ef7739972e383bdd003d33aea} [25ms] [menu.go:32:GetMenuList] SQL: SELECT
  `id`,`tenant_id`,`parent_id`,`menu_code`,`icon`,`path`,`component`,`permission`,`sort_order`,`visible`,`status`,`created_at`,`updated_at`,`redirect`,`active_menu`,`always_show`,`breadcrumb`,`remark`,`title`,`menu_type`
FROM `sys_menus`
WHERE (`id`=12) AND `deleted_at` IS NULL
LIMIT 1
ARGS: [0]=12
```

### 输出说明

- `[17:37:53.789]`：时间戳
- `{7883d51ef7739972e383bdd003d33aea}`：链路追踪ID
- `[25ms]`：执行时间（带颜色标识）
- `[menu.go:32:GetMenuList]`：调用者信息（支持IDE点击跳转）
- `SQL:`：格式化的SQL语句
- `ARGS:`：查询参数

### 颜色标识

- 🟢 **绿色**：执行时间 < 10ms（快速查询）
- 🟡 **黄色**：执行时间 10ms - 100ms（中等查询）
- 🔴 **红色**：执行时间 > 100ms（慢查询）

## IDE支持

### IntelliJ IDEA / GoLand

1. 打开 **Settings/Preferences**
2. 导航到 **Editor** > **General** > **Console**
3. 确保 **Hyperlink file paths** 选项已启用
4. 日志中的 `file:line:function` 格式将自动变为可点击链接

### VS Code

1. 安装 **Go** 扩展
2. 确保启用了 **Go: Use Language Server** 设置
3. 日志中的文件路径将自动变为可点击链接

## 配置选项

### 启用/禁用功能

```go
// 启用SQL演示功能
logger.EnableSQLDemo()

// 禁用SQL演示功能
logger.DisableSQLDemo()

// 启用SQL追踪功能
logger.EnableSQLTracing()

// 禁用SQL追踪功能
logger.DisableSQLTracing()
```

### 自定义配置

```go
// 获取默认SQL演示工具
demo := logger.GetDefaultSQLDemo()

// 自定义配置
demo.SetColors(true)  // 启用颜色输出
demo.Enable()         // 启用功能
```

## 示例代码

查看 `example_menu_with_sql_logging.go` 文件获取完整的使用示例。

## 注意事项

1. **性能影响**：SQL日志记录会有轻微的性能开销，建议仅在开发环境启用
2. **日志级别**：确保日志级别设置为 `debug` 以查看SQL日志
3. **敏感信息**：注意不要在日志中记录敏感信息（如密码、密钥等）
4. **生产环境**：生产环境建议禁用详细的SQL日志记录

## 故障排除

### 问题1：看不到SQL日志

**解决方案**：
1. 检查日志级别是否设置为 `debug`
2. 确认已调用 `logger.EnableSQLDemo()` 或 `logger.EnableSQLTracing()`
3. 检查配置文件中的日志设置

### 问题2：调用者信息显示为 "unknown"

**解决方案**：
1. 确保代码编译时包含了调试信息
2. 检查 `runtime.Caller` 的调用栈深度设置

### 问题3：IDE点击跳转不工作

**解决方案**：
1. 确认IDE支持文件路径超链接功能
2. 检查日志格式是否为 `file:line:function`
3. 确保文件路径是相对路径或绝对路径

## 扩展功能

### 添加自定义字段

```go
// 在SQL日志中添加自定义信息
func (s *MenuService) GetMenuList(ctx context.Context, status int) ([]map[string]interface{}, error) {
    start := time.Now()
    
    // 执行查询...
    
    duration := time.Since(start)
    
    // 添加自定义信息到上下文
    ctx = context.WithValue(ctx, "operation", "GetMenuList")
    ctx = context.WithValue(ctx, "user_id", "123")
    
    logger.LogSQLWithCaller(ctx, sql, args, duration)
    
    return result, nil
}
```

### 集成到现有中间件

```go
// 在现有的数据库中间件中添加SQL日志
func (m *DatabaseMiddleware) LogQuery(ctx context.Context, sql string, args []interface{}, duration time.Duration) {
    logger.LogSQLWithCaller(ctx, sql, args, duration)
}
```

## 总结

这个SQL日志记录功能为您的GoFrame应用程序提供了强大的调试和监控能力。通过显示SQL查询的调用者信息和执行时间，您可以：

1. 快速定位SQL查询的来源
2. 识别性能瓶颈
3. 监控数据库查询模式
4. 提高开发效率

使用IDE的点击跳转功能，您可以快速从日志跳转到生成SQL的代码位置，大大提高了调试效率。
