package sysout

import (
	"client-app/internal/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
)

// MenuListModel 菜单列表响应模型
type MenuListModel struct {
	List     []*MenuModel `json:"list" description:"菜单列表"`
	Total    int64        `json:"total" description:"总记录数"`
	Page     int          `json:"page" description:"当前页码"`
	PageSize int          `json:"pageSize" description:"每页数量"`
}

// MenuModel 菜单基础响应模型
type MenuModel struct {
	Id         int64       `json:"id" description:"主键ID"`
	ParentId   int64       `json:"parentId" description:"父菜单ID"`
	Title      string      `json:"title" description:"菜单标题"`
	Name       string      `json:"name" description:"菜单名称"`
	Path       string      `json:"path" description:"菜单路径"`
	Component  string      `json:"component" description:"组件路径"`
	Icon       string      `json:"icon" description:"菜单图标"`
	Type       int         `json:"type" description:"菜单类型"`
	TypeName   string      `json:"typeName" description:"菜单类型名称"`
	Sort       int         `json:"sort" description:"排序号"`
	Status     int         `json:"status" description:"状态"`
	StatusName string      `json:"statusName" description:"状态名称"`
	Visible    int         `json:"visible" description:"是否显示"`
	Permission string      `json:"permission" description:"权限标识"`
	Redirect   string      `json:"redirect" description:"重定向地址"`
	AlwaysShow int         `json:"alwaysShow" description:"是否总是显示"`
	Breadcrumb int         `json:"breadcrumb" description:"是否显示面包屑"`
	ActiveMenu string      `json:"activeMenu" description:"高亮菜单路径"`
	Remark     string      `json:"remark" description:"备注说明"`
	CreatedBy  int64       `json:"createdBy" description:"创建人ID"`
	UpdatedBy  int64       `json:"updatedBy" description:"修改人ID"`
	CreatedAt  *gtime.Time `json:"createdAt" description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updatedAt" description:"更新时间"`
}

// MenuDetailModel 菜单详情响应模型
type MenuDetailModel struct {
	MenuModel
	Children []*MenuDetailModel `json:"children,omitempty" description:"子菜单"`
}

// MenuTreeModel 菜单树响应模型
type MenuTreeModel struct {
	Id         int64            `json:"id" description:"主键ID"`
	ParentId   int64            `json:"parentId" description:"父菜单ID"`
	Title      string           `json:"title" description:"菜单标题"`
	Name       string           `json:"name" description:"菜单名称"`
	Path       string           `json:"path" description:"菜单路径"`
	Component  string           `json:"component" description:"组件路径"`
	Icon       string           `json:"icon" description:"菜单图标"`
	Type       int              `json:"type" description:"菜单类型"`
	Sort       int              `json:"sort" description:"排序号"`
	Status     int              `json:"status" description:"状态"`
	Visible    int              `json:"visible" description:"是否显示"`
	Permission string           `json:"permission" description:"权限标识"`
	Redirect   string           `json:"redirect" description:"重定向地址"`
	AlwaysShow int              `json:"alwaysShow" description:"是否总是显示"`
	Breadcrumb int              `json:"breadcrumb" description:"是否显示面包屑"`
	ActiveMenu string           `json:"activeMenu" description:"高亮菜单路径"`
	Children   []*MenuTreeModel `json:"children,omitempty" description:"子菜单"`
}

// RouterModel 前端路由模型
type RouterModel struct {
	Id         int64          `json:"id" description:"菜单ID"`
	Name       string         `json:"name" description:"路由名称"`
	Path       string         `json:"path" description:"路由路径"`
	Component  string         `json:"component" description:"组件路径"`
	Redirect   string         `json:"redirect,omitempty" description:"重定向地址"`
	AlwaysShow bool           `json:"alwaysShow" description:"是否总是显示"`
	Hidden     bool           `json:"hidden" description:"是否隐藏"`
	Meta       *RouterMeta    `json:"meta" description:"路由元信息"`
	Children   []*RouterModel `json:"children,omitempty" description:"子路由"`
}

