package service

import (
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"
	"context"
)

type IRole interface {
	// 角色基础操作
	GetRoleList(ctx context.Context, in *sysin.RoleListInp) (res *sysout.RoleListModel, err error)
	GetRoleDetail(ctx context.Context, in *sysin.RoleDetailInp) (res *sysout.RoleDetailModel, err error)
	CreateRole(ctx context.Context, in *sysin.CreateRoleInp) (res *sysout.RoleModel, err error)
	UpdateRole(ctx context.Context, in *sysin.UpdateRoleInp) (res *sysout.RoleModel, err error)
	DeleteRole(ctx context.Context, in *sysin.DeleteRoleInp) (err error)
	BatchDeleteRole(ctx context.Context, in *sysin.BatchDeleteRoleInp) (err error)
	UpdateRoleStatus(ctx context.Context, in *sysin.UpdateRoleStatusInp) (err error)
	CopyRole(ctx context.Context, in *sysin.CopyRoleInp) (res *sysout.RoleModel, err error)

	// 角色权限管理
	GetRoleMenus(ctx context.Context, in *sysin.RoleMenuInp) (res *sysout.RoleMenuModel, err error)
	UpdateRoleMenus(ctx context.Context, in *sysin.UpdateRoleMenuInp) (err error)
	GetRolePermissions(ctx context.Context, in *sysin.RolePermissionInp) (res *sysout.RolePermissionModel, err error)

	// 角色选项和统计
	GetRoleOptions(ctx context.Context, in *sysin.RoleOptionInp) (res []*sysout.RoleOptionModel, err error)
	GetRoleStats(ctx context.Context) (res *sysout.RoleStatsModel, err error)
	GetDataScopeOptions(ctx context.Context) (res []*sysout.DataScopeModel, err error)

	// 用户角色关联操作
	AssignUserRoles(ctx context.Context, userId int64, roleIds []int64, assignedBy int64) (err error)
	RemoveUserRoles(ctx context.Context, userId int64, roleIds []int64) (err error)
	GetUserRoles(ctx context.Context, userId int64) (res []*sysout.RoleModel, err error)
	SetUserPrimaryRole(ctx context.Context, userId int64, roleId int64) (err error)

	// 权限验证
	CheckUserPermission(ctx context.Context, userId int64, permission string) (bool, error)
	CheckUserRole(ctx context.Context, userId int64, roleCode string) (bool, error)
	GetUserPermissions(ctx context.Context, userId int64) ([]string, error)
	GetUserMenus(ctx context.Context, userId int64) ([]int64, error)
	GetUserDataScope(ctx context.Context, userId int64) (int, error)

	// 批量权限验证
	CheckUsersPermission(ctx context.Context, userIds []int64, permission string) (map[int64]bool, error)
	FilterUsersByPermission(ctx context.Context, userIds []int64, permission string) ([]int64, error)
}

var (
	localRole IRole
)

func Role() IRole {
	if localRole == nil {
		panic("implement not found for interface IRole, forgot register?")
	}
	return localRole
}

func RegisterRole(i IRole) {
	localRole = i
}
