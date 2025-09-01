package v1

import (
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"

	"github.com/gogf/gf/v2/frame/g"
)

// RoleMenuReq 获取角色菜单权限请求
type RoleMenuReq struct {
	g.Meta `path:"/role/{id}/menus" method:"GET" summary:"获取角色菜单权限" tags:"角色权限"`
	sysin.RoleMenuInp
}

// RoleMenuRes 获取角色菜单权限响应
type RoleMenuRes struct {
	*sysout.RoleMenuModel
}

// UpdateRoleMenuReq 更新角色菜单权限请求
type UpdateRoleMenuReq struct {
	g.Meta `path:"/role/{id}/menus" method:"PUT" summary:"更新角色菜单权限" tags:"角色权限"`
	sysin.UpdateRoleMenuInp
}

// UpdateRoleMenuRes 更新角色菜单权限响应
type UpdateRoleMenuRes struct {
	Success bool   `json:"success" description:"是否成功"`
	Message string `json:"message" description:"提示信息"`
}

// RolePermissionReq 获取角色权限详情请求
type RolePermissionReq struct {
	g.Meta `path:"/role/{id}/permissions" method:"GET" summary:"获取角色权限详情" tags:"角色权限"`
	sysin.RolePermissionInp
}

// RolePermissionRes 获取角色权限详情响应
type RolePermissionRes struct {
	*sysout.RolePermissionModel
}