// RouterMeta 路由元信息
type RouterMeta struct {
	Title       string `json:"title" description:"菜单标题"`
	Icon        string `json:"icon,omitempty" description:"菜单图标"`
	NoCache     bool   `json:"noCache" description:"是否缓存"`
	Breadcrumb  bool   `json:"breadcrumb" description:"是否显示面包屑"`
	ActiveMenu  string `json:"activeMenu,omitempty" description:"高亮菜单路径"`
	Permissions string `json:"permissions,omitempty" description:"权限标识"`
}

// MenuOptionModel 菜单选项模型（用于下拉框等）
type MenuOptionModel struct {
	Value    int64              `json:"value" description:"菜单ID"`
	Label    string             `json:"label" description:"菜单标题"`
	Type     int                `json:"type" description:"菜单类型"`
	Disabled bool               `json:"disabled" description:"是否禁用"`
	Children []*MenuOptionModel `json:"children,omitempty" description:"子菜单"`
}

// ConvertToMenuModel 将entity.Menu转换为MenuModel
func ConvertToMenuModel(menu *entity.Menu) *MenuModel {
	if menu == nil {
		return nil
	}

	return &MenuModel{
		Id:         menu.Id,
		ParentId:   menu.ParentId,
		Title:      menu.Title,
		Name:       menu.Name,
		Path:       menu.Path,
		Component:  menu.Component,
		Icon:       menu.Icon,
		Type:       menu.MenuType,
		TypeName:   menu.GetTypeName(),
		Sort:       menu.SortOrder,
		Status:     menu.Status,
		StatusName: menu.GetStatusName(),
		Visible:    menu.Visible,
		Permission: menu.Permission,
		Redirect:   menu.Redirect,
		AlwaysShow: menu.AlwaysShow,
		Breadcrumb: menu.Breadcrumb,
		ActiveMenu: menu.ActiveMenu,
		Remark:     menu.Remark,
		CreatedBy:  menu.CreatedBy,
		UpdatedBy:  menu.UpdatedBy,
		CreatedAt:  menu.CreatedAt,
		UpdatedAt:  menu.UpdatedAt,
	}
}

// ConvertToMenuTreeModel 将entity.Menu转换为MenuTreeModel
func ConvertToMenuTreeModel(menu *entity.Menu) *MenuTreeModel {
	if menu == nil {
		return nil
	}

	return &MenuTreeModel{
		Id:         menu.Id,
		ParentId:   menu.ParentId,
		Title:      menu.Title,
		Name:       menu.Name,
		Path:       menu.Path,
		Component:  menu.Component,
		Icon:       menu.Icon,
		Type:       menu.MenuType,
		Sort:       menu.SortOrder,
		Status:     menu.Status,
		Visible:    menu.Visible,
		Permission: menu.Permission,
		Redirect:   menu.Redirect,
		AlwaysShow: menu.AlwaysShow,
		Breadcrumb: menu.Breadcrumb,
		ActiveMenu: menu.ActiveMenu,
		Children:   make([]*MenuTreeModel, 0),
	}
}

// ConvertToRouterModel 将entity.Menu转换为RouterModel
func ConvertToRouterModel(menu *entity.Menu) *RouterModel {
	if menu == nil {
		return nil
	}

	router := &RouterModel{
		Id:         menu.Id,
		Name:       menu.Title,
		Path:       menu.Path,
		Component:  menu.Component,
		AlwaysShow: menu.AlwaysShow == entity.MenuAlwaysShowYes,
		Hidden:     menu.Visible == entity.MenuHidden,
		Meta: &RouterMeta{
			Title:       menu.Title,
			Icon:        menu.Icon,
			NoCache:     false,
			Breadcrumb:  menu.Breadcrumb == entity.MenuBreadcrumbVisible,
			ActiveMenu:  menu.ActiveMenu,
			Permissions: menu.Permission,
		},
		Children: make([]*RouterModel, 0),
	}

	if menu.Redirect != "" {
		router.Redirect = menu.Redirect
	}

	return router
}

// ConvertToMenuOptionModel 将entity.Menu转换为MenuOptionModel
func ConvertToMenuOptionModel(menu *entity.Menu) *MenuOptionModel {
	if menu == nil {
		return nil
	}

	return &MenuOptionModel{
		Value:    menu.Id,
		Label:    menu.Title,
		Type:     menu.MenuType,
		Disabled: menu.Status == entity.MenuStatusDisabled,
		Children: make([]*MenuOptionModel, 0),
	}
}
