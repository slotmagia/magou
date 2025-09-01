# 服务接口定义

## 服务接口概述

服务接口是连接控制器层和逻辑层的桥梁，定义了业务功能的抽象接口。项目中的服务接口定义在`internal/service`目录下，采用接口+注册机制的设计模式。

## 接口定义规范

服务接口命名以`I`开头，表示接口，包含完整的业务功能定义：

```go
// internal/service/role.go
type IRole interface {
    // 角色基础操作
    GetRoleList(ctx context.Context, in *sysin.RoleListInp) (res *sysout.RoleListModel, err error)
    GetRoleDetail(ctx context.Context, in *sysin.RoleDetailInp) (res *sysout.RoleDetailModel, err error)
    CreateRole(ctx context.Context, in *sysin.CreateRoleInp) (res *sysout.RoleModel, err error)
    UpdateRole(ctx context.Context, in *sysin.UpdateRoleInp) (res *sysout.RoleModel, err error)
    DeleteRole(ctx context.Context, in *sysin.DeleteRoleInp) (err error)
    BatchDeleteRole(ctx context.Context, in *sysin.BatchDeleteRoleInp) (err error)
    UpdateRoleStatus(ctx context.Context, in *sysin.UpdateRoleStatusInp) (err error)
    CopyRole(ctx context.Context, in *sysin.CopyRoleInp) (res *sysout.RoleModel, err error)

    // 角色权限管理
    GetRoleMenus(ctx context.Context, in *sysin.RoleMenuInp) (res *sysout.RoleMenuModel, err error)
    UpdateRoleMenus(ctx context.Context, in *sysin.UpdateRoleMenuInp) (err error)
    GetRolePermissions(ctx context.Context, in *sysin.RolePermissionInp) (res *sysout.RolePermissionModel, err error)

    // 角色选项和统计
    GetRoleOptions(ctx context.Context, in *sysin.RoleOptionInp) (res []*sysout.RoleOptionModel, err error)
    GetRoleStats(ctx context.Context) (res *sysout.RoleStatsModel, err error)
    GetDataScopeOptions(ctx context.Context) (res []*sysout.DataScopeModel, err error)

    // 用户角色关联操作
    AssignUserRoles(ctx context.Context, userId int64, roleIds []int64, assignedBy int64) (err error)
    RemoveUserRoles(ctx context.Context, userId int64, roleIds []int64) (err error)
    GetUserRoles(ctx context.Context, userId int64) (res []*sysout.RoleModel, err error)
    SetUserPrimaryRole(ctx context.Context, userId int64, roleId int64) (err error)

    // 权限验证
    CheckUserPermission(ctx context.Context, userId int64, permission string) (bool, error)
    CheckUserRole(ctx context.Context, userId int64, roleCode string) (bool, error)
    GetUserPermissions(ctx context.Context, userId int64) ([]string, error)
    GetUserMenus(ctx context.Context, userId int64) ([]int64, error)
    GetUserDataScope(ctx context.Context, userId int64) (int, error)

    // 批量权限验证
    CheckUsersPermission(ctx context.Context, userIds []int64, permission string) (map[int64]bool, error)
    FilterUsersByPermission(ctx context.Context, userIds []int64, permission string) ([]int64, error)
}
```

## 服务注册机制

### 接口实现注册

接口实现在`internal/logic`目录下，并通过`Register`函数注册：

```go
// internal/logic/api/role.go
type sRole struct{}

func NewRole() *sRole {
    return &sRole{}
}

// 在init函数中注册服务实现
func init() {
    service.RegisterRole(NewRole())
}

// 实现接口方法
func (s *sRole) GetRoleList(ctx context.Context, in *sysin.RoleListInp) (*sysout.RoleListModel, error) {
    // 参数过滤
    if err := in.Filter(ctx); err != nil {
        return nil, err
    }

    // 构建查询条件
    db := g.DB().Model("roles").Where("deleted_at IS NULL")
    
    // 具体业务逻辑实现
    // ...
    
    return result, nil
}
```

### 服务获取机制

服务层提供获取函数，供控制器调用：

```go
// internal/service/role.go
var (
    localRole IRole
)

func Role() IRole {
    if localRole == nil {
        panic("implement not found for interface IRole, forgot register?")
    }
    return localRole
}

func RegisterRole(i IRole) {
    localRole = i
}
```

