# 开发规范

## 代码风格

### 命名规范

1. **包名**: 使用小写单词，不使用下划线或混合大小写

   - 正确: `controller`, `service`, `model`
   - 错误: `Controller`, `service_api`

2. **文件名**: 使用小写单词，用下划线分隔

   - 正确: `payment.go`, `api_auth.go`
   - 错误: `Payment.go`, `ApiAuth.go`

3. **结构体名**: 使用驼峰命名法

   - 普通结构体: `PaymentReq`, `UserInfo`
   - 接口: 以`I`开头，如`IPayment`, `IMiddleware`
   - 接口实现: 以小写`s`开头，如`sPayment`, `sMiddleware`

4. **方法名**: 使用驼峰命名法，首字母大写表示公开，小写表示私有

   - 公开方法: `Submit()`, `GetUser()`
   - 私有方法: `validateToken()`, `hashPassword()`

5. **变量名**: 使用驼峰命名法，首字母小写

   - 局部变量: `userID`, `orderInfo`
   - 全局变量: `localPayment`, `defaultConfig`

6. **常量名**: 全部大写，使用下划线分隔
   - 正确: `APP_NAME`, `MAX_RETRY_COUNT`
   - 错误: `AppName`, `maxRetryCount`

### 代码格式

1. 使用`gofmt`或`goimports`格式化代码
2. 缩进使用 Tab，宽度为 4 个空格
3. 行长度尽量控制在 120 字符以内
4. 导入包分组按标准库、第三方库、项目内包顺序排列

## 注释规范

1. **包注释**: 每个包都应有包注释，位于`package`声明之前

   ```go
   // Package router 提供路由注册和管理功能
   package router
   ```

2. **函数/方法注释**: 使用完整的句子描述功能

   ```go
   // GetUser 根据用户ID获取用户信息
   // 参数:
   //   - ctx: 上下文
   //   - id: 用户ID
   // 返回:
   //   - *model.User: 用户信息
   //   - error: 错误信息
   func GetUser(ctx context.Context, id int64) (*model.User, error) {
   ```

3. **常量和变量注释**: 对非自明的常量和变量添加注释
   ```go
   // MaxRetryCount 定义最大重试次数
   const MaxRetryCount = 3
   ```

## 错误处理规范

1. **错误检查**: 每个可能返回错误的函数调用都必须检查错误

   ```go
   result, err := service.DoSomething()
   if err != nil {
       return nil, err
   }
   ```

2. **错误封装**: 使用`gerror`封装错误，保留上下文

   ```go
   if err := db.QueryRow(); err != nil {
       return nil, gerror.Wrap(err, "查询数据失败")
   }
   ```

3. **错误返回**: 尽量返回有意义的错误，避免简单的`errors.New("error")`

## 接口设计规范

1. **RESTful 原则**: 遵循 RESTful API 设计原则

   - 使用 HTTP 方法表示操作: GET(查询), POST(创建), PUT(更新), DELETE(删除)
   - URL 使用名词表示资源，如`/api/payment`而非`/api/doPayment`

2. **接口版本化**: 使用路径版本化策略

   - 接口定义放在`internal/api/api/*/v1/`目录下
   - 新版本接口创建新目录，如`v2`、`v3`

3. **参数验证**: 使用结构体标签进行验证
   ```go
   type LoginReq struct {
       Username string `json:"username" v:"required#用户名不能为空"`
       Password string `json:"password" v:"required#密码不能为空"`
   }
   ```

## 测试规范

1. **单元测试**: 关键功能必须编写单元测试

   - 测试文件命名为`*_test.go`
   - 测试函数命名为`TestXxx`

2. **测试覆盖率**: 核心模块测试覆盖率不低于 70%

3. **模拟与桩**: 使用模拟对象测试依赖外部服务的代码

## 日志规范

1. **日志级别**: 合理使用日志级别

   - Debug: 调试信息，仅在开发环境使用
   - Info: 一般信息，记录关键操作
   - Warning: 警告信息，可能出现问题
   - Error: 错误信息，影响功能但系统可继续运行
   - Fatal: 致命错误，导致系统无法继续运行

2. **日志内容**: 记录必要的上下文信息

   - 请求 ID/跟踪 ID
   - 用户信息(脱敏)
   - 操作描述
   - 错误详情

3. **日志格式**: 使用结构化日志，便于解析和分析

## 安全规范

1. **身份验证**: 所有 API 除特殊说明外必须进行身份验证
2. **输入验证**: 所有用户输入必须验证和净化
3. **敏感信息**: 敏感信息不得明文存储或记录在日志中
4. **HTTPS**: 生产环境必须使用 HTTPS

## 性能规范

1. **数据库操作**: 避免 N+1 查询问题，大批量操作使用批处理
2. **缓存使用**: 合理使用缓存减少重复计算和数据库访问
3. **并发控制**: 使用适当的锁机制或并发原语确保线程安全
4. **资源释放**: 确保所有资源(文件、连接等)正确关闭
