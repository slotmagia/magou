package sysin

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TenantListInp 租户列表查询参数
type TenantListInp struct {
	Page     int    `json:"page"     v:"min:1"              d:"1"          description:"页码"`
	PageSize int    `json:"pageSize" v:"min:1|max:100"      d:"20"         description:"每页数量"`
	Name     string `json:"name"     v:"length:0,100"                      description:"租户名称"`
	Code     string `json:"code"     v:"length:0,50"                       description:"租户编码"`
	Domain   string `json:"domain"   v:"length:0,100"                      description:"租户域名"`
	Status   int    `json:"status"   v:"in:0,1,2,3"         d:"0"          description:"状态：0=全部 1=正常 2=锁定 3=禁用"`
}

// CreateTenantInp 创建租户输入参数
type CreateTenantInp struct {
	Name         string      `json:"name"         v:"required|length:1,100#租户名称不能为空|租户名称长度不能超过100字符"`
	Code         string      `json:"code"         v:"required|length:1,50#租户编码不能为空|租户编码长度不能超过50字符"`
	Domain       string      `json:"domain"       v:"length:0,100#租户域名长度不能超过100字符"`
	MaxUsers     int         `json:"maxUsers"     v:"min:1|max:10000#最大用户数不能小于1|最大用户数不能超过10000"`
	StorageLimit int64       `json:"storageLimit" v:"min:0#存储限制不能小于0"`
	ExpireAt     *gtime.Time `json:"expireAt"     description:"过期时间"`
	AdminName    string      `json:"adminName"    v:"required|length:1,50#管理员用户名不能为空|用户名长度不能超过50字符"`
	AdminEmail   string      `json:"adminEmail"   v:"required|email#管理员邮箱不能为空|邮箱格式不正确"`
	AdminPassword string     `json:"adminPassword" v:"required|length:6,32#管理员密码不能为空|密码长度必须在6-32位之间"`
	Remark       string      `json:"remark"       v:"length:0,500#备注长度不能超过500字符"`
}

// 参数过滤和验证方法
func (in *CreateTenantInp) Filter(ctx context.Context) error {
	return g.Validator().Data(in).Run(ctx)
}

// UpdateTenantInp 更新租户输入参数
type UpdateTenantInp struct {
	Id           uint64      `json:"id"           v:"required|min:1#租户ID不能为空"`
	Name         string      `json:"name"         v:"required|length:1,100#租户名称不能为空|租户名称长度不能超过100字符"`
	Domain       string      `json:"domain"       v:"length:0,100#租户域名长度不能超过100字符"`
	MaxUsers     int         `json:"maxUsers"     v:"min:1|max:10000#最大用户数不能小于1|最大用户数不能超过10000"`
	StorageLimit int64       `json:"storageLimit" v:"min:0#存储限制不能小于0"`
	ExpireAt     *gtime.Time `json:"expireAt"     description:"过期时间"`
	Remark       string      `json:"remark"       v:"length:0,500#备注长度不能超过500字符"`
}

// 参数过滤和验证方法
func (in *UpdateTenantInp) Filter(ctx context.Context) error {
	return g.Validator().Data(in).Run(ctx)
}

// DeleteTenantInp 删除租户输入参数
type DeleteTenantInp struct {
	Id uint64 `json:"id" v:"required|min:1#租户ID不能为空"`
}

// 参数过滤和验证方法
func (in *DeleteTenantInp) Filter(ctx context.Context) error {
	return g.Validator().Data(in).Run(ctx)
}

// TenantDetailInp 租户详情查询参数
type TenantDetailInp struct {
	Id uint64 `json:"id" v:"required|min:1#租户ID不能为空"`
}

// 参数过滤和验证方法
func (in *TenantDetailInp) Filter(ctx context.Context) error {
	return g.Validator().Data(in).Run(ctx)
}

// TenantStatusInp 租户状态修改参数
type TenantStatusInp struct {
	Id     uint64 `json:"id"     v:"required|min:1#租户ID不能为空"`
	Status int    `json:"status" v:"required|in:1,2,3#状态值不正确"`
}

// 参数过滤和验证方法
func (in *TenantStatusInp) Filter(ctx context.Context) error {
	return g.Validator().Data(in).Run(ctx)
}

// TenantStatsInp 租户统计查询参数
type TenantStatsInp struct {
	Id uint64 `json:"id" v:"required|min:1#租户ID不能为空"`
}

// 参数过滤和验证方法
func (in *TenantStatsInp) Filter(ctx context.Context) error {
	return g.Validator().Data(in).Run(ctx)
}

// TenantConfigInp 租户配置输入参数
type TenantConfigInp struct {
	Id     uint64                 `json:"id"     v:"required|min:1#租户ID不能为空"`
	Config map[string]interface{} `json:"config" v:"required#配置不能为空"`
}

// 参数过滤和验证方法
func (in *TenantConfigInp) Filter(ctx context.Context) error {
	return g.Validator().Data(in).Run(ctx)
}
