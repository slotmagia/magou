package tenant

import (
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"
	"github.com/gogf/gf/v2/frame/g"
)

// 租户列表请求
type TenantListReq struct {
	g.Meta `path:"/tenant/list" method:"get" summary:"获取租户列表" tags:"租户管理"`
	sysin.TenantListInp
}

type TenantListRes struct {
	*sysout.TenantListModel
}

// 创建租户请求
type CreateTenantReq struct {
	g.Meta `path:"/tenant/create" method:"post" summary:"创建租户" tags:"租户管理"`
	sysin.CreateTenantInp
}

type CreateTenantRes struct {
	*sysout.TenantModel
}

// 更新租户请求
type UpdateTenantReq struct {
	g.Meta `path:"/tenant/update" method:"put" summary:"更新租户" tags:"租户管理"`
	sysin.UpdateTenantInp
}

type UpdateTenantRes struct {
	*sysout.TenantModel
}

// 删除租户请求
type DeleteTenantReq struct {
	g.Meta `path:"/tenant/delete" method:"delete" summary:"删除租户" tags:"租户管理"`
	sysin.DeleteTenantInp
}

type DeleteTenantRes struct{}

// 获取租户详情请求
type TenantDetailReq struct {
	g.Meta `path:"/tenant/detail" method:"get" summary:"获取租户详情" tags:"租户管理"`
	sysin.TenantDetailInp
}

type TenantDetailRes struct {
	*sysout.TenantDetailModel
}

// 更新租户状态请求
type TenantStatusReq struct {
	g.Meta `path:"/tenant/status" method:"put" summary:"更新租户状态" tags:"租户管理"`
	sysin.TenantStatusInp
}

type TenantStatusRes struct{}

// 获取租户统计请求
type TenantStatsReq struct {
	g.Meta `path:"/tenant/stats" method:"get" summary:"获取租户统计" tags:"租户管理"`
	sysin.TenantStatsInp
}

type TenantStatsRes struct {
	*sysout.TenantStatsModel
}

// 更新租户配置请求
type TenantConfigReq struct {
	g.Meta `path:"/tenant/config" method:"put" summary:"更新租户配置" tags:"租户管理"`
	sysin.TenantConfigInp
}

type TenantConfigRes struct{}

// 获取租户选项请求
type TenantOptionsReq struct {
	g.Meta `path:"/tenant/options" method:"get" summary:"获取租户选项" tags:"租户管理"`
}

type TenantOptionsRes struct {
	*sysout.TenantOptionsModel
}
