package service

import (
	"context"

	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"
)

// IMenu 菜单服务接口
type IMenu interface {
	// GetMenuList 获取菜单列表
	GetMenuList(ctx context.Context, in *sysin.MenuListInp) (res *sysout.MenuListModel, err error)

	// GetMenuTree 获取菜单树
	GetMenuTree(ctx context.Context, in *sysin.MenuTreeInp) (res []*sysout.MenuTreeModel, err error)

	// GetMenuDetail 获取菜单详情
	GetMenuDetail(ctx context.Context, in *sysin.MenuDetailInp) (res *sysout.MenuDetailModel, err error)

	// CreateMenu 创建菜单
	CreateMenu(ctx context.Context, in *sysin.CreateMenuInp) (res *sysout.MenuModel, err error)

	// UpdateMenu 更新菜单
	UpdateMenu(ctx context.Context, in *sysin.UpdateMenuInp) (res *sysout.MenuModel, err error)

	// DeleteMenu 删除菜单
	DeleteMenu(ctx context.Context, in *sysin.DeleteMenuInp) (err error)

	// BatchDeleteMenu 批量删除菜单
	BatchDeleteMenu(ctx context.Context, in *sysin.BatchDeleteMenuInp) (err error)

	// UpdateMenuStatus 更新菜单状态
	UpdateMenuStatus(ctx context.Context, in *sysin.UpdateMenuStatusInp) (err error)

	// GetMenuOptions 获取菜单选项（用于下拉框等）
	GetMenuOptions(ctx context.Context, in *sysin.MenuOptionInp) (res []*sysout.MenuOptionModel, err error)

	// GetRouters 获取前端路由列表
	GetRouters(ctx context.Context) (res []*sysout.RouterModel, err error)

	// ValidateMenuPermission 验证菜单权限
	ValidateMenuPermission(ctx context.Context, userId int64, permission string) (bool, error)

	// GetUserMenus 获取用户可访问的菜单列表
	GetUserMenus(ctx context.Context, userId int64) (res []*sysout.MenuTreeModel, err error)

	// GetUserRouters 获取用户可访问的前端路由
	GetUserRouters(ctx context.Context, userId int64) (res []*sysout.RouterModel, err error)

	// BuildMenuTree 构建菜单树结构
	BuildMenuTree(ctx context.Context, menus []*sysout.MenuTreeModel, parentId int64) []*sysout.MenuTreeModel

	// BuildRouterTree 构建路由树结构
	BuildRouterTree(ctx context.Context, routers []*sysout.RouterModel, parentId int64) []*sysout.RouterModel

	// CheckMenuExists 检查菜单是否存在
	CheckMenuExists(ctx context.Context, id int64) (bool, error)

	// CheckMenuNameExists 检查菜单名称是否存在
	CheckMenuNameExists(ctx context.Context, name string, excludeId int64) (bool, error)

	// CheckMenuPathExists 检查菜单路径是否存在
	CheckMenuPathExists(ctx context.Context, path string, excludeId int64) (bool, error)

	// GetChildrenIds 获取指定菜单的所有子菜单ID（递归）
	GetChildrenIds(ctx context.Context, parentId int64) ([]int64, error)

	// IsParentMenu 判断菜单是否可以作为父菜单
	IsParentMenu(ctx context.Context, menuType int) bool
}

var localMenu IMenu

// Menu 获取菜单服务实例
func Menu() IMenu {
	if localMenu == nil {
		panic("implement not found for interface IMenu, forgot register?")
	}
	return localMenu
}

// RegisterMenu 注册菜单服务实现
func RegisterMenu(i IMenu) {
	localMenu = i
}
