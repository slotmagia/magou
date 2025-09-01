# 多租户系统实施指南

## 概述

本文档详细说明了如何在现有的角色权限管理系统基础上实现多租户架构。多租户系统允许多个客户（租户）共享同一个应用实例，同时保持数据隔离和安全性。

## 一、实施方案概览

### 1.1 多租户模式选择

我们采用**共享数据库、共享数据表**的模式，通过 `tenant_id` 字段实现数据隔离：

- **优势**：资源利用率高、维护成本低、扩展性好
- **劣势**：需要严格的数据隔离控制
- **适用场景**：中小型SaaS应用、客户数量较多但单个客户数据量不大

### 1.2 架构调整

```
原有架构：用户 ← 角色 ← 菜单
多租户架构：租户 → 用户 ← 角色 ← 菜单
```

每个租户拥有独立的：
- 用户体系
- 角色权限体系
- 菜单体系
- 业务数据

## 二、数据库结构调整

### 2.1 新增租户表

```sql
-- 执行 internal/sql/tenants.sql
-- 创建租户管理表，包含租户基本信息、配置、限制等
```

### 2.2 现有表结构调整

```sql
-- 执行 internal/sql/add_tenant_id.sql
-- 为所有业务表添加 tenant_id 字段
-- 调整唯一索引以支持租户隔离
-- 添加租户相关的索引优化
```

关键变更：
- 所有业务表添加 `tenant_id` 字段
- 唯一索引调整为 `(tenant_id, 原字段)` 组合
- 外键关系保持租户一致性

### 2.3 数据迁移

```sql
-- 将现有数据迁移到系统租户（tenant_id = 1）
UPDATE users SET tenant_id = 1 WHERE tenant_id IS NULL;
UPDATE roles SET tenant_id = 1 WHERE tenant_id IS NULL;
UPDATE menus SET tenant_id = 1 WHERE tenant_id IS NULL;
UPDATE user_roles SET tenant_id = 1 WHERE tenant_id IS NULL;
UPDATE role_menus SET tenant_id = 1 WHERE tenant_id IS NULL;
```

## 三、业务逻辑调整

### 3.1 租户服务层

新增租户管理服务：
- `internal/service/tenant.go` - 租户服务接口
- `internal/logic/api/tenant.go` - 租户业务逻辑实现

主要功能：
- 租户CRUD操作
- 租户状态管理
- 租户统计信息
- 租户配置管理

### 3.2 多租户中间件

新增租户过滤中间件：
- `internal/logic/middleware/tenant.go` - 租户过滤和权限验证

功能包括：
- 租户识别（Header、域名、参数）
- 租户状态验证
- 租户权限控制
- 数据隔离保证

### 3.3 用户认证调整

扩展用户认证逻辑：
- `internal/logic/api/user_tenant_patch.go` - 多租户用户认证

调整内容：
- JWT Token 包含租户信息
- 登录验证包含租户过滤
- 权限获取包含租户隔离

### 3.4 数据模型调整

- `internal/model/entity/tenants.go` - 租户实体模型
- `internal/model/input/sysin/tenant.go` - 租户输入模型  
- `internal/model/output/sysout/tenant.go` - 租户输出模型
- `internal/model/context.go` - 上下文模型调整

## 四、API 接口调整

### 4.1 租户管理接口

新增租户管理相关接口：
- `GET /api/tenant/list` - 获取租户列表
- `POST /api/tenant/create` - 创建租户
- `PUT /api/tenant/update` - 更新租户
- `DELETE /api/tenant/delete` - 删除租户
- `GET /api/tenant/detail` - 获取租户详情
- `PUT /api/tenant/status` - 更新租户状态
- `GET /api/tenant/stats` - 获取租户统计
- `PUT /api/tenant/config` - 更新租户配置
- `GET /api/tenant/options` - 获取租户选项

### 4.2 现有接口调整

所有业务接口需要：
1. 添加租户过滤中间件
2. 查询条件包含 `tenant_id` 过滤
3. 创建数据时自动设置 `tenant_id`

## 五、实施步骤

### 步骤 1：数据库调整

```bash
# 1. 创建租户表
mysql -u root -p admin < internal/sql/tenants.sql

# 2. 调整现有表结构
mysql -u root -p admin < internal/sql/add_tenant_id.sql

# 3. 验证表结构
mysql -u root -p admin -e "DESCRIBE tenants;"
mysql -u root -p admin -e "SHOW INDEX FROM users;"
```

