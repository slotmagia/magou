package api

import (
	"context"
	"fmt"

	"client-app/internal/consts"
	"client-app/internal/model/entity"
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"
	"client-app/internal/service"
	"client-app/utility/logger"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sMenu struct{}

func init() {
	service.RegisterMenu(NewMenu())
}

func NewMenu() *sMenu {
	return &sMenu{}
}

// GetMenuList 获取菜单列表
func (s *sMenu) GetMenuList(ctx context.Context, in *sysin.MenuListInp) (res *sysout.MenuListModel, err error) {
	// 构建查询条件
	m := g.DB().Model("sys_menus")

	// 标题模糊查询
	if in.Title != "" {
		m = m.WhereLike("title", "%"+in.Title+"%")
	}

	// 名称模糊查询
	if in.Name != "" {
		m = m.WhereLike("name", "%"+in.Name+"%")
	}

	// 状态过滤
	if in.Status >= 0 {
		m = m.Where("status", in.Status)
	}

	// 类型过滤
	if in.Type > 0 {
		m = m.Where("type", in.Type)
	}

	// 父菜单过滤
	if in.ParentId >= 0 {
		m = m.Where("parent_id", in.ParentId)
	}

	// 统计总数
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Newf("查询菜单总数失败: %v", err)
	}

	// 排序
	orderBy := in.OrderBy
	if orderBy == "" {
		orderBy = "sort_order"
	}
	m = m.Order(fmt.Sprintf("%s %s", orderBy, in.OrderType))

	// 分页查询
	var menuEntities []*entity.Menu
	err = m.Page(in.Page, in.PageSize).Scan(&menuEntities)
	if err != nil {
		return nil, gerror.Newf("查询菜单列表失败: %v", err)
	}

	// 转换为输出模型
	list := make([]*sysout.MenuModel, 0, len(menuEntities))
	for _, menu := range menuEntities {
		list = append(list, sysout.ConvertToMenuModel(menu))
	}

	return &sysout.MenuListModel{
		List:     list,
		Total:    int64(total),
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}

// GetMenuTree 获取菜单树
func (s *sMenu) GetMenuTree(ctx context.Context, in *sysin.MenuTreeInp) (res []*sysout.MenuTreeModel, err error) {
	m := g.DB().Model("sys_menus")

	// 状态过滤
	if in.Status >= 0 {
		m = m.Where("status", in.Status)
	}

	// 类型过滤
	if in.Type > 0 {
		m = m.Where("type", in.Type)
	}

	// 查询所有菜单
	var menuEntities []*entity.Menu
	err = m.Order("sort ASC, id ASC").Scan(&menuEntities)
	if err != nil {
		return nil, gerror.Newf("查询菜单列表失败: %v", err)
	}

	// 转换为树形模型
	menuModels := make([]*sysout.MenuTreeModel, 0, len(menuEntities))
	for _, menu := range menuEntities {
		menuModels = append(menuModels, sysout.ConvertToMenuTreeModel(menu))
	}

	// 构建树结构
	return s.BuildMenuTree(ctx, menuModels, 0), nil
}

// GetMenuDetail 获取菜单详情
func (s *sMenu) GetMenuDetail(ctx context.Context, in *sysin.MenuDetailInp) (res *sysout.MenuDetailModel, err error) {
	var menuEntity *entity.Menu
	err = g.DB().Model("sys_menus").Where("id", in.Id).Scan(&menuEntity)
	if err != nil {
		return nil, gerror.Newf("查询菜单失败: %v", err)
	}

	if menuEntity == nil {
		return nil, gerror.New("菜单不存在")
	}

	// 查询子菜单
	var childEntities []*entity.Menu
	err = g.DB().Model("sys_menus").Where("parent_id", in.Id).Order("sort ASC").Scan(&childEntities)
	if err != nil {
		return nil, gerror.Newf("查询子菜单失败: %v", err)
	}

	// 构建详情模型
	detail := &sysout.MenuDetailModel{
		MenuModel: *sysout.ConvertToMenuModel(menuEntity),
		Children:  make([]*sysout.MenuDetailModel, 0, len(childEntities)),
	}

	// 添加子菜单
	for _, child := range childEntities {
		detail.Children = append(detail.Children, &sysout.MenuDetailModel{
			MenuModel: *sysout.ConvertToMenuModel(child),
		})
	}

	return detail, nil
}

// CreateMenu 创建菜单
func (s *sMenu) CreateMenu(ctx context.Context, in *sysin.CreateMenuInp) (res *sysout.MenuModel, err error) {
	// 验证父菜单
	if in.ParentId > 0 {
		exists, err := s.CheckMenuExists(ctx, in.ParentId)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, gerror.New("父菜单不存在")
		}
	}

	// 检查名称是否重复
	nameExists, err := s.CheckMenuNameExists(ctx, in.Name, 0)
	if err != nil {
		return nil, err
	}
	if nameExists {
		return nil, gerror.New("菜单名称已存在")
	}

	// 检查路径是否重复（非按钮类型）
	if in.Type != entity.MenuTypeButton && in.Path != "" {
		pathExists, err := s.CheckMenuPathExists(ctx, in.Path, 0)
		if err != nil {
			return nil, err
		}
		if pathExists {
			return nil, gerror.New("菜单路径已存在")
		}
	}

	// 获取当前用户ID
	userId := g.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId, 0).Int64()

	// 构建菜单实体
	menu := &entity.Menu{
		ParentId:   in.ParentId,
		Title:      in.Title,
		Name:       in.Name,
		Path:       in.Path,
		Component:  in.Component,
		Icon:       in.Icon,
		MenuType:   in.Type,
		SortOrder:  in.Sort,
		Status:     in.Status,
		Visible:    in.Visible,
		Permission: in.Permission,
		Redirect:   in.Redirect,
		AlwaysShow: in.AlwaysShow,
		Breadcrumb: in.Breadcrumb,
		ActiveMenu: in.ActiveMenu,
		Remark:     in.Remark,
		CreatedBy:  userId,
		UpdatedBy:  userId,
		CreatedAt:  gtime.Now(),
		UpdatedAt:  gtime.Now(),
	}

	// 插入数据库
	result, err := g.DB().Model("sys_menus").Data(menu).Insert()
	if err != nil {
		return nil, gerror.Newf("创建菜单失败: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Newf("获取菜单ID失败: %v", err)
	}

	menu.Id = id
	return sysout.ConvertToMenuModel(menu), nil
}

