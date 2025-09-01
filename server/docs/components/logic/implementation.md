# 业务逻辑层实现

## 概述

业务逻辑层位于 `internal/logic` 目录，是整个应用的核心，负责实现具体的业务逻辑。该层通过实现服务接口，为控制器层提供业务功能支持。

## 组织结构

```
internal/logic/
├── api/                # API业务逻辑
│   ├── user.go         # 用户管理业务逻辑
│   └── role.go         # 角色管理业务逻辑
├── middleware/         # 中间件业务逻辑
│   ├── api_auth.go     # API认证中间件
│   ├── cors.go         # 跨域处理中间件
│   ├── ctx.go          # 上下文处理中间件
│   ├── pre_filter.go   # 预过滤中间件
│   └── response.go     # 响应处理中间件
├── hook/               # 钩子函数
│   └── access_log.go   # 访问日志钩子
├── sys/                # 系统级业务逻辑
└── logic.go            # 逻辑层初始化
```

## 业务逻辑实现模式

### 1. 服务实现结构

```go
// internal/logic/api/role.go
package api

import (
    "client-app/internal/model/entity"
    "client-app/internal/model/input/sysin"
    "client-app/internal/model/output/sysout"
    "client-app/internal/service"
    "context"
    "fmt"
    "testing"
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/errors/gerror"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gtime"
    "github.com/gogf/gf/v2/util/gconv"
    "github.com/stretchr/testify/assert"
)

type sRole struct{}

func NewRole() *sRole {
    return &sRole{}
}

func init() {
    service.RegisterRole(NewRole())
}
```

### 2. 角色列表查询业务逻辑

```go
// GetRoleList 获取角色列表
func (s *sRole) GetRoleList(ctx context.Context, in *sysin.RoleListInp) (*sysout.RoleListModel, error) {
    // 1. 参数过滤和验证
    if err := in.Filter(ctx); err != nil {
        return nil, err
    }

    // 2. 构建查询条件
    db := g.DB().Model("roles").Where("deleted_at IS NULL")

    // 3. 添加搜索条件
    if in.Status >= 0 {
        db = db.Where("status = ?", in.Status)
    }
    if in.DataScope > 0 {
        db = db.Where("data_scope = ?", in.DataScope)
    }
    if in.Name != "" {
        db = db.WhereLike("name", "%"+in.Name+"%")
    }
    if in.Code != "" {
        db = db.WhereLike("code", "%"+in.Code+"%")
    }

    // 4. 查询总数
    totalCount, err := db.Count()
    if err != nil {
        return nil, gerror.Newf("查询角色总数失败: %v", err)
    }

    // 5. 如果没有数据，直接返回空结果
    if totalCount == 0 {
        return &sysout.RoleListModel{
            List:     []*sysout.RoleModel{},
            Total:    0,
            Page:     in.Page,
            PageSize: in.PageSize,
        }, nil
    }

    // 6. 分页查询
    offset := (in.Page - 1) * in.PageSize
    db = db.Offset(offset).Limit(in.PageSize)

    // 7. 排序
    orderBy := fmt.Sprintf("%s %s", in.OrderBy, in.OrderType)
    db = db.Order(orderBy)

    // 8. 执行查询
    var roles []*entity.Role
    if err := db.Scan(&roles); err != nil {
        return nil, gerror.Newf("查询角色列表失败: %v", err)
    }

    // 9. 转换为输出模型
    list := make([]*sysout.RoleModel, len(roles))
    for i, role := range roles {
        list[i] = sysout.ConvertToRoleModel(role)
    }

    return &sysout.RoleListModel{
        List:     list,
        Total:    int64(totalCount),
        Page:     in.Page,
        PageSize: in.PageSize,
    }, nil
}
```

### 3. 角色详情查询业务逻辑

