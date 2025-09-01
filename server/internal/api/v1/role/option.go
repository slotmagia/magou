package v1

import (
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"

	"github.com/gogf/gf/v2/frame/g"
)

// RoleOptionReq 获取角色选项请求
type RoleOptionReq struct {
	g.Meta `path:"/role/options" method:"GET" summary:"获取角色选项" tags:"角色选项"`
	sysin.RoleOptionInp
}

// RoleOptionRes 获取角色选项响应
type RoleOptionRes struct {
	List []*sysout.RoleOptionModel `json:"list" description:"角色选项列表"`
}

// RoleStatsReq 获取角色统计请求
type RoleStatsReq struct {
	g.Meta `path:"/role/stats" method:"GET" summary:"获取角色统计" tags:"角色统计"`
}

// RoleStatsRes 获取角色统计响应
type RoleStatsRes struct {
	*sysout.RoleStatsModel
}

// DataScopeOptionReq 获取数据权限范围选项请求
type DataScopeOptionReq struct {
	g.Meta `path:"/role/data-scope-options" method:"GET" summary:"获取数据权限范围选项" tags:"角色选项"`
}

// DataScopeOptionRes 获取数据权限范围选项响应
type DataScopeOptionRes struct {
	List []*sysout.DataScopeModel `json:"list" description:"数据权限范围选项列表"`
}