// UpdateMenu 更新菜单
func (s *sMenu) UpdateMenu(ctx context.Context, in *sysin.UpdateMenuInp) (res *sysout.MenuModel, err error) {
	// 检查菜单是否存在
	exists, err := s.CheckMenuExists(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, gerror.New("菜单不存在")
	}

	// 验证父菜单
	if in.ParentId > 0 {
		// 不能设置自己为父菜单
		if in.ParentId == in.Id {
			return nil, gerror.New("不能设置自己为父菜单")
		}

		parentExists, err := s.CheckMenuExists(ctx, in.ParentId)
		if err != nil {
			return nil, err
		}
		if !parentExists {
			return nil, gerror.New("父菜单不存在")
		}

		// 检查是否会形成循环依赖
		childIds, err := s.GetChildrenIds(ctx, in.Id)
		if err != nil {
			return nil, err
		}
		for _, childId := range childIds {
			if childId == in.ParentId {
				return nil, gerror.New("不能设置子菜单为父菜单，会形成循环依赖")
			}
		}
	}

	// 检查名称是否重复
	nameExists, err := s.CheckMenuNameExists(ctx, in.Name, in.Id)
	if err != nil {
		return nil, err
	}
	if nameExists {
		return nil, gerror.New("菜单名称已存在")
	}

	// 检查路径是否重复（非按钮类型）
	if in.Type != entity.MenuTypeButton && in.Path != "" {
		pathExists, err := s.CheckMenuPathExists(ctx, in.Path, in.Id)
		if err != nil {
			return nil, err
		}
		if pathExists {
			return nil, gerror.New("菜单路径已存在")
		}
	}

	// 获取当前用户ID
	userId := g.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId, 0).Int64()

	// 更新数据
	updateData := g.Map{
		"parent_id":   in.ParentId,
		"title":       in.Title,
		"name":        in.Name,
		"path":        in.Path,
		"component":   in.Component,
		"icon":        in.Icon,
		"type":        in.Type,
		"sort":        in.Sort,
		"status":      in.Status,
		"visible":     in.Visible,
		"permission":  in.Permission,
		"redirect":    in.Redirect,
		"always_show": in.AlwaysShow,
		"breadcrumb":  in.Breadcrumb,
		"active_menu": in.ActiveMenu,
		"remark":      in.Remark,
		"updated_by":  userId,
		"updated_at":  gtime.Now(),
	}

	_, err = g.DB().Model("sys_menus").Where("id", in.Id).Data(updateData).Update()
	if err != nil {
		return nil, gerror.Newf("更新菜单失败: %v", err)
	}

	// 查询更新后的菜单
	var menu *entity.Menu
	err = g.DB().Model("sys_menus").Where("id", in.Id).Scan(&menu)
	if err != nil {
		return nil, gerror.Newf("查询更新后菜单失败: %v", err)
	}

	return sysout.ConvertToMenuModel(menu), nil
}

