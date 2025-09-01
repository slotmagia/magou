package api

import (
	"context"

	v1 "client-app/internal/api/v1/menu"
	"client-app/internal/service"
)

var (
	Menu = cMenu{}
)

// cMenu 菜单控制器
type cMenu struct{}

// GetMenuList 获取菜单列表
func (c *cMenu) GetMenuList(ctx context.Context, req *v1.MenuListReq) (res *v1.MenuListRes, err error) {
	out, err := service.Menu().GetMenuList(ctx, &req.MenuListInp)
	if err != nil {
		return nil, err
	}
	return &v1.MenuListRes{MenuListModel: out}, nil
}

// GetMenuTree 获取菜单树
func (c *cMenu) GetMenuTree(ctx context.Context, req *v1.MenuTreeReq) (res *v1.MenuTreeRes, err error) {
	out, err := service.Menu().GetMenuTree(ctx, &req.MenuTreeInp)
	if err != nil {
		return nil, err
	}
	return &v1.MenuTreeRes{List: out}, nil
}

// GetMenuDetail 获取菜单详情
func (c *cMenu) GetMenuDetail(ctx context.Context, req *v1.MenuDetailReq) (res *v1.MenuDetailRes, err error) {
	out, err := service.Menu().GetMenuDetail(ctx, &req.MenuDetailInp)
	if err != nil {
		return nil, err
	}
	return &v1.MenuDetailRes{MenuDetailModel: out}, nil
}

// CreateMenu 创建菜单
func (c *cMenu) CreateMenu(ctx context.Context, req *v1.CreateMenuReq) (res *v1.CreateMenuRes, err error) {
	out, err := service.Menu().CreateMenu(ctx, &req.CreateMenuInp)
	if err != nil {
		return nil, err
	}
	return &v1.CreateMenuRes{MenuModel: out}, nil
}

// UpdateMenu 更新菜单
func (c *cMenu) UpdateMenu(ctx context.Context, req *v1.UpdateMenuReq) (res *v1.UpdateMenuRes, err error) {
	out, err := service.Menu().UpdateMenu(ctx, &req.UpdateMenuInp)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateMenuRes{MenuModel: out}, nil
}

// DeleteMenu 删除菜单
func (c *cMenu) DeleteMenu(ctx context.Context, req *v1.DeleteMenuReq) (res *v1.DeleteMenuRes, err error) {
	err = service.Menu().DeleteMenu(ctx, &req.DeleteMenuInp)
	if err != nil {
		return &v1.DeleteMenuRes{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &v1.DeleteMenuRes{
		Success: true,
		Message: "删除成功",
	}, nil
}

// BatchDeleteMenu 批量删除菜单
func (c *cMenu) BatchDeleteMenu(ctx context.Context, req *v1.BatchDeleteMenuReq) (res *v1.BatchDeleteMenuRes, err error) {
	err = service.Menu().BatchDeleteMenu(ctx, &req.BatchDeleteMenuInp)
	if err != nil {
		return &v1.BatchDeleteMenuRes{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &v1.BatchDeleteMenuRes{
		Success: true,
		Message: "批量删除成功",
	}, nil
}

// UpdateMenuStatus 更新菜单状态
func (c *cMenu) UpdateMenuStatus(ctx context.Context, req *v1.UpdateMenuStatusReq) (res *v1.UpdateMenuStatusRes, err error) {
	err = service.Menu().UpdateMenuStatus(ctx, &req.UpdateMenuStatusInp)
	if err != nil {
		return &v1.UpdateMenuStatusRes{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &v1.UpdateMenuStatusRes{
		Success: true,
		Message: "状态更新成功",
	}, nil
}

// GetMenuOptions 获取菜单选项
func (c *cMenu) GetMenuOptions(ctx context.Context, req *v1.MenuOptionsReq) (res *v1.MenuOptionsRes, err error) {
	out, err := service.Menu().GetMenuOptions(ctx, &req.MenuOptionInp)
	if err != nil {
		return nil, err
	}
	return &v1.MenuOptionsRes{List: out}, nil
}

// GetRouters 获取前端路由
func (c *cMenu) GetRouters(ctx context.Context, req *v1.RoutersReq) (res *v1.RoutersRes, err error) {
	out, err := service.Menu().GetRouters(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.RoutersRes{List: out}, nil
}
