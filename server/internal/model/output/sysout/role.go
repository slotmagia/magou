package sysout

import (
	"client-app/internal/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
)

// RoleListModel 角色列表响应模型
type RoleListModel struct {
	List     []*RoleModel `json:"list" description:"角色列表"`
	Total    int64        `json:"total" description:"总记录数"`
	Page     int          `json:"page" description:"当前页码"`
	PageSize int          `json:"pageSize" description:"每页数量"`
}

// RoleModel 角色基础响应模型
type RoleModel struct {
	Id            int64       `json:"id" description:"主键ID"`
	Name          string      `json:"name" description:"角色名称"`
	Code          string      `json:"code" description:"角色编码"`
	Description   string      `json:"description" description:"角色描述"`
	Status        int         `json:"status" description:"状态"`
	StatusName    string      `json:"statusName" description:"状态名称"`
	Sort          int         `json:"sort" description:"排序号"`
	DataScope     int         `json:"dataScope" description:"数据权限范围"`
	DataScopeName string      `json:"dataScopeName" description:"数据权限范围名称"`
	Remark        string      `json:"remark" description:"备注说明"`
	IsBuiltIn     bool        `json:"isBuiltIn" description:"是否内置角色"`
	CreatedBy     int64       `json:"createdBy" description:"创建人ID"`
	UpdatedBy     int64       `json:"updatedBy" description:"修改人ID"`
	CreatedAt     *gtime.Time `json:"createdAt" description:"创建时间"`
	UpdatedAt     *gtime.Time `json:"updatedAt" description:"更新时间"`
}

// RoleDetailModel 角色详情响应模型
type RoleDetailModel struct {
	RoleModel
	MenuIds     []int64      `json:"menuIds" description:"拥有的菜单ID列表"`
	Permissions []string     `json:"permissions" description:"权限标识列表"`
	MenuList    []*MenuModel `json:"menuList,omitempty" description:"菜单详情列表"`
}

// RoleMenuModel 角色菜单权限响应模型
type RoleMenuModel struct {
	RoleId      int64            `json:"roleId" description:"角色ID"`
	RoleName    string           `json:"roleName" description:"角色名称"`
	RoleCode    string           `json:"roleCode" description:"角色编码"`
	MenuIds     []int64          `json:"menuIds" description:"拥有的菜单ID列表"`
	MenuTree    []*MenuTreeModel `json:"menuTree" description:"菜单树结构"`
	Permissions []string         `json:"permissions" description:"权限标识列表"`
}

// RoleOptionModel 角色选项模型（用于下拉框等）
type RoleOptionModel struct {
	Value     int64  `json:"value" description:"角色ID"`
	Label     string `json:"label" description:"角色名称"`
	Code      string `json:"code" description:"角色编码"`
	DataScope int    `json:"dataScope" description:"数据权限范围"`
	Disabled  bool   `json:"disabled" description:"是否禁用"`
	IsBuiltIn bool   `json:"isBuiltIn" description:"是否内置角色"`
}

// RolePermissionModel 角色权限详情模型
type RolePermissionModel struct {
	RoleId      int64    `json:"roleId" description:"角色ID"`
	RoleName    string   `json:"roleName" description:"角色名称"`
	RoleCode    string   `json:"roleCode" description:"角色编码"`
	DataScope   int      `json:"dataScope" description:"数据权限范围"`
	MenuIds     []int64  `json:"menuIds" description:"菜单ID列表"`
	Permissions []string `json:"permissions" description:"权限标识列表"`
	Features    []string `json:"features" description:"功能特性列表"`
}

// DataScopeModel 数据权限范围选项模型
type DataScopeModel struct {
	Value int    `json:"value" description:"权限范围值"`
	Label string `json:"label" description:"权限范围名称"`
}

// RoleStatsModel 角色统计模型
type RoleStatsModel struct {
	TotalCount    int64 `json:"totalCount" description:"总角色数"`
	EnabledCount  int64 `json:"enabledCount" description:"启用角色数"`
	DisabledCount int64 `json:"disabledCount" description:"禁用角色数"`
	BuiltInCount  int64 `json:"builtInCount" description:"内置角色数"`
	CustomCount   int64 `json:"customCount" description:"自定义角色数"`
}

// ConvertToRoleModel 将entity.Role转换为RoleModel
func ConvertToRoleModel(role *entity.Role) *RoleModel {
	if role == nil {
		return nil
	}

	return &RoleModel{
		Id:            role.Id,
		Name:          role.Name,
		Code:          role.Code,
		Description:   role.Description,
		Status:        role.Status,
		StatusName:    role.GetStatusName(),
		Sort:          role.Sort,
		DataScope:     role.DataScope,
		DataScopeName: role.GetDataScopeName(),
		Remark:        role.Remark,
		IsBuiltIn:     role.IsBuiltIn(),
		CreatedBy:     role.CreatedBy,
		UpdatedBy:     role.UpdatedBy,
		CreatedAt:     role.CreatedAt,
		UpdatedAt:     role.UpdatedAt,
	}
}

