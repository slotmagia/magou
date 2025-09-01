# 控制器结构

## 控制器概述

控制器是分层架构中的关键层，负责接收 HTTP 请求、参数验证、调用服务层接口、返回响应结果。在本项目中，控制器位于 `internal/controller` 目录下，按照业务模块进行分组。

## 控制器组织结构

控制器按照业务领域进行组织：

```
internal/controller/
├── api/               # API 控制器
│   ├── user.go        # 用户管理控制器
│   └── role.go        # 角色管理控制器
└── ... (其他业务领域)
```

## 控制器定义方式

控制器采用结构体定义，通过全局变量提供单例访问：

```go
// internal/controller/api/role.go
package api

import (
    role "client-app/internal/api/v1/role"
    "client-app/internal/service"
    "context"
    "fmt"
)

var (
    Role = cRole{} // 全局单例
)

type cRole struct{}

// GetRoleList 获取角色列表
func (c *cRole) GetRoleList(ctx context.Context, req *role.RoleListReq) (res *role.RoleListRes, err error) {
    out, err := service.Role().GetRoleList(ctx, &req.RoleListInp)
    if err != nil {
        return nil, err
    }
    return &role.RoleListRes{
        RoleListModel: out,
    }, nil
}
```

## 控制器接口定义

控制器对应的API接口定义在 `internal/api/v1/` 目录下，包含请求和响应结构：

```go
// internal/api/v1/role/role.go
package v1

import (
    "client-app/internal/model/input/sysin"
    "client-app/internal/model/output/sysout"
    "github.com/gogf/gf/v2/frame/g"
)

// 角色列表请求
type RoleListReq struct {
    g.Meta `path:"/role/list" method:"GET" summary:"获取角色列表" tags:"角色管理"`
    sysin.RoleListInp
}

// 角色列表响应
type RoleListRes struct {
    *sysout.RoleListModel
}

// 角色详情请求
type RoleDetailReq struct {
    g.Meta `path:"/role/{id}" method:"GET" summary:"获取角色详情" tags:"角色管理"`
    sysin.RoleDetailInp
}

// 角色详情响应
type RoleDetailRes struct {
    *sysout.RoleDetailModel
}
```

## 请求处理流程

控制器方法的典型处理流程如下：

1. **接收请求**：控制器方法接收上下文和请求参数
2. **参数验证**：框架基于结构体标签自动进行参数验证
3. **调用服务**：通过服务接口调用业务逻辑
4. **处理响应**：封装业务处理结果
5. **返回结果**：将结果返回给客户端

### 基础CRUD操作示例

```go
// 获取角色列表
func (c *cRole) GetRoleList(ctx context.Context, req *role.RoleListReq) (res *role.RoleListRes, err error) {
    out, err := service.Role().GetRoleList(ctx, &req.RoleListInp)
    if err != nil {
        return nil, err
    }
    return &role.RoleListRes{
        RoleListModel: out,
    }, nil
}

// 获取角色详情
func (c *cRole) GetRoleDetail(ctx context.Context, req *role.RoleDetailReq) (res *role.RoleDetailRes, err error) {
    out, err := service.Role().GetRoleDetail(ctx, &req.RoleDetailInp)
    if err != nil {
        return nil, err
    }
    return &role.RoleDetailRes{
        RoleDetailModel: out,
    }, nil
}

// 创建角色
func (c *cRole) CreateRole(ctx context.Context, req *role.CreateRoleReq) (res *role.CreateRoleRes, err error) {
    out, err := service.Role().CreateRole(ctx, &req.CreateRoleInp)
    if err != nil {
        return nil, err
    }
    return &role.CreateRoleRes{
        RoleModel: out,
    }, nil
}

// 更新角色
func (c *cRole) UpdateRole(ctx context.Context, req *role.UpdateRoleReq) (res *role.UpdateRoleRes, err error) {
    out, err := service.Role().UpdateRole(ctx, &req.UpdateRoleInp)
    if err != nil {
        return nil, err
    }
    return &role.UpdateRoleRes{
        RoleModel: out,
    }, nil
}
```

### 状态操作示例