## 服务使用方式

控制器通过服务接口获取功能：

```go
// internal/controller/api/role.go
func (c *cRole) GetRoleList(ctx context.Context, req *role.RoleListReq) (res *role.RoleListRes, err error) {
    // 通过服务接口调用业务逻辑
    out, err := service.Role().GetRoleList(ctx, &req.RoleListInp)
    if err != nil {
        return nil, err
    }
    return &role.RoleListRes{RoleListModel: out}, nil
}
```

## 接口功能分类

### 基础CRUD操作

```go
// 查询操作
GetRoleList(ctx context.Context, in *sysin.RoleListInp) (res *sysout.RoleListModel, err error)
GetRoleDetail(ctx context.Context, in *sysin.RoleDetailInp) (res *sysout.RoleDetailModel, err error)

// 创建和更新
CreateRole(ctx context.Context, in *sysin.CreateRoleInp) (res *sysout.RoleModel, err error)
UpdateRole(ctx context.Context, in *sysin.UpdateRoleInp) (res *sysout.RoleModel, err error)
CopyRole(ctx context.Context, in *sysin.CopyRoleInp) (res *sysout.RoleModel, err error)

// 删除操作
DeleteRole(ctx context.Context, in *sysin.DeleteRoleInp) (err error)
BatchDeleteRole(ctx context.Context, in *sysin.BatchDeleteRoleInp) (err error)

// 状态管理
UpdateRoleStatus(ctx context.Context, in *sysin.UpdateRoleStatusInp) (err error)
```

### 权限管理

```go
// 菜单权限
GetRoleMenus(ctx context.Context, in *sysin.RoleMenuInp) (res *sysout.RoleMenuModel, err error)
UpdateRoleMenus(ctx context.Context, in *sysin.UpdateRoleMenuInp) (err error)

// 权限详情
GetRolePermissions(ctx context.Context, in *sysin.RolePermissionInp) (res *sysout.RolePermissionModel, err error)

// 权限验证
CheckUserPermission(ctx context.Context, userId int64, permission string) (bool, error)
CheckUserRole(ctx context.Context, userId int64, roleCode string) (bool, error)
GetUserPermissions(ctx context.Context, userId int64) ([]string, error)
GetUserMenus(ctx context.Context, userId int64) ([]int64, error)
GetUserDataScope(ctx context.Context, userId int64) (int, error)

// 批量权限验证
CheckUsersPermission(ctx context.Context, userIds []int64, permission string) (map[int64]bool, error)
FilterUsersByPermission(ctx context.Context, userIds []int64, permission string) ([]int64, error)
```

### 用户角色关联

```go
// 角色分配
AssignUserRoles(ctx context.Context, userId int64, roleIds []int64, assignedBy int64) (err error)
RemoveUserRoles(ctx context.Context, userId int64, roleIds []int64) (err error)
GetUserRoles(ctx context.Context, userId int64) (res []*sysout.RoleModel, err error)
SetUserPrimaryRole(ctx context.Context, userId int64, roleId int64) (err error)
```

### 统计和选项

```go
// 选项数据
GetRoleOptions(ctx context.Context, in *sysin.RoleOptionInp) (res []*sysout.RoleOptionModel, err error)
GetDataScopeOptions(ctx context.Context) (res []*sysout.DataScopeModel, err error)

// 统计数据
GetRoleStats(ctx context.Context) (res *sysout.RoleStatsModel, err error)
```

## 接口设计原则

1. **接口隔离原则**: 接口应该小而精，按功能模块分组
2. **依赖倒置原则**: 控制器依赖抽象接口，而不是具体实现
3. **上下文传递**: 所有接口方法必须接收`context.Context`作为第一个参数
4. **错误处理**: 接口方法应返回标准错误类型，便于上层统一处理
5. **参数模型**: 使用专门的输入输出模型，而不是直接使用实体类型

## 参数模型设计

### 输入参数模型

