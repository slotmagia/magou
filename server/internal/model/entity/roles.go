package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Role 角色实体
type Role struct {
	Id          int64       `json:"id"          description:"主键ID"`
	Name        string      `json:"name"        description:"角色名称"`
	Code        string      `json:"code"        description:"角色编码"`
	Description string      `json:"description" description:"角色描述"`
	Status      int         `json:"status"      description:"状态：1=启用 0=禁用"`
	Sort        int         `json:"sort"        description:"排序号，数字越小越靠前"`
	DataScope   int         `json:"dataScope"   description:"数据权限范围"`
	Remark      string      `json:"remark"      description:"备注说明"`
	CreatedBy   int64       `json:"createdBy"   description:"创建人ID"`
	UpdatedBy   int64       `json:"updatedBy"   description:"修改人ID"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"更新时间"`
}

// RoleMenu 角色菜单关联实体
type RoleMenu struct {
	Id        int64       `json:"id"        description:"主键ID"`
	RoleId    int64       `json:"roleId"    description:"角色ID"`
	MenuId    int64       `json:"menuId"    description:"菜单ID"`
	CreatedAt *gtime.Time `json:"createdAt" description:"创建时间"`
}

// RoleStatus 角色状态常量
const (
	RoleStatusDisabled = 0 // 禁用
	RoleStatusEnabled  = 1 // 启用
	RoleStatusNormal   = 1 // 正常（别名）
)

// DataScope 数据权限范围常量
const (
	DataScopeAll        = 1 // 全部数据
	DataScopeDept       = 2 // 部门数据
	DataScopeDeptAndSub = 3 // 部门及以下数据
	DataScopeSelf       = 4 // 仅本人数据
	DataScopeCustom     = 5 // 自定义数据权限
)

// RoleCode 角色编码常量
const (
	RoleCodeFinanceAdmin    = "finance_admin"    // 财务管理员
	RoleCodeOperator        = "operator"         // 运营人员
	RoleCodeCustomerService = "customer_service" // 客服人员
	RoleCodeAuditor         = "auditor"          // 审计人员
)

// IsEnabled 判断角色是否启用
func (r *Role) IsEnabled() bool {
	return r.Status == RoleStatusEnabled
}

// IsDisabled 判断角色是否禁用
func (r *Role) IsDisabled() bool {
	return r.Status == RoleStatusDisabled
}

// IsSuperAdmin 判断是否为超级管理员
func (r *Role) IsSuperAdmin() bool {
	return r.Code == RoleCodeSuperAdmin
}

// IsSystemAdmin 判断是否为系统管理员
func (r *Role) IsSystemAdmin() bool {
	return r.Code == RoleCodeSystemAdmin
}

// IsBuiltIn 判断是否为内置角色（不可删除）
func (r *Role) IsBuiltIn() bool {
	builtInRoles := []string{
		RoleCodeSuperAdmin,
		RoleCodeSystemAdmin,

		RoleCodeFinanceAdmin,
		RoleCodeOperator,
		RoleCodeCustomerService,
		RoleCodeAuditor,
	}

	for _, code := range builtInRoles {
		if r.Code == code {
			return true
		}
	}
	return false
}

// GetStatusName 获取状态名称
func (r *Role) GetStatusName() string {
	switch r.Status {
	case RoleStatusEnabled:
		return "启用"
	case RoleStatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}

// GetDataScopeName 获取数据权限范围名称
func (r *Role) GetDataScopeName() string {
	switch r.DataScope {
	case DataScopeAll:
		return "全部数据"
	case DataScopeDept:
		return "部门数据"
	case DataScopeDeptAndSub:
		return "部门及以下数据"
	case DataScopeSelf:
		return "仅本人数据"
	case DataScopeCustom:
		return "自定义权限"
	default:
		return "未知"
	}
}

// HasDataScope 判断是否有指定的数据权限
func (r *Role) HasDataScope(scope int) bool {
	return r.DataScope == scope
}

// CanAccessAllData 判断是否可以访问全部数据
func (r *Role) CanAccessAllData() bool {
	return r.DataScope == DataScopeAll
}

// CanAccessDeptData 判断是否可以访问部门数据
func (r *Role) CanAccessDeptData() bool {
	return r.DataScope == DataScopeDept || r.DataScope == DataScopeDeptAndSub
}

// CanOnlyAccessSelfData 判断是否只能访问自己的数据
func (r *Role) CanOnlyAccessSelfData() bool {
	return r.DataScope == DataScopeSelf
}

// RoleWithMenus 带菜单权限的角色
type RoleWithMenus struct {
	Role
	MenuIds []int64 `json:"menuIds" description:"拥有的菜单ID列表"`
}

// RoleTree 角色树结构体（用于组织架构等场景）
type RoleTree struct {
	Role
	Children []*RoleTree `json:"children,omitempty" description:"子角色"`
}

// RolePermission 角色权限详情
type RolePermission struct {
	RoleId      int64    `json:"roleId"     description:"角色ID"`
	RoleName    string   `json:"roleName"   description:"角色名称"`
	RoleCode    string   `json:"roleCode"   description:"角色编码"`
	MenuIds     []int64  `json:"menuIds"    description:"菜单ID列表"`
	Permissions []string `json:"permissions" description:"权限标识列表"`
}

// ValidateDataScope 验证数据权限范围是否有效
func ValidateDataScope(scope int) bool {
	return scope >= DataScopeAll && scope <= DataScopeCustom
}

// GetAllDataScopes 获取所有数据权限范围选项
func GetAllDataScopes() map[int]string {
	return map[int]string{
		DataScopeAll:        "全部数据",
		DataScopeDept:       "部门数据",
		DataScopeDeptAndSub: "部门及以下数据",
		DataScopeSelf:       "仅本人数据",
		DataScopeCustom:     "自定义权限",
	}
}