```go
// 删除角色
func (c *cRole) DeleteRole(ctx context.Context, req *role.DeleteRoleReq) (res *role.DeleteRoleRes, err error) {
    err = service.Role().DeleteRole(ctx, &req.DeleteRoleInp)
    if err != nil {
        return &role.DeleteRoleRes{
            Success: false,
            Message: err.Error(),
        }, nil
    }
    return &role.DeleteRoleRes{
        Success: true,
        Message: "删除成功",
    }, nil
}

// 批量删除角色
func (c *cRole) BatchDeleteRole(ctx context.Context, req *role.BatchDeleteRoleReq) (res *role.BatchDeleteRoleRes, err error) {
    err = service.Role().BatchDeleteRole(ctx, &req.BatchDeleteRoleInp)
    if err != nil {
        return &role.BatchDeleteRoleRes{
            Success: false,
            Message: err.Error(),
        }, nil
    }
    return &role.BatchDeleteRoleRes{
        Success: true,
        Message: "批量删除成功",
    }, nil
}

// 更新角色状态
func (c *cRole) UpdateRoleStatus(ctx context.Context, req *role.UpdateRoleStatusReq) (res *role.UpdateRoleStatusRes, err error) {
    err = service.Role().UpdateRoleStatus(ctx, &req.UpdateRoleStatusInp)
    if err != nil {
        return &role.UpdateRoleStatusRes{
            Success: false,
            Message: err.Error(),
        }, nil
    }
    return &role.UpdateRoleStatusRes{
        Success: true,
        Message: "状态更新成功",
    }, nil
}
```

### 权限管理示例

```go
// 获取角色菜单权限
func (c *cRole) GetRoleMenus(ctx context.Context, req *role.RoleMenuReq) (res *role.RoleMenuRes, err error) {
    out, err := service.Role().GetRoleMenus(ctx, &req.RoleMenuInp)
    if err != nil {
        return nil, err
    }
    return &role.RoleMenuRes{
        RoleMenuModel: out,
    }, nil
}

// 更新角色菜单权限
func (c *cRole) UpdateRoleMenus(ctx context.Context, req *role.UpdateRoleMenuReq) (res *role.UpdateRoleMenuRes, err error) {
    err = service.Role().UpdateRoleMenus(ctx, &req.UpdateRoleMenuInp)
    if err != nil {
        return &role.UpdateRoleMenuRes{
            Success: false,
            Message: err.Error(),
        }, nil
    }
    return &role.UpdateRoleMenuRes{
        Success: true,
        Message: "权限更新成功",
    }, nil
}

// 获取角色权限详情
func (c *cRole) GetRolePermissions(ctx context.Context, req *role.RolePermissionReq) (res *role.RolePermissionRes, err error) {
    out, err := service.Role().GetRolePermissions(ctx, &req.RolePermissionInp)
    if err != nil {
        return nil, err
    }
    return &role.RolePermissionRes{
        RolePermissionModel: out,
    }, nil
}
```

## 参数绑定和验证

控制器方法的参数绑定由框架自动完成，基于结构体标签进行验证：

```go
// 请求参数结构体
type CreateRoleReq struct {
    g.Meta `path:"/role" method:"POST" summary:"创建角色" tags:"角色管理"`
    sysin.CreateRoleInp
}

// 输入参数结构体
type CreateRoleInp struct {
    Name        string `json:"name" v:"required|length:1,50" dc:"角色名称"`
    Code        string `json:"code" v:"required|length:1,50" dc:"角色编码"`
    Description string `json:"description" dc:"角色描述"`
    Status      int    `json:"status" v:"in:0,1" dc:"状态：0-禁用，1-启用"`
    Sort        int    `json:"sort" dc:"排序"`
    DataScope   int    `json:"dataScope" v:"in:1,2,3,4,5" dc:"数据权限范围"`
    MenuIds     []int64 `json:"menuIds" dc:"菜单ID列表"`
}
```

## 错误处理

控制器的错误处理策略：

1. **服务层错误直接返回**：让框架统一处理
2. **业务逻辑错误**：返回带错误信息的响应结构
3. **统一响应格式**：通过中间件统一封装