// DeleteMenu 删除菜单
func (s *sMenu) DeleteMenu(ctx context.Context, in *sysin.DeleteMenuInp) (err error) {
	// 检查菜单是否存在
	exists, err := s.CheckMenuExists(ctx, in.Id)
	if err != nil {
		return err
	}
	if !exists {
		return gerror.New("菜单不存在")
	}

	// 检查是否有子菜单
	var count int
	count, err = g.DB().Model("menus").Where("parent_id", in.Id).Count()
	if err != nil {
		return gerror.Newf("检查子菜单失败: %v", err)
	}
	if count > 0 {
		return gerror.New("存在子菜单，无法删除")
	}

	// 检查是否有角色关联
	count, err = g.DB().Model("role_menus").Where("menu_id", in.Id).Count()
	if err != nil {
		return gerror.Newf("检查角色关联失败: %v", err)
	}
	if count > 0 {
		return gerror.New("菜单已被角色使用，无法删除")
	}

	// 删除菜单
	_, err = g.DB().Model("menus").Where("id", in.Id).Delete()
	if err != nil {
		return gerror.Newf("删除菜单失败: %v", err)
	}

	return nil
}

// BatchDeleteMenu 批量删除菜单
func (s *sMenu) BatchDeleteMenu(ctx context.Context, in *sysin.BatchDeleteMenuInp) (err error) {
	if len(in.Ids) == 0 {
		return gerror.New("请选择要删除的菜单")
	}

	// 事务处理
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, id := range in.Ids {
			// 检查是否有子菜单
			count, err := tx.Model("menus").Where("parent_id", id).Count()
			if err != nil {
				return gerror.Newf("检查菜单[%d]子菜单失败: %v", id, err)
			}
			if count > 0 {
				return gerror.Newf("菜单[%d]存在子菜单，无法删除", id)
			}

			// 检查是否有角色关联
			count, err = tx.Model("role_menus").Where("menu_id", id).Count()
			if err != nil {
				return gerror.Newf("检查菜单[%d]角色关联失败: %v", id, err)
			}
			if count > 0 {
				return gerror.Newf("菜单[%d]已被角色使用，无法删除", id)
			}

			// 删除菜单
			_, err = tx.Model("menus").Where("id", id).Delete()
			if err != nil {
				return gerror.Newf("删除菜单[%d]失败: %v", id, err)
			}
		}
		return nil
	})
}

// UpdateMenuStatus 更新菜单状态
func (s *sMenu) UpdateMenuStatus(ctx context.Context, in *sysin.UpdateMenuStatusInp) (err error) {
	// 检查菜单是否存在
	exists, err := s.CheckMenuExists(ctx, in.Id)
	if err != nil {
		return err
	}
	if !exists {
		return gerror.New("菜单不存在")
	}

	// 获取当前用户ID
	userId := g.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId, 0).Int64()

	// 更新状态
	_, err = g.DB().Model("menus").Where("id", in.Id).Data(g.Map{
		"status":     in.Status,
		"updated_by": userId,
		"updated_at": gtime.Now(),
	}).Update()

	if err != nil {
		return gerror.Newf("更新菜单状态失败: %v", err)
	}

	return nil
}

// GetMenuOptions 获取菜单选项
func (s *sMenu) GetMenuOptions(ctx context.Context, in *sysin.MenuOptionInp) (res []*sysout.MenuOptionModel, err error) {
	m := g.DB().Model("sys_menus")

	// 状态过滤
	if in.Status >= 0 {
		m = m.Where("status", in.Status)
	} else {
		// 默认只显示启用的菜单
		m = m.Where("status", entity.MenuStatusEnabled)
	}

	// 类型过滤
	if in.Type > 0 {
		m = m.Where("type", in.Type)
	} else if in.ParentOnly {
		// 只返回可作为父菜单的选项（目录和菜单）
		m = m.WhereIn("type", []int{entity.MenuTypeDir, entity.MenuTypeMenu})
	}

	// 排除指定菜单
	if in.ExcludeId > 0 {
		m = m.WhereNot("id", in.ExcludeId)
	}

	// 查询菜单列表
	var menuEntities []*entity.Menu
	err = m.Order("sort ASC, id ASC").Scan(&menuEntities)
	if err != nil {
		return nil, gerror.Newf("查询菜单选项失败: %v", err)
	}

	// 转换为选项模型
	options := make([]*sysout.MenuOptionModel, 0, len(menuEntities))
	for _, menu := range menuEntities {
		options = append(options, sysout.ConvertToMenuOptionModel(menu))
	}

	// 构建树结构
	return s.buildMenuOptionTree(options, 0), nil
}