// ConvertToRoleDetailModel 将entity.Role转换为RoleDetailModel
func ConvertToRoleDetailModel(role *entity.Role, menuIds []int64, permissions []string) *RoleDetailModel {
	if role == nil {
		return nil
	}

	return &RoleDetailModel{
		RoleModel:   *ConvertToRoleModel(role),
		MenuIds:     menuIds,
		Permissions: permissions,
	}
}

// ConvertToRoleOptionModel 将entity.Role转换为RoleOptionModel
func ConvertToRoleOptionModel(role *entity.Role) *RoleOptionModel {
	if role == nil {
		return nil
	}

	return &RoleOptionModel{
		Value:     role.Id,
		Label:     role.Name,
		Code:      role.Code,
		DataScope: role.DataScope,
		Disabled:  role.IsDisabled(),
		IsBuiltIn: role.IsBuiltIn(),
	}
}

// ConvertToRolePermissionModel 将entity.Role转换为RolePermissionModel
func ConvertToRolePermissionModel(role *entity.Role, menuIds []int64, permissions []string) *RolePermissionModel {
	if role == nil {
		return nil
	}

	// 根据角色类型生成功能特性列表
	features := generateRoleFeatures(role)

	return &RolePermissionModel{
		RoleId:      role.Id,
		RoleName:    role.Name,
		RoleCode:    role.Code,
		DataScope:   role.DataScope,
		MenuIds:     menuIds,
		Permissions: permissions,
		Features:    features,
	}
}

// GetAllDataScopes 获取所有数据权限范围选项
func GetAllDataScopes() []*DataScopeModel {
	scopes := entity.GetAllDataScopes()
	result := make([]*DataScopeModel, 0, len(scopes))

	// 按顺序排列
	order := []int{
		entity.DataScopeAll,
		entity.DataScopeDept,
		entity.DataScopeDeptAndSub,
		entity.DataScopeSelf,
		entity.DataScopeCustom,
	}

	for _, value := range order {
		if label, exists := scopes[value]; exists {
			result = append(result, &DataScopeModel{
				Value: value,
				Label: label,
			})
		}
	}

	return result
}

// generateRoleFeatures 根据角色类型生成功能特性列表
func generateRoleFeatures(role *entity.Role) []string {
	features := make([]string, 0)

	switch role.Code {
	case entity.RoleCodeSuperAdmin:
		features = append(features, "系统最高权限", "所有功能访问", "用户管理", "系统配置")
	case entity.RoleCodeSystemAdmin:
		features = append(features, "系统管理", "用户管理", "菜单管理", "日志管理")

	case entity.RoleCodeFinanceAdmin:
		features = append(features, "财务统计", "报表导出", "数据分析", "订单查看")
	case entity.RoleCodeOperator:
		features = append(features, "日常操作", "订单处理", "基础查询")
	case entity.RoleCodeCustomerService:
		features = append(features, "客户服务", "订单查询", "退款处理")
	case entity.RoleCodeAuditor:
		features = append(features, "审计监督", "只读权限", "日志查看")
	default:
		features = append(features, "自定义角色")
	}

	// 根据数据权限范围添加特性
	switch role.DataScope {
	case entity.DataScopeAll:
		features = append(features, "全部数据访问")
	case entity.DataScopeDept:
		features = append(features, "部门数据访问")
	case entity.DataScopeDeptAndSub:
		features = append(features, "部门及下级数据访问")
	case entity.DataScopeSelf:
		features = append(features, "仅本人数据访问")
	case entity.DataScopeCustom:
		features = append(features, "自定义数据权限")
	}

	return features
}

// BuildRoleStatsModel 构建角色统计模型
func BuildRoleStatsModel(roles []*entity.Role) *RoleStatsModel {
	stats := &RoleStatsModel{}

	for _, role := range roles {
		stats.TotalCount++

		if role.IsEnabled() {
			stats.EnabledCount++
		} else {
			stats.DisabledCount++
		}

		if role.IsBuiltIn() {
			stats.BuiltInCount++
		} else {
			stats.CustomCount++
		}
	}

	return stats
}

// FilterRolesByPermission 根据权限过滤角色（用于权限控制）
func FilterRolesByPermission(roles []*RoleModel, userDataScope int, userDeptId int64) []*RoleModel {
	if userDataScope == entity.DataScopeAll {
		return roles // 全部数据权限，返回所有角色
	}

	filtered := make([]*RoleModel, 0)
	for _, role := range roles {
		// 根据用户的数据权限范围过滤角色
		switch userDataScope {
		case entity.DataScopeDept, entity.DataScopeDeptAndSub:
			// 部门权限，可以看到同级或下级角色
			if role.DataScope >= userDataScope {
				filtered = append(filtered, role)
			}
		case entity.DataScopeSelf:
			// 仅本人权限，只能看到本人创建的角色
			if role.CreatedBy == userDeptId {
				filtered = append(filtered, role)
			}
		}
	}

	return filtered
}