```go
// GetRoleDetail 获取角色详情
func (s *sRole) GetRoleDetail(ctx context.Context, in *sysin.RoleDetailInp) (*sysout.RoleDetailModel, error) {
    // 1. 参数验证
    if err := in.Filter(ctx); err != nil {
        return nil, err
    }

    // 2. 查询角色信息
    var role *entity.Role
    err := g.DB().Model("roles").Where("id = ? AND deleted_at IS NULL", in.Id).Scan(&role)
    if err != nil {
        return nil, gerror.Newf("查询角色详情失败: %v", err)
    }
    if role == nil {
        return nil, gerror.New("角色不存在或已删除")
    }

    // 3. 查询角色拥有的菜单权限
    menuIds, err := s.getRoleMenuIds(ctx, in.Id)
    if err != nil {
        return nil, err
    }

    // 4. 查询权限标识列表
    permissions, err := s.getRolePermissions(ctx, in.Id)
    if err != nil {
        return nil, err
    }

    // 5. 转换为输出模型
    return sysout.ConvertToRoleDetailModel(role, menuIds, permissions), nil
}
```

### 4. 角色创建业务逻辑

```go
// CreateRole 创建角色
func (s *sRole) CreateRole(ctx context.Context, in *sysin.CreateRoleInp) (*sysout.RoleModel, error) {
    // 1. 参数验证
    if err := in.Filter(ctx); err != nil {
        return nil, err
    }

    // 2. 检查角色编码是否已存在
    exists, err := s.checkRoleCodeExists(ctx, in.Code, 0)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, gerror.New("角色编码已存在")
    }

    // 3. 检查角色名称是否已存在
    exists, err = s.checkRoleNameExists(ctx, in.Name, 0)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, gerror.New("角色名称已存在")
    }

    var resultRole *sysout.RoleModel

    // 4. 开启事务
    err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
        // 5. 构建角色数据
        roleData := &entity.Role{
            Name:        in.Name,
            Code:        in.Code,
            Description: in.Description,
            Status:      in.Status,
            Sort:        in.Sort,
            DataScope:   in.DataScope,
            Remark:      in.Remark,
            CreatedAt:   gtime.Now(),
            UpdatedAt:   gtime.Now(),
        }

        // 6. 获取当前用户ID
        if userId := s.getCurrentUserId(ctx); userId > 0 {
            roleData.CreatedBy = userId
            roleData.UpdatedBy = userId
        }

        // 7. 插入角色记录
        result, err := tx.Model("roles").Data(roleData).Insert()
        if err != nil {
            return gerror.Newf("创建角色失败: %v", err)
        }

        roleId, err := result.LastInsertId()
        if err != nil {
            return gerror.Newf("获取角色ID失败: %v", err)
        }

        roleData.Id = roleId

        // 8. 分配菜单权限
        if len(in.MenuIds) > 0 {
            if err := s.assignRoleMenus(ctx, tx, roleId, in.MenuIds); err != nil {
                return err
            }
        }

        // 9. 转换为输出模型
        resultRole = sysout.ConvertToRoleModel(roleData)
        return nil
    })

    if err != nil {
        return nil, err
    }

    return resultRole, nil
}
```

### 5. 角色权限管理业务逻辑

```go
// GetRoleMenus 获取角色菜单权限
func (s *sRole) GetRoleMenus(ctx context.Context, in *sysin.RoleMenuInp) (*sysout.RoleMenuModel, error) {
    // 1. 参数验证
    if err := in.Filter(ctx); err != nil {
        return nil, err
    }

    // 2. 获取角色菜单ID列表
    menuIds, err := s.getRoleMenuIds(ctx, in.RoleId)
    if err != nil {
        return nil, err
    }

    // 3. 查询菜单详情
    var menus []*entity.Menu
    if len(menuIds) > 0 {
        err = g.DB().Model("menus").WhereIn("id", menuIds).Where("deleted_at IS NULL").Scan(&menus)
        if err != nil {
            return nil, gerror.Newf("查询菜单信息失败: %v", err)
        }
    }

    // 4. 转换为输出模型
    return &sysout.RoleMenuModel{
        RoleId:  in.RoleId,
        MenuIds: menuIds,
        Menus:   menus,
    }, nil
}

// UpdateRoleMenus 更新角色菜单权限
func (s *sRole) UpdateRoleMenus(ctx context.Context, in *sysin.UpdateRoleMenuInp) error {
    // 1. 参数验证
    if err := in.Filter(ctx); err != nil {
        return err
    }

    // 2. 检查角色是否存在
    exists, err := s.checkRoleExists(ctx, in.RoleId)
    if err != nil {
        return err
    }
    if !exists {
        return gerror.New("角色不存在")
    }

    // 3. 开启事务更新菜单权限
    return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
        return s.updateRoleMenus(ctx, tx, in.RoleId, in.MenuIds)
    })
}
```