### 步骤 2：业务逻辑部署

```bash
# 1. 部署新增文件
# - internal/model/entity/tenants.go
# - internal/model/input/sysin/tenant.go
# - internal/model/output/sysout/tenant.go
# - internal/service/tenant.go
# - internal/logic/api/tenant.go
# - internal/logic/api/user_tenant_patch.go
# - internal/logic/middleware/tenant.go
# - internal/controller/api/tenant.go
# - internal/api/v1/tenant/tenant.go

# 2. 修改现有文件
# - internal/model/context.go
# - utility/simple/simple.go

# 3. 重新编译应用
go build -o admin main.go
```

### 步骤 3：路由配置调整

修改 `internal/router/api.go`，添加租户路由：

```go
func Api(ctx context.Context, group *ghttp.RouterGroup) {
    group.Group(simple.RouterPrefix(ctx, consts.AppApi), func(group *ghttp.RouterGroup) {
        // 租户过滤中间件（所有接口生效）
        group.Middleware(service.Middleware().TenantFilter)
        
        // 绑定控制器
        group.Bind(
            api.NewUser(),
            api.NewRole(),
            api.NewMenu(),
            api.NewTenant(), // 新增租户控制器
        )
        
        // 认证中间件
        group.Middleware(service.Middleware().ApiAuth)
        // 租户权限验证中间件
        group.Middleware(service.Middleware().TenantAuth)
    })
}
```

### 步骤 4：中间件服务注册

修改 `internal/service/middleware.go`，添加租户中间件：

```go
type IMiddleware interface {
    // ... 现有方法
    TenantFilter(r *ghttp.Request)   // 租户过滤
    TenantAuth(r *ghttp.Request)     // 租户权限验证
}
```

### 步骤 5：现有业务逻辑调整

需要调整的文件示例：

1. **用户管理** (`internal/logic/api/user.go`)
```go
// 原有查询
m := g.DB().Model("users").Where("deleted_at IS NULL")

// 调整后查询
tenantId := middleware.GetCurrentTenantId(r)
m := g.DB().Model("users").Where("tenant_id = ? AND deleted_at IS NULL", tenantId)
```

2. **角色管理** (`internal/logic/api/role.go`)
```go
// 类似调整，所有查询添加租户过滤
```

3. **菜单管理** (`internal/logic/api/menu.go`)
```go
// 类似调整，所有查询添加租户过滤
```

### 步骤 6：测试验证

1. **基础功能测试**
```bash
# 创建租户
curl -X POST "http://localhost:8888/api/tenant/create" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试租户",
    "code": "test",
    "maxUsers": 100,
    "storageLimit": 1073741824,
    "adminName": "admin",
    "adminEmail": "admin@test.com",
    "adminPassword": "MTIzNDU2", 
    "remark": "测试租户"
  }'

# 租户用户登录
curl -X POST "http://localhost:8888/api/user/login" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Id: 2" \
  -d '{
    "username": "admin",
    "password": "MTIzNDU2",
    "captchaId": "xxx",
    "captcha": "1234"
  }'
```

2. **数据隔离测试**
- 验证不同租户用户无法看到其他租户数据
- 验证租户管理员权限范围
- 验证系统管理员可以管理所有租户

## 六、配置说明

### 6.1 租户识别配置

支持多种租户识别方式：

1. **HTTP Header**
```
X-Tenant-Id: 123
```

2. **域名映射**
```
tenant1.example.com -> tenant_id: 1
tenant2.example.com -> tenant_id: 2
```

3. **URL参数**（开发环境）
```
http://localhost:8888/api/user/list?tenant_id=1
```

### 6.2 租户配置管理

租户配置通过 JSON 格式存储：

```json
{
  "features": {
    "advancedReports": true,
    "apiAccess": false,
    "customBranding": true
  },
  "limitations": {
    "maxApiCalls": 10000,
    "maxStorage": 5368709120
  },
  "settings": {
    "theme": "dark",
    "language": "zh-CN",
    "timezone": "Asia/Shanghai"
  }
}
```

## 七、性能优化

### 7.1 数据库优化

1. **索引优化**
```sql
-- 复合索引优化查询性能
CREATE INDEX idx_tenant_status ON users(tenant_id, status);
CREATE INDEX idx_tenant_created ON users(tenant_id, created_at);
```