```go
// 错误处理示例
func (c *cRole) DeleteRole(ctx context.Context, req *role.DeleteRoleReq) (res *role.DeleteRoleRes, err error) {
    err = service.Role().DeleteRole(ctx, &req.DeleteRoleInp)
    if err != nil {
        // 返回业务错误响应，而不是直接返回错误
        return &role.DeleteRoleRes{
            Success: false,
            Message: err.Error(),
        }, nil
    }
    
    return &role.DeleteRoleRes{
        Success: true,
        Message: "删除成功",
    }, nil
}
```

## 依赖注入

控制器通过服务注册机制调用服务层：

```go
// 调用服务层接口
out, err := service.Role().GetRoleList(ctx, &req.RoleListInp)

// 服务层接口定义在 internal/service/role.go
// 具体实现在 internal/logic/api/role.go
```

## 控制器命名规范

1. **文件名**：使用业务模块名称，如 `role.go`、`user.go`
2. **结构体名**：使用 `c` 前缀 + 业务模块名称，如 `cRole`、`cUser`
3. **方法名**：使用业务动作名称，如 `GetRoleList`、`CreateRole`、`UpdateRole`
4. **全局变量**：使用业务模块名称，如 `Role`、`User`

## 控制器完整示例

```go
package api

import (
    role "client-app/internal/api/v1/role"
    "client-app/internal/service"
    "context"
)

var (
    Role = cRole{}
)

type cRole struct{}

// 角色基础操作
func (c *cRole) GetRoleList(ctx context.Context, req *role.RoleListReq) (res *role.RoleListRes, err error) {
    out, err := service.Role().GetRoleList(ctx, &req.RoleListInp)
    if err != nil {
        return nil, err
    }
    return &role.RoleListRes{RoleListModel: out}, nil
}

func (c *cRole) GetRoleDetail(ctx context.Context, req *role.RoleDetailReq) (res *role.RoleDetailRes, err error) {
    out, err := service.Role().GetRoleDetail(ctx, &req.RoleDetailInp)
    if err != nil {
        return nil, err
    }
    return &role.RoleDetailRes{RoleDetailModel: out}, nil
}

func (c *cRole) CreateRole(ctx context.Context, req *role.CreateRoleReq) (res *role.CreateRoleRes, err error) {
    out, err := service.Role().CreateRole(ctx, &req.CreateRoleInp)
    if err != nil {
        return nil, err
    }
    return &role.CreateRoleRes{RoleModel: out}, nil
}

// 角色权限管理
func (c *cRole) GetRoleMenus(ctx context.Context, req *role.RoleMenuReq) (res *role.RoleMenuRes, err error) {
    out, err := service.Role().GetRoleMenus(ctx, &req.RoleMenuInp)
    if err != nil {
        return nil, err
    }
    return &role.RoleMenuRes{RoleMenuModel: out}, nil
}

// 统计和选项
func (c *cRole) GetRoleOptions(ctx context.Context, req *role.RoleOptionReq) (res *role.RoleOptionRes, err error) {
    out, err := service.Role().GetRoleOptions(ctx, &req.RoleOptionInp)
    if err != nil {
        return nil, err
    }
    return &role.RoleOptionRes{List: out}, nil
}

func (c *cRole) GetRoleStats(ctx context.Context, req *role.RoleStatsReq) (res *role.RoleStatsRes, err error) {
    out, err := service.Role().GetRoleStats(ctx)
    if err != nil {
        return nil, err
    }
    return &role.RoleStatsRes{RoleStatsModel: out}, nil
}
```

## 总结

控制器层在本项目中承担着以下职责：

1. **HTTP接口实现**：实现具体的API端点
2. **参数验证**：利用框架的自动验证机制
3. **服务调用**：通过服务接口调用业务逻辑
4. **响应封装**：将业务结果封装为HTTP响应
5. **错误处理**：处理业务错误并返回适当的响应

通过这种设计，控制器层保持了简洁性，专注于HTTP协议相关的处理，而将业务逻辑委托给服务层处理。