### 6. 权限验证业务逻辑

```go
// CheckUserPermission 检查用户权限
func (s *sRole) CheckUserPermission(ctx context.Context, userId int64, permission string) (bool, error) {
    // 1. 获取用户权限列表
    permissions, err := s.GetUserPermissions(ctx, userId)
    if err != nil {
        return false, err
    }

    // 2. 检查权限是否存在
    for _, p := range permissions {
        if p == permission {
            return true, nil
        }
    }

    return false, nil
}

// GetUserPermissions 获取用户权限列表
func (s *sRole) GetUserPermissions(ctx context.Context, userId int64) ([]string, error) {
    // 1. 查询用户角色
    var userRoles []*entity.UserRole
    err := g.DB().Model("user_roles").Where("user_id = ? AND deleted_at IS NULL", userId).Scan(&userRoles)
    if err != nil {
        return nil, gerror.Newf("查询用户角色失败: %v", err)
    }

    if len(userRoles) == 0 {
        return []string{}, nil
    }

    // 2. 提取角色ID列表
    roleIds := make([]int64, len(userRoles))
    for i, ur := range userRoles {
        roleIds[i] = ur.RoleId
    }

    // 3. 查询角色权限
    var permissions []string
    query := `
        SELECT DISTINCT m.permission
        FROM role_menus rm
        JOIN menus m ON rm.menu_id = m.id
        WHERE rm.role_id IN (?) AND rm.deleted_at IS NULL AND m.deleted_at IS NULL AND m.permission != ''
    `
    err = g.DB().Raw(query, roleIds).Scan(&permissions)
    if err != nil {
        return nil, gerror.Newf("查询角色权限失败: %v", err)
    }

    return permissions, nil
}
```

### 7. 辅助方法实现

```go
// getRoleMenuIds 获取角色菜单ID列表
func (s *sRole) getRoleMenuIds(ctx context.Context, roleId int64) ([]int64, error) {
    var menuIds []int64
    err := g.DB().Model("role_menus").
        Fields("menu_id").
        Where("role_id = ? AND deleted_at IS NULL", roleId).
        Scan(&menuIds)
    if err != nil {
        return nil, gerror.Newf("查询角色菜单权限失败: %v", err)
    }
    return menuIds, nil
}

// checkRoleCodeExists 检查角色编码是否存在
func (s *sRole) checkRoleCodeExists(ctx context.Context, code string, excludeId int64) (bool, error) {
    db := g.DB().Model("roles").Where("code = ? AND deleted_at IS NULL", code)
    if excludeId > 0 {
        db = db.Where("id != ?", excludeId)
    }
    
    count, err := db.Count()
    if err != nil {
        return false, gerror.Newf("检查角色编码失败: %v", err)
    }
    return count > 0, nil
}

// checkRoleNameExists 检查角色名称是否存在
func (s *sRole) checkRoleNameExists(ctx context.Context, name string, excludeId int64) (bool, error) {
    db := g.DB().Model("roles").Where("name = ? AND deleted_at IS NULL", name)
    if excludeId > 0 {
        db = db.Where("id != ?", excludeId)
    }
    
    count, err := db.Count()
    if err != nil {
        return false, gerror.Newf("检查角色名称失败: %v", err)
    }
    return count > 0, nil
}

// assignRoleMenus 分配角色菜单权限
func (s *sRole) assignRoleMenus(ctx context.Context, tx gdb.TX, roleId int64, menuIds []int64) error {
    // 1. 构建菜单权限数据
    menuData := make([]map[string]interface{}, len(menuIds))
    for i, menuId := range menuIds {
        menuData[i] = map[string]interface{}{
            "role_id":    roleId,
            "menu_id":    menuId,
            "created_at": gtime.Now(),
            "updated_at": gtime.Now(),
        }
    }

    // 2. 批量插入菜单权限
    _, err := tx.Model("role_menus").Data(menuData).Insert()
    if err != nil {
        return gerror.Newf("分配菜单权限失败: %v", err)
    }

    return nil
}

// getCurrentUserId 获取当前用户ID
func (s *sRole) getCurrentUserId(ctx context.Context) int64 {
    // 从上下文获取用户ID（具体实现可能因项目而异）
    if userId := ctx.Value("user_id"); userId != nil {
        return gconv.Int64(userId)
    }
    return 0
}
```

