# API鉴权中间件使用指南

## 概述

API鉴权中间件是本项目的核心安全组件，实现了基于JWT Token的身份认证和基于RBAC的权限控制。该中间件严格按照项目的分层架构设计规范开发，提供了完整的用户身份验证、权限检查和数据权限控制功能。

## 核心功能

### 1. JWT Token验证
- 支持标准的JWT Token格式验证
- 自动检查Token签名、过期时间
- 支持从Authorization头获取Bearer Token

### 2. 用户身份验证
- 自动从数据库验证用户状态
- 检查用户是否被禁用、锁定或删除
- 将用户信息设置到请求上下文中

### 3. 权限控制
- 基于RBAC模型的细粒度权限控制
- 支持API级别的权限验证
- 可配置的权限例外路由

### 4. 路由保护
- 支持不需要登录的公开路由
- 支持已登录但不需要权限验证的路由
- 灵活的路由配置机制

## 使用方法

### 1. 中间件注册

中间件已在HTTP服务启动时自动注册：

```go
// internal/router/api.go
func Api(ctx context.Context, group *ghttp.RouterGroup) {
    group.Group(simple.RouterPrefix(ctx, consts.AppApi), func(group *ghttp.RouterGroup) {
        group.Bind(
            api.NewRole(),
        )
        
        // API鉴权中间件在这里生效
        group.Middleware(service.Middleware().ApiAuth)
        
        group.Bind()
    })
}
```

### 2. 路由权限配置

在 `manifest/config/config.yaml` 中配置路由权限：

```yaml
router:
  api:
    prefix: "/api"
    # 不需要登录验证的路由
    exceptLogin:
      - "/login"
      - "/register" 
      - "/captcha"
      - "/ping"
      - "/health"
    # 不需要权限验证的路由（已登录但不验证具体权限）
    exceptAuth:
      - "/user/profile"
      - "/user/logout" 
      - "/user/refresh-token"
```

### 3. 客户端请求示例

#### 登录后获取Token
```javascript
const loginResponse = await fetch('/api/login', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        username: 'admin',
        password: 'password'
    })
});

const { accessToken } = await loginResponse.json();
```

#### 携带Token访问受保护的API
```javascript
const response = await fetch('/api/role/list', {
    method: 'GET',
    headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
    }
});
```

### 4. 获取当前用户信息

在控制器或业务逻辑中获取当前登录用户：

```go
// 在控制器中
func (c *Role) List(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error) {
    // 获取当前用户信息
    user := service.Middleware().GetCurrentUser(ctx)
    if user == nil {
        return nil, gerror.New("用户未登录")
    }
    
    g.Log().Infof(ctx, "当前用户: %s(%d)", user.Username, user.Id)
    
    // 业务逻辑处理...
}
```

## 权限设计

### 1. 权限标识生成规则

API权限标识按以下规则生成：
- 移除路由前缀：`/api/role/list` → `role/list`
- 将斜杠替换为冒号：`role/list` → `role:list`

### 2. 权限配置示例

在菜单管理中为API配置权限：

```sql
-- 角色管理相关权限
INSERT INTO menus (title, permission, type) VALUES 
('角色列表', 'role:list', 3),
('创建角色', 'role:create', 3),
('更新角色', 'role:update', 3),
('删除角色', 'role:delete', 3);
```

### 3. 数据权限控制

用户的数据权限范围通过角色的 `data_scope` 字段控制：

- `1`: 全部数据权限
- `2`: 部门数据权限  
- `3`: 部门及以下数据权限
- `4`: 仅本人数据权限
- `5`: 自定义数据权限

## 错误处理

### 1. 认证错误响应

```json
{
    "code": 401,
    "message": "访问令牌已过期，请重新登录",
    "timestamp": 1640995200,
    "traceID": "trace_id_here"
}
```

### 2. 权限错误响应

```json
{
    "code": 401, 
    "message": "您没有访问该接口的权限",
    "timestamp": 1640995200,
    "traceID": "trace_id_here"
}
```

## 安全特性

### 1. Token安全
- 使用HMAC-SHA256签名算法
- 支持Token过期时间控制
- 自动验证Token完整性

### 2. 用户状态检查
- 实时验证用户状态（正常/锁定/禁用）
- 支持软删除用户的访问拦截
- 记录认证失败日志

### 3. 请求日志
- 记录所有认证失败的请求
- 包含IP地址、User-Agent等信息
- 便于安全审计和异常排查

## 扩展配置

### 1. 自定义权限验证

可以在业务逻辑中添加额外的权限检查：

```go
// 检查用户是否有特定角色
func checkUserRole(ctx context.Context, roleCode string) error {
    user := service.Middleware().GetCurrentUser(ctx)
    if user == nil {
        return gerror.New("用户未登录")
    }
    
    hasRole, err := service.Role().CheckUserRole(ctx, user.Id, roleCode)
    if err != nil {
        return err
    }
    
    if !hasRole {
        return gerror.New("权限不足")
    }
    
    return nil
}
```

### 2. 数据权限过滤

在查询数据时应用数据权限过滤：

```go
func (s *sRole) GetRoleList(ctx context.Context, in *sysin.RoleListInp) (*sysout.RoleListModel, error) {
    user := service.Middleware().GetCurrentUser(ctx)
    
    // 根据用户数据权限范围过滤数据
    query := g.DB().Model("roles")
    
    // 获取用户数据权限范围
    dataScope, err := service.Role().GetUserDataScope(ctx, user.Id)
    if err != nil {
        return nil, err
    }
    
    switch dataScope {
    case entity.DataScopeDept:
        query = query.Where("dept_id = ?", user.DeptId)
    case entity.DataScopeSelf:
        query = query.Where("created_by = ?", user.Id)
    // 其他数据权限范围处理...
    }
    
    // 执行查询...
}
```

## 最佳实践

### 1. Token管理
- 建议Token过期时间设置为24小时
- 前端应实现Token自动刷新机制
- 敏感操作可要求重新验证

### 2. 权限粒度
- API权限应细化到具体的操作
- 避免过于粗粒度的权限设计
- 遵循最小权限原则

### 3. 安全审计
- 定期检查权限配置的合理性
- 监控异常的权限访问模式
- 及时清理无效的用户和权限

## 故障排查

### 1. 常见问题

**Token验证失败**
- 检查Token格式是否正确
- 验证签名密钥配置
- 确认Token是否过期

**权限验证失败**
- 检查用户是否分配了正确的角色
- 验证菜单权限配置
- 确认权限标识生成规则

**用户状态异常**
- 检查用户状态字段值
- 验证用户是否被软删除
- 确认角色权限是否过期

### 2. 调试技巧

启用详细日志记录：

```yaml
system:
  log:
    level: "debug"
    switch: true
```

查看认证失败日志：

```bash
tail -f logs/access.log | grep "API认证失败"
```

通过以上文档，开发人员可以完全了解和正确使用API鉴权中间件系统。 