// GetRouters 获取前端路由列表
func (s *sMenu) GetRouters(ctx context.Context) (res []*sysout.RouterModel, err error) {
	// 查询启用的菜单
	var menuEntities []*entity.Menu
	start := gtime.Now()
	err = g.DB().Model("sys_menus").
		Where("status", entity.MenuStatusEnabled).
		WhereIn("menu_type", []int{entity.MenuTypeDir, entity.MenuTypeMenu}).
		Order("sort_order ASC, id ASC").
		Scan(&menuEntities)
	// 打印可点击的调用位置，便于从控制台跳转到此SQL生成处
	logger.LogSQLWithCaller(ctx, "(SQL见上一条Gf调试输出)", nil, gtime.Now().Sub(start))

	if err != nil {
		return nil, gerror.Newf("查询路由列表失败: %v", err)
	}

	// 转换为路由模型
	routers := make([]*sysout.RouterModel, 0, len(menuEntities))
	for _, menu := range menuEntities {
		routers = append(routers, sysout.ConvertToRouterModel(menu))
	}

	// 构建树结构
	return s.BuildRouterTree(ctx, routers, 0), nil
}

// ValidateMenuPermission 验证菜单权限
func (s *sMenu) ValidateMenuPermission(ctx context.Context, userId int64, permission string) (bool, error) {
	// 通过角色服务验证权限
	return service.Role().CheckUserPermission(ctx, userId, permission)
}

// GetUserMenus 获取用户可访问的菜单列表
func (s *sMenu) GetUserMenus(ctx context.Context, userId int64) (res []*sysout.MenuTreeModel, err error) {
	// 获取用户菜单ID列表
	menuIds, err := service.Role().GetUserMenus(ctx, userId)
	if err != nil {
		return nil, err
	}

	if len(menuIds) == 0 {
		return []*sysout.MenuTreeModel{}, nil
	}

	// 查询用户可访问的菜单
	var menuEntities []*entity.Menu
	err = g.DB().Model("sys_menus").
		Where("status", entity.MenuStatusEnabled).
		Where("visible", entity.MenuVisible).
		WhereIn("id", menuIds).
		Order("sort ASC, id ASC").
		Scan(&menuEntities)

	if err != nil {
		return nil, gerror.Newf("查询用户菜单失败: %v", err)
	}

	// 转换为树形模型
	menuModels := make([]*sysout.MenuTreeModel, 0, len(menuEntities))
	for _, menu := range menuEntities {
		menuModels = append(menuModels, sysout.ConvertToMenuTreeModel(menu))
	}

	// 构建树结构
	return s.BuildMenuTree(ctx, menuModels, 0), nil
}

// GetUserRouters 获取用户可访问的前端路由
func (s *sMenu) GetUserRouters(ctx context.Context, userId int64) (res []*sysout.RouterModel, err error) {
	// 获取用户菜单ID列表
	menuIds, err := service.Role().GetUserMenus(ctx, userId)
	if err != nil {
		return nil, err
	}

	if len(menuIds) == 0 {
		return []*sysout.RouterModel{}, nil
	}

	// 查询用户可访问的路由菜单
	var menuEntities []*entity.Menu
	err = g.DB().Model("sys_menus").
		Where("status", entity.MenuStatusEnabled).
		WhereIn("type", []int{entity.MenuTypeDir, entity.MenuTypeMenu}).
		WhereIn("id", menuIds).
		Order("sort ASC, id ASC").
		Scan(&menuEntities)

	if err != nil {
		return nil, gerror.Newf("查询用户路由失败: %v", err)
	}

	// 转换为路由模型
	routers := make([]*sysout.RouterModel, 0, len(menuEntities))
	for _, menu := range menuEntities {
		routers = append(routers, sysout.ConvertToRouterModel(menu))
	}

	// 构建树结构
	return s.BuildRouterTree(ctx, routers, 0), nil
}