```go
// 角色列表查询参数
type RoleListInp struct {
    Page      int    `json:"page" v:"min:1" dc:"当前页码"`
    PageSize  int    `json:"pageSize" v:"min:1|max:100" dc:"每页数量"`
    Name      string `json:"name" dc:"角色名称"`
    Code      string `json:"code" dc:"角色编码"`
    Status    int    `json:"status" dc:"状态"`
    DataScope int    `json:"dataScope" dc:"数据权限范围"`
    OrderBy   string `json:"orderBy" dc:"排序字段"`
    OrderType string `json:"orderType" dc:"排序类型"`
}

// 角色创建参数
type CreateRoleInp struct {
    Name        string   `json:"name" v:"required|length:1,50" dc:"角色名称"`
    Code        string   `json:"code" v:"required|length:1,50" dc:"角色编码"`
    Description string   `json:"description" dc:"角色描述"`
    Status      int      `json:"status" v:"in:0,1" dc:"状态"`
    Sort        int      `json:"sort" dc:"排序"`
    DataScope   int      `json:"dataScope" v:"in:1,2,3,4,5" dc:"数据权限范围"`
    Remark      string   `json:"remark" dc:"备注"`
    MenuIds     []int64  `json:"menuIds" dc:"菜单ID列表"`
}
```

### 输出结果模型

```go
// 角色列表响应模型
type RoleListModel struct {
    List     []*RoleModel `json:"list" dc:"角色列表"`
    Total    int64        `json:"total" dc:"总数"`
    Page     int          `json:"page" dc:"当前页"`
    PageSize int          `json:"pageSize" dc:"每页数量"`
}

// 角色详情模型
type RoleDetailModel struct {
    *RoleModel
    MenuIds     []int64  `json:"menuIds" dc:"菜单ID列表"`
    Permissions []string `json:"permissions" dc:"权限列表"`
}

// 角色统计模型
type RoleStatsModel struct {
    TotalCount    int64 `json:"totalCount" dc:"总角色数"`
    ActiveCount   int64 `json:"activeCount" dc:"活跃角色数"`
    InactiveCount int64 `json:"inactiveCount" dc:"禁用角色数"`
}
```

## 现有服务接口

| 接口名称    | 文件路径                        | 主要功能     |
| ----------- | ------------------------------- | ------------ |
| IRole       | internal/service/role.go        | 角色管理功能 |
| IUser       | internal/service/user.go        | 用户管理功能 |
| IMiddleware | internal/service/middleware.go  | 中间件功能   |
| IHook       | internal/service/hook.go        | 钩子函数     |
| IView       | internal/service/view.go        | 视图服务     |

## 添加新服务接口

1. 在`internal/service`目录下创建新的接口定义文件
2. 定义接口方法，遵循命名规范和设计原则
3. 提供服务获取函数和注册函数
4. 在`internal/logic`目录下实现接口
5. 在实现的`init()`函数中注册服务
6. 在控制器中通过服务获取函数调用

示例：

```go
// 1. 定义服务接口
// internal/service/menu.go
type IMenu interface {
    GetMenuList(ctx context.Context, in *sysin.MenuListInp) (res *sysout.MenuListModel, err error)
    CreateMenu(ctx context.Context, in *sysin.CreateMenuInp) (res *sysout.MenuModel, err error)
}

var localMenu IMenu

func Menu() IMenu {
    if localMenu == nil {
        panic("implement not found for interface IMenu, forgot register?")
    }
    return localMenu
}

func RegisterMenu(i IMenu) {
    localMenu = i
}

// 2. 实现服务接口
// internal/logic/api/menu.go
type sMenu struct{}

func NewMenu() *sMenu {
    return &sMenu{}
}

func init() {
    service.RegisterMenu(NewMenu())
}

func (s *sMenu) GetMenuList(ctx context.Context, in *sysin.MenuListInp) (*sysout.MenuListModel, error) {
    // 实现业务逻辑
    return result, nil
}

// 3. 控制器中使用
// internal/controller/api/menu.go
func (c *cMenu) GetMenuList(ctx context.Context, req *menu.MenuListReq) (res *menu.MenuListRes, err error) {
    out, err := service.Menu().GetMenuList(ctx, &req.MenuListInp)
    if err != nil {
        return nil, err
    }
    return &menu.MenuListRes{MenuListModel: out}, nil
}
```

## 总结

服务接口层在本项目中起到了关键的解耦作用：

1. **抽象化业务逻辑**：通过接口定义清晰的业务边界
2. **依赖注入**：通过注册机制实现依赖注入
3. **测试友好**：便于单元测试和模拟测试
4. **扩展性**：易于替换具体实现
5. **类型安全**：编译时检查接口实现

这种设计使得系统具有良好的可维护性和可测试性，同时保持了架构的清晰和简洁。