## 业务逻辑实现特点

### 1. 事务管理

对于涉及多表操作的业务逻辑，使用数据库事务确保数据一致性：

```go
err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
    // 业务逻辑操作
    if err := s.doSomething(ctx, tx); err != nil {
        return err // 自动回滚
    }
    return nil // 提交事务
})
```

### 2. 参数验证

每个业务方法都进行参数验证，确保数据有效性：

```go
func (s *sRole) CreateRole(ctx context.Context, in *sysin.CreateRoleInp) (*sysout.RoleModel, error) {
    // 参数验证
    if err := in.Filter(ctx); err != nil {
        return nil, err
    }
    
    // 业务逻辑验证
    if exists, err := s.checkRoleCodeExists(ctx, in.Code, 0); err != nil {
        return nil, err
    } else if exists {
        return nil, gerror.New("角色编码已存在")
    }
    
    // 继续处理...
}
```

### 3. 错误处理

使用GoFrame的错误处理机制，提供清晰的错误信息：

```go
if err != nil {
    return nil, gerror.Newf("查询角色列表失败: %v", err)
}
```

### 4. 数据转换

在业务逻辑层进行数据模型转换，隔离不同层之间的数据结构：

```go
// 实体转换为输出模型
list := make([]*sysout.RoleModel, len(roles))
for i, role := range roles {
    list[i] = sysout.ConvertToRoleModel(role)
}
```

### 5. 分页处理

标准的分页逻辑实现：

```go
// 计算偏移量
offset := (in.Page - 1) * in.PageSize
db = db.Offset(offset).Limit(in.PageSize)

// 返回分页结果
return &sysout.RoleListModel{
    List:     list,
    Total:    int64(totalCount),
    Page:     in.Page,
    PageSize: in.PageSize,
}, nil
```

### 6. 权限控制

实现基于角色的权限控制：

```go
// 检查用户权限
func (s *sRole) CheckUserPermission(ctx context.Context, userId int64, permission string) (bool, error) {
    permissions, err := s.GetUserPermissions(ctx, userId)
    if err != nil {
        return false, err
    }
    
    for _, p := range permissions {
        if p == permission {
            return true, nil
        }
    }
    
    return false, nil
}
```

## 测试支持

业务逻辑层支持单元测试：

```go
func TestCreateRole(t *testing.T) {
    // 创建服务实例
    roleService := NewRole()
    
    // 准备测试数据
    input := &sysin.CreateRoleInp{
        Name:        "测试角色",
        Code:        "test_role",
        Description: "测试角色描述",
        Status:      1,
    }
    
    // 执行测试
    result, err := roleService.CreateRole(ctx, input)
    
    // 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "测试角色", result.Name)
}
```

## 最佳实践

1. **单一职责**：每个方法专注于一个业务功能
2. **事务管理**：合理使用数据库事务
3. **参数验证**：充分验证输入参数
4. **错误处理**：提供清晰的错误信息
5. **数据转换**：隔离不同层的数据结构
6. **代码复用**：抽取公共的辅助方法
7. **性能优化**：合理使用数据库查询和缓存
8. **测试友好**：编写可测试的业务逻辑

## 总结

业务逻辑层在本项目中承担着以下职责：

1. **业务规则实现**：实现具体的业务规则和流程
2. **数据处理**：处理数据的增删改查操作
3. **权限控制**：实现基于角色的权限控制
4. **事务管理**：保证数据操作的一致性
5. **错误处理**：提供清晰的错误信息
6. **数据转换**：在不同数据模型之间进行转换

通过这种设计，业务逻辑层实现了高内聚低耦合，便于维护和扩展。