2. **分区表**（可选）
```sql
-- 大数据量情况下可考虑按租户分区
ALTER TABLE users PARTITION BY HASH(tenant_id) PARTITIONS 16;
```

### 7.2 缓存策略

1. **租户信息缓存**
```go
// 缓存租户基础信息
cacheKey := fmt.Sprintf("tenant:info:%d", tenantId)
gcache.Set(ctx, cacheKey, tenant, 24*time.Hour)
```

2. **权限信息缓存**
```go
// 缓存用户权限信息
cacheKey := fmt.Sprintf("user:permissions:%d:%d", tenantId, userId)
gcache.Set(ctx, cacheKey, permissions, 1*time.Hour)
```

## 八、安全考虑

### 8.1 数据隔离

1. **强制租户过滤**
- 所有数据库查询必须包含 `tenant_id` 条件
- 通过中间件自动添加租户过滤
- 代码审查确保没有遗漏

2. **权限验证**
- 验证用户是否属于当前租户
- 系统管理员可以跨租户访问
- 租户管理员只能管理本租户

### 8.2 接口安全

1. **租户验证**
- 验证租户状态（正常、锁定、禁用）
- 验证租户是否过期
- 验证用户权限范围

2. **数据校验**
- 创建数据时验证租户限制
- 检查存储空间限制
- 检查用户数量限制

## 九、监控和运维

### 9.1 租户监控

1. **资源使用监控**
```go
// 监控租户资源使用情况
type TenantMetrics struct {
    TenantId     uint64 `json:"tenantId"`
    UserCount    int    `json:"userCount"`
    StorageUsed  int64  `json:"storageUsed"`
    ApiCalls     int    `json:"apiCalls"`
    LastActive   time.Time `json:"lastActive"`
}
```

2. **告警配置**
- 租户资源使用超限告警
- 租户过期时间提醒
- 异常访问模式检测

### 9.2 数据备份

1. **租户数据备份**
```bash
# 按租户备份数据
mysqldump --where="tenant_id=1" admin users roles menus > tenant_1_backup.sql
```

2. **数据恢复**
```bash
# 租户数据恢复
mysql admin < tenant_1_backup.sql
```

## 十、故障排除

### 10.1 常见问题

1. **租户数据泄露**
- 检查查询条件是否包含 `tenant_id`
- 检查中间件是否正确配置
- 检查用户权限验证逻辑

2. **性能问题**
- 检查数据库索引是否正确
- 检查查询是否高效
- 检查缓存是否生效

3. **租户创建失败**
- 检查租户编码是否重复
- 检查管理员用户名是否重复
- 检查数据库事务是否正确

### 10.2 调试工具

1. **日志配置**
```go
// 添加租户信息到日志上下文
g.Log().Debugf(ctx, "Tenant %d: Processing request", tenantId)
```

2. **SQL调试**
```yaml
# config.yaml
database:
  default:
    debug: true  # 开启SQL调试
```

## 十一、升级和维护

### 11.1 版本升级

1. **数据库升级脚本**
```sql
-- v1.1.0 升级脚本
ALTER TABLE tenants ADD COLUMN new_field varchar(100) DEFAULT NULL;
UPDATE tenants SET new_field = 'default_value';
```

2. **兼容性处理**
- 保持向后兼容
- 渐进式升级
- 回滚方案准备

### 11.2 维护操作

1. **租户数据清理**
```sql
-- 清理已删除租户的数据
DELETE FROM users WHERE tenant_id IN (
    SELECT id FROM tenants WHERE deleted_at IS NOT NULL
);
```

2. **性能优化**
```sql
-- 定期分析表结构
ANALYZE TABLE users, roles, menus;

-- 优化查询计划
EXPLAIN SELECT * FROM users WHERE tenant_id = 1 AND status = 1;
```

## 总结

本多租户实施方案采用共享数据库模式，通过 `tenant_id` 字段实现数据隔离，具有以下特点：

**优势：**
- 资源利用率高
- 维护成本低  
- 快速部署新租户
- 统一的功能更新

**注意事项：**
- 严格的数据隔离控制
- 完善的权限验证机制
- 充分的性能优化
- 全面的监控和告警

通过以上实施方案，可以在现有系统基础上快速实现多租户架构，满足SaaS应用的业务需求。
