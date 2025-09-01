package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Menu 菜单实体
type Menu struct {
	Id         int64       `json:"id"          description:"主键ID"`
	TenantId   int64       `json:"tenant_id"   description:"租户id"`
	ParentId   int64       `json:"parentId"    description:"父菜单ID，0表示顶级菜单"`
	MenuCode   string      `json:"menu_code"   description:"菜单编码"`
	Icon       string      `json:"icon"        description:"菜单图标"`
	Path       string      `json:"path"        description:"菜单路径"`
	Component  string      `json:"component"   description:"组件路径"`
	Permission string      `json:"permission"  description:"权限标识"`
	SortOrder  int         `json:"sort_order"  description:"排序号，数字越小越靠前"`
	Visible    int         `json:"visible"     description:"是否显示：1=显示 0=隐藏"`
	Status     int         `json:"status"      description:"状态：1=启用 0=禁用"`
	CreatedBy  int64       `json:"createdBy"   description:"创建人ID"`
	UpdatedBy  int64       `json:"updatedBy"   description:"修改人ID"`
	CreatedAt  *gtime.Time `json:"createdAt"   description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updatedAt"   description:"更新时间"`
	Redirect   string      `json:"redirect"    description:"重定向地址"`
	ActiveMenu string      `json:"activeMenu"  description:"高亮菜单路径"`
	AlwaysShow int         `json:"alwaysShow"  description:"是否总是显示：1=是 0=否"`
	Breadcrumb int         `json:"breadcrumb"  description:"是否显示面包屑：1=显示 0=隐藏"`
	Remark     string      `json:"remark"      description:"备注说明"`

	Title    string `json:"title"       description:"菜单标题"`
	Name     string `json:"name"        description:"菜单名称，用于路由name"`
	MenuType int    `json:"menu_type"   description:"菜单类型：1=目录 2=菜单 3=按钮"`
}

// MenuType 菜单类型常量
const (
	MenuTypeDir    = 1 // 目录
	MenuTypeMenu   = 2 // 菜单
	MenuTypeButton = 3 // 按钮
)

// MenuStatus 菜单状态常量
const (
	MenuStatusDisabled = 0 // 禁用
	MenuStatusEnabled  = 1 // 启用
	MenuStatusNormal   = 1 // 正常（别名）
)

// MenuVisible 菜单可见性常量
const (
	MenuHidden  = 0 // 隐藏
	MenuVisible = 1 // 显示
)

// MenuAlwaysShow 菜单总是显示常量
const (
	MenuAlwaysShowNo  = 0 // 否
	MenuAlwaysShowYes = 1 // 是
)

// MenuBreadcrumb 面包屑显示常量
const (
	MenuBreadcrumbHidden  = 0 // 隐藏
	MenuBreadcrumbVisible = 1 // 显示
)

// MenuTree 菜单树结构体，用于构建层级菜单
type MenuTree struct {
	Menu
	Children []*MenuTree `json:"children,omitempty" description:"子菜单"`
}

// IsDir 判断是否为目录
func (m *Menu) IsDir() bool {
	return m.MenuType == MenuTypeDir
}

// IsMenu 判断是否为菜单
func (m *Menu) IsMenu() bool {
	return m.MenuType == MenuTypeMenu
}

// IsButton 判断是否为按钮
func (m *Menu) IsButton() bool {
	return m.MenuType == MenuTypeButton
}

// IsEnabled 判断是否启用
func (m *Menu) IsEnabled() bool {
	return m.Status == MenuStatusEnabled
}

// IsVisible 判断是否可见
func (m *Menu) IsVisible() bool {
	return m.Visible == MenuVisible
}

// HasPermission 判断是否有权限标识
func (m *Menu) HasPermission() bool {
	return m.Permission != ""
}

// GetTypeName 获取菜单类型名称
func (m *Menu) GetTypeName() string {
	switch m.MenuType {
	case MenuTypeDir:
		return "目录"
	case MenuTypeMenu:
		return "菜单"
	case MenuTypeButton:
		return "按钮"
	default:
		return "未知"
	}
}

// GetStatusName 获取状态名称
func (m *Menu) GetStatusName() string {
	switch m.Status {
	case MenuStatusEnabled:
		return "启用"
	case MenuStatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}
