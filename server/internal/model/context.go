package model

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Context 请求上下文结构
type Context struct {
	Module    string    // 应用模块 admin｜sysin｜home｜websocket
	AddonName string    // 插件名称 如果不是插件模块请求，可能为空
	User      *Identity // 上下文用户信息
	Response  *Response // 请求响应
	Data      g.Map     // 自定kv变量 业务模块根据需要设置，不固定
}

// Identity 通用身份模型
type Identity struct {
	Id         int64       `json:"id"              description:"用户ID"`
	TenantId   int64       `json:"tenantId"        description:"租户ID"`
	TenantCode string      `json:"tenantCode"      description:"租户编码"`
	Pid        int64       `json:"pid"             description:"上级ID"`
	DeptId     int64       `json:"deptId"          description:"部门ID"`
	DeptType   string      `json:"deptType"        description:"部门类型"`
	RoleId     int64       `json:"roleId"          description:"角色ID"`
	RoleKey    string      `json:"roleKey"         description:"角色唯一标识符"`
	Username   string      `json:"username"        description:"用户名"`
	RealName   string      `json:"realName"        description:"姓名"`
	Avatar     string      `json:"avatar"          description:"头像"`
	Email      string      `json:"email"           description:"邮箱"`
	Mobile     string      `json:"mobile"          description:"手机号码"`
	App        string      `json:"app"             description:"登录应用"`
	LoginAt    *gtime.Time `json:"loginAt"         description:"登录时间"`
}

// IsTenantAdmin 判断是否为租户管理员
func (i *Identity) IsTenantAdmin() bool {
	return i.RoleKey == "tenant_admin"
}

// IsSystemAdmin 判断是否为系统管理员
func (i *Identity) IsSystemAdmin() bool {
	return i.TenantCode == "system" && (i.RoleKey == "super_admin" || i.RoleKey == "system_admin")
}
