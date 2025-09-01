package v1

import (
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"

	"github.com/gogf/gf/v2/frame/g"
)

// RoleListReq 角色列表请求
type RoleListReq struct {
	g.Meta `path:"/role/list" method:"GET" summary:"获取角色列表" tags:"角色管理"`
	sysin.RoleListInp
}

// RoleListRes 角色列表响应
type RoleListRes struct {
	*sysout.RoleListModel
}

// RoleDetailReq 角色详情请求
type RoleDetailReq struct {
	g.Meta `path:"/role/{id}" method:"GET" summary:"获取角色详情" tags:"角色管理"`
	sysin.RoleDetailInp
}

// RoleDetailRes 角色详情响应
type RoleDetailRes struct {
	*sysout.RoleDetailModel
}

// CreateRoleReq 创建角色请求
type CreateRoleReq struct {
	g.Meta `path:"/role" method:"POST" summary:"创建角色" tags:"角色管理"`
	sysin.CreateRoleInp
}

// CreateRoleRes 创建角色响应
type CreateRoleRes struct {
	*sysout.RoleModel
}

// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	g.Meta `path:"/role/{id}" method:"PUT" summary:"更新角色" tags:"角色管理"`
	sysin.UpdateRoleInp
}

// UpdateRoleRes 更新角色响应
type UpdateRoleRes struct {
	*sysout.RoleModel
}

// DeleteRoleReq 删除角色请求
type DeleteRoleReq struct {
	g.Meta `path:"/role/{id}" method:"DELETE" summary:"删除角色" tags:"角色管理"`
	sysin.DeleteRoleInp
}

// DeleteRoleRes 删除角色响应
type DeleteRoleRes struct {
	Success bool   `json:"success" description:"是否成功"`
	Message string `json:"message" description:"提示信息"`
}

// BatchDeleteRoleReq 批量删除角色请求
type BatchDeleteRoleReq struct {
	g.Meta `path:"/role/batch" method:"DELETE" summary:"批量删除角色" tags:"角色管理"`
	sysin.BatchDeleteRoleInp
}

// BatchDeleteRoleRes 批量删除角色响应
type BatchDeleteRoleRes struct {
	Success bool   `json:"success" description:"是否成功"`
	Message string `json:"message" description:"提示信息"`
}

// UpdateRoleStatusReq 更新角色状态请求
type UpdateRoleStatusReq struct {
	g.Meta `path:"/role/{id}/status" method:"PUT" summary:"更新角色状态" tags:"角色管理"`
	sysin.UpdateRoleStatusInp
}

// UpdateRoleStatusRes 更新角色状态响应
type UpdateRoleStatusRes struct {
	Success bool   `json:"success" description:"是否成功"`
	Message string `json:"message" description:"提示信息"`
}

// CopyRoleReq 复制角色请求
type CopyRoleReq struct {
	g.Meta `path:"/role/{id}/copy" method:"POST" summary:"复制角色" tags:"角色管理"`
	sysin.CopyRoleInp
}

// CopyRoleRes 复制角色响应
type CopyRoleRes struct {
	*sysout.RoleModel
}