// BuildMenuTree 构建菜单树结构
func (s *sMenu) BuildMenuTree(ctx context.Context, menus []*sysout.MenuTreeModel, parentId int64) []*sysout.MenuTreeModel {
	var tree []*sysout.MenuTreeModel

	for _, menu := range menus {
		if menu.ParentId == parentId {
			// 递归构建子菜单
			children := s.BuildMenuTree(ctx, menus, menu.Id)
			if len(children) > 0 {
				menu.Children = children
			}
			tree = append(tree, menu)
		}
	}

	return tree
}

// BuildRouterTree 构建路由树结构
func (s *sMenu) BuildRouterTree(ctx context.Context, routers []*sysout.RouterModel, parentId int64) []*sysout.RouterModel {
	var tree []*sysout.RouterModel

	for _, router := range routers {
		// 通过Path判断父子关系，或者可以添加ParentId字段
		if s.isChildRouter(router, parentId) {
			// 递归构建子路由
			children := s.BuildRouterTree(ctx, routers, router.Id)
			if len(children) > 0 {
				router.Children = children
			}
			tree = append(tree, router)
		}
	}

	return tree
}

// CheckMenuExists 检查菜单是否存在
func (s *sMenu) CheckMenuExists(ctx context.Context, id int64) (bool, error) {
	count, err := g.DB().Model("sys_menus").Where("id", id).Count()
	if err != nil {
		return false, gerror.Newf("检查菜单是否存在失败: %v", err)
	}
	return count > 0, nil
}

// CheckMenuNameExists 检查菜单名称是否存在
func (s *sMenu) CheckMenuNameExists(ctx context.Context, name string, excludeId int64) (bool, error) {
	m := g.DB().Model("sys_menus").Where("name", name)
	if excludeId > 0 {
		m = m.WhereNot("id", excludeId)
	}

	count, err := m.Count()
	if err != nil {
		return false, gerror.Newf("检查菜单名称是否存在失败: %v", err)
	}
	return count > 0, nil
}

// CheckMenuPathExists 检查菜单路径是否存在
func (s *sMenu) CheckMenuPathExists(ctx context.Context, path string, excludeId int64) (bool, error) {
	m := g.DB().Model("sys_menus").Where("path", path)
	if excludeId > 0 {
		m = m.WhereNot("id", excludeId)
	}

	count, err := m.Count()
	if err != nil {
		return false, gerror.Newf("检查菜单路径是否存在失败: %v", err)
	}
	return count > 0, nil
}

// GetChildrenIds 获取指定菜单的所有子菜单ID（递归）
func (s *sMenu) GetChildrenIds(ctx context.Context, parentId int64) ([]int64, error) {
	var childIds []int64

	// 查询直接子菜单
	var directChildren []int64
	_, err := g.DB().Model("sys_menus").Fields("id").Where("parent_id", parentId).Array(&directChildren)
	if err != nil {
		return nil, gerror.Newf("查询子菜单失败: %v", err)
	}

	// 递归查询每个子菜单的子菜单
	for _, childId := range directChildren {
		childIds = append(childIds, childId)

		// 递归获取子菜单的子菜单
		grandChildIds, err := s.GetChildrenIds(ctx, childId)
		if err != nil {
			return nil, err
		}
		childIds = append(childIds, grandChildIds...)
	}

	return childIds, nil
}

// IsParentMenu 判断菜单是否可以作为父菜单
func (s *sMenu) IsParentMenu(ctx context.Context, menuType int) bool {
	// 只有目录和菜单可以作为父菜单，按钮不能
	return menuType == entity.MenuTypeDir || menuType == entity.MenuTypeMenu
}

// buildMenuOptionTree 构建菜单选项树结构
func (s *sMenu) buildMenuOptionTree(options []*sysout.MenuOptionModel, parentId int64) []*sysout.MenuOptionModel {
	var tree []*sysout.MenuOptionModel

	for _, option := range options {
		// 通过菜单实体查询父子关系
		var menu *entity.Menu
		g.DB().Model("sys_menus").Where("id", option.Value).Scan(&menu)

		if menu != nil && menu.ParentId == parentId {
			// 递归构建子选项
			children := s.buildMenuOptionTree(options, option.Value)
			if len(children) > 0 {
				option.Children = children
			}
			tree = append(tree, option)
		}
	}

	return tree
}

// isChildRouter 判断路由是否为指定父路由的子路由
func (s *sMenu) isChildRouter(router *sysout.RouterModel, parentId int64) bool {
	// 这里需要根据实际业务逻辑判断
	// 可以通过数据库查询菜单的parent_id字段
	var menu *entity.Menu
	err := g.DB().Model("sys_menus").Where("id", router.Id).Scan(&menu)
	if err != nil || menu == nil {
		return false
	}
	return menu.ParentId == parentId
}
