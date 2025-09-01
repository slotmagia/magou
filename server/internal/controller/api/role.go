package api

import (
	role "client-app/internal/api/v1/role"
	"client-app/internal/service"
	"context"
)

var (
	Role = cRole{}
)

type cRole struct{}

// GetRoleList 获取角色列表
func (c *cRole) GetRoleList(ctx context.Context, req *role.RoleListReq) (res *role.RoleListRes, err error) {
	// 记录业务日志
	service.Middleware().LogBusiness(ctx, "GetRoleList", req, "开始")

	out, err := service.Role().GetRoleList(ctx, &req.RoleListInp)
	if err != nil {
		service.Middleware().LogError(ctx, err, "获取角色列表失败")
		return nil, err
	}

	service.Middleware().LogBusiness(ctx, "GetRoleList", out, "成功")
	return &role.RoleListRes{
		RoleListModel: out,
	}, nil
}

// GetRoleDetail 获取角色详情
func (c *cRole) GetRoleDetail(ctx context.Context, req *role.RoleDetailReq) (res *role.RoleDetailRes, err error) {
	out, err := service.Role().GetRoleDetail(ctx, &req.RoleDetailInp)
	if err != nil {
		return nil, err
	}

	return &role.RoleDetailRes{
		RoleDetailModel: out,
	}, nil
}

// CreateRole 创建角色
func (c *cRole) CreateRole(ctx context.Context, req *role.CreateRoleReq) (res *role.CreateRoleRes, err error) {
	// 记录业务日志
	service.Middleware().LogBusiness(ctx, "CreateRole", req, "开始")

	out, err := service.Role().CreateRole(ctx, &req.CreateRoleInp)
	if err != nil {
		service.Middleware().LogError(ctx, err, "创建角色失败")
		return nil, err
	}

	service.Middleware().LogBusiness(ctx, "CreateRole", out, "成功")
	service.Middleware().LogAudit(ctx, "CREATE", "ROLE", "SUCCESS", "创建角色", out.Id)
	return &role.CreateRoleRes{
		RoleModel: out,
	}, nil
}

// UpdateRole 更新角色
func (c *cRole) UpdateRole(ctx context.Context, req *role.UpdateRoleReq) (res *role.UpdateRoleRes, err error) {
	out, err := service.Role().UpdateRole(ctx, &req.UpdateRoleInp)
	if err != nil {
		return nil, err
	}

	return &role.UpdateRoleRes{
		RoleModel: out,
	}, nil
}

// DeleteRole 删除角色
func (c *cRole) DeleteRole(ctx context.Context, req *role.DeleteRoleReq) (res *role.DeleteRoleRes, err error) {
	// 记录业务日志
	service.Middleware().LogBusiness(ctx, "DeleteRole", req, "开始")

	err = service.Role().DeleteRole(ctx, &req.DeleteRoleInp)
	if err != nil {
		service.Middleware().LogError(ctx, err, "删除角色失败")
		service.Middleware().LogAudit(ctx, "DELETE", "ROLE", "FAILED", "删除角色失败", req.Id)
		return &role.DeleteRoleRes{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	service.Middleware().LogBusiness(ctx, "DeleteRole", req.Id, "成功")
	service.Middleware().LogAudit(ctx, "DELETE", "ROLE", "SUCCESS", "删除角色", req.Id)
	service.Middleware().LogSecurity(ctx, "ROLE_DELETED", "medium", "角色被删除", req.Id)
	return &role.DeleteRoleRes{
		Success: true,
		Message: "删除成功",
	}, nil
}

// BatchDeleteRole 批量删除角色
func (c *cRole) BatchDeleteRole(ctx context.Context, req *role.BatchDeleteRoleReq) (res *role.BatchDeleteRoleRes, err error) {
	err = service.Role().BatchDeleteRole(ctx, &req.BatchDeleteRoleInp)
	if err != nil {
		return &role.BatchDeleteRoleRes{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &role.BatchDeleteRoleRes{
		Success: true,
		Message: "批量删除成功",
	}, nil
}

// UpdateRoleStatus 更新角色状态
func (c *cRole) UpdateRoleStatus(ctx context.Context, req *role.UpdateRoleStatusReq) (res *role.UpdateRoleStatusRes, err error) {
	err = service.Role().UpdateRoleStatus(ctx, &req.UpdateRoleStatusInp)
	if err != nil {
		return &role.UpdateRoleStatusRes{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &role.UpdateRoleStatusRes{
		Success: true,
		Message: "状态更新成功",
	}, nil
}

// CopyRole 复制角色
func (c *cRole) CopyRole(ctx context.Context, req *role.CopyRoleReq) (res *role.CopyRoleRes, err error) {
	out, err := service.Role().CopyRole(ctx, &req.CopyRoleInp)
	if err != nil {
		return nil, err
	}

	return &role.CopyRoleRes{
		RoleModel: out,
	}, nil
}

// GetRoleMenus 获取角色菜单权限
func (c *cRole) GetRoleMenus(ctx context.Context, req *role.RoleMenuReq) (res *role.RoleMenuRes, err error) {
	out, err := service.Role().GetRoleMenus(ctx, &req.RoleMenuInp)
	if err != nil {
		return nil, err
	}

	return &role.RoleMenuRes{
		RoleMenuModel: out,
	}, nil
}

// UpdateRoleMenus 更新角色菜单权限
func (c *cRole) UpdateRoleMenus(ctx context.Context, req *role.UpdateRoleMenuReq) (res *role.UpdateRoleMenuRes, err error) {
	err = service.Role().UpdateRoleMenus(ctx, &req.UpdateRoleMenuInp)
	if err != nil {
		return &role.UpdateRoleMenuRes{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &role.UpdateRoleMenuRes{
		Success: true,
		Message: "权限更新成功",
	}, nil
}

// GetRolePermissions 获取角色权限详情
func (c *cRole) GetRolePermissions(ctx context.Context, req *role.RolePermissionReq) (res *role.RolePermissionRes, err error) {
	out, err := service.Role().GetRolePermissions(ctx, &req.RolePermissionInp)
	if err != nil {
		return nil, err
	}

	return &role.RolePermissionRes{
		RolePermissionModel: out,
	}, nil
}

// GetRoleOptions 获取角色选项
func (c *cRole) GetRoleOptions(ctx context.Context, req *role.RoleOptionReq) (res *role.RoleOptionRes, err error) {
	out, err := service.Role().GetRoleOptions(ctx, &req.RoleOptionInp)
	if err != nil {
		return nil, err
	}

	return &role.RoleOptionRes{
		List: out,
	}, nil
}

// GetRoleStats 获取角色统计
func (c *cRole) GetRoleStats(ctx context.Context, req *role.RoleStatsReq) (res *role.RoleStatsRes, err error) {
	out, err := service.Role().GetRoleStats(ctx)
	if err != nil {
		return nil, err
	}

	return &role.RoleStatsRes{
		RoleStatsModel: out,
	}, nil
}

// GetDataScopeOptions 获取数据权限范围选项
func (c *cRole) GetDataScopeOptions(ctx context.Context, req *role.DataScopeOptionReq) (res *role.DataScopeOptionRes, err error) {
	out, err := service.Role().GetDataScopeOptions(ctx)
	if err != nil {
		return nil, err
	}

	return &role.DataScopeOptionRes{
		List: out,
	}, nil
}
