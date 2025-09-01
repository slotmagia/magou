package menu

import (
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"

	"github.com/gogf/gf/v2/frame/g"
)

// MenuListReq 菜单列表查询请求
type MenuListReq struct {
	g.Meta `path:"/menu/list" method:"GET" summary:"获取菜单列表" tags:"菜单管理"`
	sysin.MenuListInp
}

// MenuListRes 菜单列表查询响应
type MenuListRes struct {
	*sysout.MenuListModel
}

// MenuTreeReq 菜单树查询请求
type MenuTreeReq struct {
	g.Meta `path:"/menu/tree" method:"GET" summary:"获取菜单树" tags:"菜单管理"`
	sysin.MenuTreeInp
}

// MenuTreeRes 菜单树查询响应
type MenuTreeRes struct {
	List []*sysout.MenuTreeModel `json:"list" description:"菜单树列表"`
}

// MenuDetailReq 菜单详情查询请求
type MenuDetailReq struct {
	g.Meta `path:"/menu/{id}" method:"GET" summary:"获取菜单详情" tags:"菜单管理"`
	sysin.MenuDetailInp
}

// MenuDetailRes 菜单详情查询响应
type MenuDetailRes struct {
	*sysout.MenuDetailModel
}

// CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	g.Meta `path:"/menu" method:"POST" summary:"创建菜单" tags:"菜单管理"`
	sysin.CreateMenuInp
}

// CreateMenuRes 创建菜单响应
type CreateMenuRes struct {
	*sysout.MenuModel
}

// UpdateMenuReq 更新菜单请求
type UpdateMenuReq struct {
	g.Meta `path:"/menu/{id}" method:"PUT" summary:"更新菜单" tags:"菜单管理"`
	sysin.UpdateMenuInp
}

// UpdateMenuRes 更新菜单响应
type UpdateMenuRes struct {
	*sysout.MenuModel
}

// DeleteMenuReq 删除菜单请求
type DeleteMenuReq struct {
	g.Meta `path:"/menu/{id}" method:"DELETE" summary:"删除菜单" tags:"菜单管理"`
	sysin.DeleteMenuInp
}

// DeleteMenuRes 删除菜单响应
type DeleteMenuRes struct {
	Success bool   `json:"success" description:"是否成功"`
	Message string `json:"message" description:"提示信息"`
}

// UpdateMenuStatusReq 更新菜单状态请求
type UpdateMenuStatusReq struct {
	g.Meta `path:"/menu/{id}/status" method:"PUT" summary:"更新菜单状态" tags:"菜单管理"`
	sysin.UpdateMenuStatusInp
}

// UpdateMenuStatusRes 更新菜单状态响应
type UpdateMenuStatusRes struct {
	Success bool   `json:"success" description:"是否成功"`
	Message string `json:"message" description:"提示信息"`
}

// BatchDeleteMenuReq 批量删除菜单请求
type BatchDeleteMenuReq struct {
	g.Meta `path:"/menu/batch/delete" method:"DELETE" summary:"批量删除菜单" tags:"菜单管理"`
	sysin.BatchDeleteMenuInp
}

// BatchDeleteMenuRes 批量删除菜单响应
type BatchDeleteMenuRes struct {
	Success bool   `json:"success" description:"是否成功"`
	Message string `json:"message" description:"提示信息"`
}

// MenuOptionsReq 菜单选项请求
type MenuOptionsReq struct {
	g.Meta `path:"/menu/options" method:"GET" summary:"获取菜单选项" tags:"菜单管理"`
	sysin.MenuOptionInp
}

// MenuOptionsRes 菜单选项响应
type MenuOptionsRes struct {
	List []*sysout.MenuOptionModel `json:"list" description:"菜单选项列表"`
}

// RoutersReq 前端路由请求
type RoutersReq struct {
	g.Meta `path:"/menu/routers" method:"GET" summary:"获取前端路由" tags:"菜单管理"`
}

// RoutersRes 前端路由响应
type RoutersRes struct {
	List []*sysout.RouterModel `json:"list" description:"路由列表"`
}
