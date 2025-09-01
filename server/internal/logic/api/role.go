package api

import (
	"client-app/internal/model/entity"
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"
	"client-app/internal/service"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sRole struct{}

func NewRole() *sRole {
	return &sRole{}
}

func init() {
	service.RegisterRole(NewRole())
}

// GetRoleList 获取角色列表
func (s *sRole) GetRoleList(ctx context.Context, in *sysin.RoleListInp) (*sysout.RoleListModel, error) {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return nil, err
	}

	// 构建查询条件
	db := g.DB().Model("sys_roles").Where("deleted_at IS NULL")

	// 状态筛选
	if in.Status >= 0 {
		db = db.Where("status = ?", in.Status)
	}

	// 数据权限范围筛选
	if in.DataScope > 0 {
		db = db.Where("data_scope = ?", in.DataScope)
	}

	// 角色名称模糊查询
	if in.Name != "" {
		db = db.WhereLike("name", "%"+in.Name+"%")
	}

	// 角色编码模糊查询
	if in.Code != "" {
		db = db.WhereLike("code", "%"+in.Code+"%")
	}

	// 查询总数
	totalCount, err := db.Count()
	if err != nil {
		return nil, gerror.Newf("查询角色总数失败: %v", err)
	}

	// 如果没有数据，直接返回空结果
	if totalCount == 0 {
		return &sysout.RoleListModel{
			List:     []*sysout.RoleModel{},
			Total:    0,
			Page:     in.Page,
			PageSize: in.PageSize,
		}, nil
	}

	// 分页查询
	offset := (in.Page - 1) * in.PageSize
	db = db.Offset(offset).Limit(in.PageSize)

	// 排序
	orderBy := fmt.Sprintf("%s %s", in.OrderBy, in.OrderType)
	db = db.Order(orderBy)

	// 执行查询
	var roles []*entity.Role
	if err := db.Scan(&roles); err != nil {
		return nil, gerror.Newf("查询角色列表失败: %v", err)
	}

	// 转换为输出模型
	list := make([]*sysout.RoleModel, len(roles))
	for i, role := range roles {
		list[i] = sysout.ConvertToRoleModel(role)
	}

	return &sysout.RoleListModel{
		List:     list,
		Total:    int64(totalCount),
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}

// GetRoleDetail 获取角色详情
func (s *sRole) GetRoleDetail(ctx context.Context, in *sysin.RoleDetailInp) (*sysout.RoleDetailModel, error) {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return nil, err
	}

	// 查询角色信息
	var role *entity.Role
	err := g.DB().Model("sys_roles").Where("id = ? AND deleted_at IS NULL", in.Id).Scan(&role)
	if err != nil {
		return nil, gerror.Newf("查询角色详情失败: %v", err)
	}
	if role == nil {
		return nil, gerror.New("角色不存在或已删除")
	}

	// 查询角色拥有的菜单权限
	menuIds, err := s.getRoleMenuIds(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 查询权限标识列表
	permissions, err := s.getRolePermissions(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 转换为输出模型
	return sysout.ConvertToRoleDetailModel(role, menuIds, permissions), nil
}

// CreateRole 创建角色
func (s *sRole) CreateRole(ctx context.Context, in *sysin.CreateRoleInp) (*sysout.RoleModel, error) {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return nil, err
	}

	// 检查角色编码是否已存在
	exists, err := s.checkRoleCodeExists(ctx, in.Code, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, gerror.New("角色编码已存在")
	}

	// 检查角色名称是否已存在
	exists, err = s.checkRoleNameExists(ctx, in.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, gerror.New("角色名称已存在")
	}

	var resultRole *sysout.RoleModel

	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 插入角色记录
		roleData := &entity.Role{
			Name:        in.Name,
			Code:        in.Code,
			Description: in.Description,
			Status:      in.Status,
			Sort:        in.Sort,
			DataScope:   in.DataScope,
			Remark:      in.Remark,
			CreatedAt:   gtime.Now(),
			UpdatedAt:   gtime.Now(),
		}

		// 获取当前用户ID（从上下文或其他方式）
		if userId := s.getCurrentUserId(ctx); userId > 0 {
			roleData.CreatedBy = userId
			roleData.UpdatedBy = userId
		}

		result, err := tx.Model("sys_roles").Data(roleData).Insert()
		if err != nil {
			return gerror.Newf("创建角色失败: %v", err)
		}

		roleId, err := result.LastInsertId()
		if err != nil {
			return gerror.Newf("获取角色ID失败: %v", err)
		}

		roleData.Id = roleId

		// 分配菜单权限
		if len(in.MenuIds) > 0 {
			if err := s.assignRoleMenus(ctx, tx, roleId, in.MenuIds); err != nil {
				return err
			}
		}

		resultRole = sysout.ConvertToRoleModel(roleData)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resultRole, nil
}

// UpdateRole 更新角色
func (s *sRole) UpdateRole(ctx context.Context, in *sysin.UpdateRoleInp) (*sysout.RoleModel, error) {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return nil, err
	}

	// 检查角色是否存在
	exists, err := s.checkRoleExists(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, gerror.New("角色不存在")
	}

	// 检查是否为内置角色（内置角色的编码不能修改）
	role, err := s.getRoleById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if role.IsBuiltIn() && role.Code != in.Code {
		return nil, gerror.New("内置角色编码不允许修改")
	}

	// 检查角色编码是否已存在（排除自己）
	exists, err = s.checkRoleCodeExists(ctx, in.Code, in.Id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, gerror.New("角色编码已存在")
	}

	// 检查角色名称是否已存在（排除自己）
	exists, err = s.checkRoleNameExists(ctx, in.Name, in.Id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, gerror.New("角色名称已存在")
	}

	var resultRole *sysout.RoleModel

	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 更新角色记录
		updateData := g.Map{
			"name":        in.Name,
			"code":        in.Code,
			"description": in.Description,
			"status":      in.Status,
			"sort":        in.Sort,
			"data_scope":  in.DataScope,
			"remark":      in.Remark,
			"updated_at":  gtime.Now(),
		}

		// 获取当前用户ID
		if userId := s.getCurrentUserId(ctx); userId > 0 {
			updateData["updated_by"] = userId
		}

		_, err := tx.Model("sys_roles").Where("id = ?", in.Id).Data(updateData).Update()
		if err != nil {
			return gerror.Newf("更新角色失败: %v", err)
		}

		// 更新菜单权限
		if err := s.updateRoleMenus(ctx, tx, in.Id, in.MenuIds); err != nil {
			return err
		}

		// 查询更新后的角色信息
		var updatedRole *entity.Role
		err = tx.Model("sys_roles").Where("id = ?", in.Id).Scan(&updatedRole)
		if err != nil {
			return gerror.Newf("查询更新后角色信息失败: %v", err)
		}

		resultRole = sysout.ConvertToRoleModel(updatedRole)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resultRole, nil
}

// DeleteRole 删除角色
func (s *sRole) DeleteRole(ctx context.Context, in *sysin.DeleteRoleInp) error {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return err
	}

	// 检查角色是否存在
	role, err := s.getRoleById(ctx, in.Id)
	if err != nil {
		return err
	}

	// 检查是否为内置角色
	if role.IsBuiltIn() {
		return gerror.New("内置角色不允许删除")
	}

	// 检查是否有用户正在使用该角色
	hasUsers, err := s.checkRoleHasUsers(ctx, in.Id)
	if err != nil {
		return err
	}
	if hasUsers {
		return gerror.New("该角色正在被用户使用，无法删除")
	}

	// 开启事务删除
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 软删除角色
		updateData := g.Map{
			"deleted_at": gtime.Now(),
			"updated_at": gtime.Now(),
		}

		if userId := s.getCurrentUserId(ctx); userId > 0 {
			updateData["updated_by"] = userId
		}

		_, err := tx.Model("sys_roles").Where("id = ?", in.Id).Data(updateData).Update()
		if err != nil {
			return gerror.Newf("删除角色失败: %v", err)
		}

		// 删除角色菜单关联
		_, err = tx.Model("sys_role_menus").Where("role_id = ?", in.Id).Delete()
		if err != nil {
			return gerror.Newf("删除角色菜单关联失败: %v", err)
		}

		return nil
	})
}

// BatchDeleteRole 批量删除角色
func (s *sRole) BatchDeleteRole(ctx context.Context, in *sysin.BatchDeleteRoleInp) error {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return err
	}

	// 检查所有角色
	for _, roleId := range in.Ids {
		role, err := s.getRoleById(ctx, roleId)
		if err != nil {
			return gerror.Newf("角色ID %d 不存在", roleId)
		}

		if role.IsBuiltIn() {
			return gerror.Newf("角色 %s 为内置角色，不允许删除", role.Name)
		}

		hasUsers, err := s.checkRoleHasUsers(ctx, roleId)
		if err != nil {
			return err
		}
		if hasUsers {
			return gerror.Newf("角色 %s 正在被用户使用，无法删除", role.Name)
		}
	}

	// 开启事务批量删除
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 批量软删除角色
		updateData := g.Map{
			"deleted_at": gtime.Now(),
			"updated_at": gtime.Now(),
		}

		if userId := s.getCurrentUserId(ctx); userId > 0 {
			updateData["updated_by"] = userId
		}

		_, err := tx.Model("sys_roles").Where("id IN(?)", in.Ids).Data(updateData).Update()
		if err != nil {
			return gerror.Newf("批量删除角色失败: %v", err)
		}

		// 删除角色菜单关联
		_, err = tx.Model("sys_role_menus").Where("role_id IN(?)", in.Ids).Delete()
		if err != nil {
			return gerror.Newf("删除角色菜单关联失败: %v", err)
		}

		return nil
	})
}

// UpdateRoleStatus 更新角色状态
func (s *sRole) UpdateRoleStatus(ctx context.Context, in *sysin.UpdateRoleStatusInp) error {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return err
	}

	// 检查角色是否存在
	exists, err := s.checkRoleExists(ctx, in.Id)
	if err != nil {
		return err
	}
	if !exists {
		return gerror.New("角色不存在")
	}

	// 更新状态
	updateData := g.Map{
		"status":     in.Status,
		"updated_at": gtime.Now(),
	}

	if userId := s.getCurrentUserId(ctx); userId > 0 {
		updateData["updated_by"] = userId
	}

	_, err = g.DB().Model("sys_roles").Where("id = ?", in.Id).Data(updateData).Update()
	if err != nil {
		return gerror.Newf("更新角色状态失败: %v", err)
	}

	return nil
}

// CopyRole 复制角色
func (s *sRole) CopyRole(ctx context.Context, in *sysin.CopyRoleInp) (*sysout.RoleModel, error) {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return nil, err
	}

	// 获取源角色信息
	sourceRole, err := s.getRoleById(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 检查新角色编码是否已存在
	exists, err := s.checkRoleCodeExists(ctx, in.Code, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, gerror.New("角色编码已存在")
	}

	// 检查新角色名称是否已存在
	exists, err = s.checkRoleNameExists(ctx, in.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, gerror.New("角色名称已存在")
	}

	// 获取源角色的菜单权限
	menuIds, err := s.getRoleMenuIds(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	var resultRole *sysout.RoleModel

	// 开启事务复制
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 创建新角色
		newRole := &entity.Role{
			Name:        in.Name,
			Code:        in.Code,
			Description: sourceRole.Description,
			Status:      sourceRole.Status,
			Sort:        sourceRole.Sort,
			DataScope:   sourceRole.DataScope,
			Remark:      "从角色 " + sourceRole.Name + " 复制",
			CreatedAt:   gtime.Now(),
			UpdatedAt:   gtime.Now(),
		}

		if userId := s.getCurrentUserId(ctx); userId > 0 {
			newRole.CreatedBy = userId
			newRole.UpdatedBy = userId
		}

		result, err := tx.Model("sys_roles").Data(newRole).Insert()
		if err != nil {
			return gerror.Newf("复制角色失败: %v", err)
		}

		newRoleId, err := result.LastInsertId()
		if err != nil {
			return gerror.Newf("获取新角色ID失败: %v", err)
		}

		newRole.Id = newRoleId

		// 复制菜单权限
		if len(menuIds) > 0 {
			if err := s.assignRoleMenus(ctx, tx, newRoleId, menuIds); err != nil {
				return err
			}
		}

		resultRole = sysout.ConvertToRoleModel(newRole)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resultRole, nil
}

// GetRoleMenus 获取角色菜单权限
func (s *sRole) GetRoleMenus(ctx context.Context, in *sysin.RoleMenuInp) (*sysout.RoleMenuModel, error) {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return nil, err
	}

	// 检查角色是否存在
	role, err := s.getRoleById(ctx, in.RoleId)
	if err != nil {
		return nil, err
	}

	// 获取角色菜单权限
	menuIds, err := s.getRoleMenuIds(ctx, in.RoleId)
	if err != nil {
		return nil, err
	}

	// 获取权限标识列表
	permissions, err := s.getRolePermissions(ctx, in.RoleId)
	if err != nil {
		return nil, err
	}

	return &sysout.RoleMenuModel{
		RoleId:      role.Id,
		RoleName:    role.Name,
		RoleCode:    role.Code,
		MenuIds:     menuIds,
		Permissions: permissions,
	}, nil
}

// UpdateRoleMenus 更新角色菜单权限
func (s *sRole) UpdateRoleMenus(ctx context.Context, in *sysin.UpdateRoleMenuInp) error {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return err
	}

	// 检查角色是否存在
	exists, err := s.checkRoleExists(ctx, in.RoleId)
	if err != nil {
		return err
	}
	if !exists {
		return gerror.New("角色不存在")
	}

	// 开启事务更新
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return s.updateRoleMenus(ctx, tx, in.RoleId, in.MenuIds)
	})
}

// GetRolePermissions 获取角色权限详情
func (s *sRole) GetRolePermissions(ctx context.Context, in *sysin.RolePermissionInp) (*sysout.RolePermissionModel, error) {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return nil, err
	}

	// 获取角色信息
	role, err := s.getRoleById(ctx, in.RoleId)
	if err != nil {
		return nil, err
	}

	// 获取菜单权限
	menuIds, err := s.getRoleMenuIds(ctx, in.RoleId)
	if err != nil {
		return nil, err
	}

	// 获取权限标识
	permissions, err := s.getRolePermissions(ctx, in.RoleId)
	if err != nil {
		return nil, err
	}

	return sysout.ConvertToRolePermissionModel(role, menuIds, permissions), nil
}

// GetRoleOptions 获取角色选项
func (s *sRole) GetRoleOptions(ctx context.Context, in *sysin.RoleOptionInp) ([]*sysout.RoleOptionModel, error) {
	// 参数过滤
	if err := in.Filter(ctx); err != nil {
		return nil, err
	}

	// 构建查询条件
	db := g.DB().Model("sys_roles").Where("deleted_at IS NULL")

	if in.Status >= 0 {
		db = db.Where("status = ?", in.Status)
	}

	// 查询角色列表
	var roles []*entity.Role
	if err := db.Order("sort ASC, id ASC").Scan(&roles); err != nil {
		return nil, gerror.Newf("查询角色选项失败: %v", err)
	}

	// 转换为选项模型
	options := make([]*sysout.RoleOptionModel, len(roles))
	for i, role := range roles {
		options[i] = sysout.ConvertToRoleOptionModel(role)
	}

	return options, nil
}

// GetRoleStats 获取角色统计
func (s *sRole) GetRoleStats(ctx context.Context) (*sysout.RoleStatsModel, error) {
	// 查询所有角色
	var roles []*entity.Role
	err := g.DB().Model("sys_roles").Where("deleted_at IS NULL").Scan(&roles)
	if err != nil {
		return nil, gerror.Newf("查询角色统计失败: %v", err)
	}

	return sysout.BuildRoleStatsModel(roles), nil
}

// GetDataScopeOptions 获取数据权限范围选项
func (s *sRole) GetDataScopeOptions(ctx context.Context) ([]*sysout.DataScopeModel, error) {
	return sysout.GetAllDataScopes(), nil
}

// 下面是辅助方法的实现

// getRoleById 根据ID获取角色
func (s *sRole) getRoleById(ctx context.Context, roleId int64) (*entity.Role, error) {
	var role *entity.Role
	err := g.DB().Model("sys_roles").Where("id = ? AND deleted_at IS NULL", roleId).Scan(&role)
	if err != nil {
		return nil, gerror.Newf("查询角色失败: %v", err)
	}
	if role == nil {
		return nil, gerror.New("角色不存在")
	}
	return role, nil
}

// checkRoleExists 检查角色是否存在
func (s *sRole) checkRoleExists(ctx context.Context, roleId int64) (bool, error) {
	count, err := g.DB().Model("sys_roles").Where("id = ? AND deleted_at IS NULL", roleId).Count()
	if err != nil {
		return false, gerror.Newf("检查角色存在性失败: %v", err)
	}
	return count > 0, nil
}

// checkRoleCodeExists 检查角色编码是否存在
func (s *sRole) checkRoleCodeExists(ctx context.Context, code string, excludeId int64) (bool, error) {
	db := g.DB().Model("sys_roles").Where("code = ? AND deleted_at IS NULL", code)
	if excludeId > 0 {
		db = db.Where("id != ?", excludeId)
	}

	count, err := db.Count()
	if err != nil {
		return false, gerror.Newf("检查角色编码失败: %v", err)
	}
	return count > 0, nil
}

// checkRoleNameExists 检查角色名称是否存在
func (s *sRole) checkRoleNameExists(ctx context.Context, name string, excludeId int64) (bool, error) {
	db := g.DB().Model("sys_roles").Where("name = ? AND deleted_at IS NULL", name)
	if excludeId > 0 {
		db = db.Where("id != ?", excludeId)
	}

	count, err := db.Count()
	if err != nil {
		return false, gerror.Newf("检查角色名称失败: %v", err)
	}
	return count > 0, nil
}

// checkRoleHasUsers 检查角色是否有用户使用
func (s *sRole) checkRoleHasUsers(ctx context.Context, roleId int64) (bool, error) {
	count, err := g.DB().Model("sys_user_roles").Where("role_id = ?", roleId).Count()
	if err != nil {
		return false, gerror.Newf("检查角色用户关联失败: %v", err)
	}
	return count > 0, nil
}

// getRoleMenuIds 获取角色的菜单ID列表
func (s *sRole) getRoleMenuIds(ctx context.Context, roleId int64) ([]int64, error) {
	var menuIds []int64
	result, err := g.DB().Model("sys_role_menus").Fields("menu_id").Where("role_id = ?", roleId).Array()
	if err != nil {
		return nil, gerror.Newf("查询角色菜单权限失败: %v", err)
	}

	// 转换结果
	for _, v := range result {
		if id := gconv.Int64(v); id > 0 {
			menuIds = append(menuIds, id)
		}
	}

	return menuIds, nil
}

// getRolePermissions 获取角色的权限标识列表
func (s *sRole) getRolePermissions(ctx context.Context, roleId int64) ([]string, error) {
	var permissions []string

	sql := `SELECT DISTINCT m.permission 
			FROM sys_role_menus rm 
			JOIN menus m ON rm.menu_id = m.id 
			WHERE rm.role_id = ? AND m.status = 1 AND m.permission != ''`

	result, err := g.DB().Raw(sql, roleId).Array()
	if err != nil {
		return nil, gerror.Newf("查询角色权限标识失败: %v", err)
	}

	// 转换结果
	for _, v := range result {
		if perm := gconv.String(v); perm != "" {
			permissions = append(permissions, perm)
		}
	}

	return permissions, nil
}

// assignRoleMenus 分配角色菜单权限
func (s *sRole) assignRoleMenus(ctx context.Context, tx gdb.TX, roleId int64, menuIds []int64) error {
	if len(menuIds) == 0 {
		return nil
	}

	// 批量插入角色菜单关联
	var data []g.Map
	for _, menuId := range menuIds {
		data = append(data, g.Map{
			"role_id":    roleId,
			"menu_id":    menuId,
			"created_at": gtime.Now(),
		})
	}

	_, err := tx.Model("sys_role_menus").Data(data).Insert()
	if err != nil {
		return gerror.Newf("分配角色菜单权限失败: %v", err)
	}

	return nil
}

// updateRoleMenus 更新角色菜单权限
func (s *sRole) updateRoleMenus(ctx context.Context, tx gdb.TX, roleId int64, menuIds []int64) error {
	// 先删除现有权限
	_, err := tx.Model("sys_role_menus").Where("role_id = ?", roleId).Delete()
	if err != nil {
		return gerror.Newf("删除角色原有菜单权限失败: %v", err)
	}

	// 重新分配权限
	if len(menuIds) > 0 {
		if err := s.assignRoleMenus(ctx, tx, roleId, menuIds); err != nil {
			return err
		}
	}

	return nil
}

// getCurrentUserId 获取当前用户ID（这里需要根据实际的认证机制实现）
func (s *sRole) getCurrentUserId(ctx context.Context) int64 {
	// 这里应该从JWT token或session中获取当前用户ID
	// 暂时返回1作为默认值
	return 1
}

// AssignUserRoles 分配用户角色
func (s *sRole) AssignUserRoles(ctx context.Context, userId int64, roleIds []int64, assignedBy int64) error {
	if len(roleIds) == 0 {
		return gerror.New("角色ID列表不能为空")
	}

	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 批量插入用户角色关联
		var data []g.Map
		for i, roleId := range roleIds {
			data = append(data, g.Map{
				"user_id":     userId,
				"role_id":     roleId,
				"is_primary":  gconv.Int(i == 0), // 第一个角色设为主要角色
				"assigned_by": assignedBy,
				"created_at":  gtime.Now(),
				"updated_at":  gtime.Now(),
			})
		}

		_, err := tx.Model("sys_user_roles").Data(data).Insert()
		if err != nil {
			return gerror.Newf("分配用户角色失败: %v", err)
		}

		return nil
	})
}

// RemoveUserRoles 移除用户角色
func (s *sRole) RemoveUserRoles(ctx context.Context, userId int64, roleIds []int64) error {
	if len(roleIds) == 0 {
		return gerror.New("角色ID列表不能为空")
	}

	_, err := g.DB().Model("sys_user_roles").Where("user_id = ? AND role_id IN(?)", userId, roleIds).Delete()
	if err != nil {
		return gerror.Newf("移除用户角色失败: %v", err)
	}

	return nil
}

// GetUserRoles 获取用户角色列表
func (s *sRole) GetUserRoles(ctx context.Context, userId int64) ([]*sysout.RoleModel, error) {
	var roles []*entity.Role

	sql := `SELECT r.* FROM roles r 
			JOIN user_roles ur ON r.id = ur.role_id 
			WHERE ur.user_id = ? AND r.deleted_at IS NULL 
			ORDER BY ur.is_primary DESC, r.sort ASC`

	err := g.DB().Raw(sql, userId).Scan(&roles)
	if err != nil {
		return nil, gerror.Newf("查询用户角色失败: %v", err)
	}

	// 转换为输出模型
	result := make([]*sysout.RoleModel, len(roles))
	for i, role := range roles {
		result[i] = sysout.ConvertToRoleModel(role)
	}

	return result, nil
}

// SetUserPrimaryRole 设置用户主要角色
func (s *sRole) SetUserPrimaryRole(ctx context.Context, userId int64, roleId int64) error {
	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先将所有角色设为非主要角色
		_, err := tx.Model("sys_user_roles").Where("user_id = ?", userId).Data(g.Map{
			"is_primary": 0,
			"updated_at": gtime.Now(),
		}).Update()
		if err != nil {
			return gerror.Newf("更新用户角色失败: %v", err)
		}

		// 设置指定角色为主要角色
		_, err = tx.Model("sys_user_roles").Where("user_id = ? AND role_id = ?", userId, roleId).Data(g.Map{
			"is_primary": 1,
			"updated_at": gtime.Now(),
		}).Update()
		if err != nil {
			return gerror.Newf("设置主要角色失败: %v", err)
		}

		return nil
	})
}

// CheckUserPermission 检查用户权限
func (s *sRole) CheckUserPermission(ctx context.Context, userId int64, permission string) (bool, error) {
	sql := `SELECT COUNT(*) FROM sys_user_roles ur
			JOIN sys_role_menus rm ON ur.role_id = rm.role_id
			JOIN sys_menus m ON rm.menu_id = m.id
			JOIN sys_roles r ON ur.role_id = r.id
			WHERE ur.user_id = ? AND m.permission = ? 
			AND r.status = 1 AND m.status = 1 
			AND r.deleted_at IS NULL`

	count, err := g.DB().Raw(sql, userId, permission).Count()
	if err != nil {
		return false, gerror.Newf("检查用户权限失败: %v", err)
	}

	return count > 0, nil
}

// CheckUserRole 检查用户角色
func (s *sRole) CheckUserRole(ctx context.Context, userId int64, roleCode string) (bool, error) {
	sql := `SELECT COUNT(*) FROM sys_user_roles ur
			JOIN sys_roles r ON ur.role_id = r.id
			WHERE ur.user_id = ? AND r.code = ? 
			AND r.status = 1 AND r.deleted_at IS NULL`

	count, err := g.DB().Raw(sql, userId, roleCode).Count()
	if err != nil {
		return false, gerror.Newf("检查用户角色失败: %v", err)
	}

	return count > 0, nil
}

// GetUserPermissions 获取用户权限列表
func (s *sRole) GetUserPermissions(ctx context.Context, userId int64) ([]string, error) {
	var permissions []string

	sql := `SELECT DISTINCT m.permission FROM sys_user_roles ur
			JOIN sys_role_menus rm ON ur.role_id = rm.role_id
			JOIN sys_menus m ON rm.menu_id = m.id
			JOIN sys_roles r ON ur.role_id = r.id
			WHERE ur.user_id = ? AND m.permission != ''
			AND r.status = 1 AND m.status = 1 
			AND r.deleted_at IS NULL`

	result, err := g.DB().Raw(sql, userId).Array()
	if err != nil {
		return nil, gerror.Newf("获取用户权限列表失败: %v", err)
	}

	// 转换结果
	for _, v := range result {
		if perm := gconv.String(v); perm != "" {
			permissions = append(permissions, perm)
		}
	}

	return permissions, nil
}

// GetUserMenus 获取用户菜单ID列表
func (s *sRole) GetUserMenus(ctx context.Context, userId int64) ([]int64, error) {
	var menuIds []int64

	sql := `SELECT DISTINCT rm.menu_id FROM sys_user_roles ur
			JOIN sys_role_menus rm ON ur.role_id = rm.role_id
			JOIN sys_roles r ON ur.role_id = r.id
			WHERE ur.user_id = ? AND r.status = 1 AND r.deleted_at IS NULL`

	result, err := g.DB().Raw(sql, userId).Array()
	if err != nil {
		return nil, gerror.Newf("获取用户菜单列表失败: %v", err)
	}

	// 转换结果
	for _, v := range result {
		if id := gconv.Int64(v); id > 0 {
			menuIds = append(menuIds, id)
		}
	}

	return menuIds, nil
}

// GetUserDataScope 获取用户数据权限范围
func (s *sRole) GetUserDataScope(ctx context.Context, userId int64) (int, error) {
	sql := `SELECT MIN(r.data_scope) FROM sys_user_roles ur
			JOIN sys_roles r ON ur.role_id = r.id
			WHERE ur.user_id = ? AND r.status = 1 AND r.deleted_at IS NULL`

	var dataScope int
	err := g.DB().Raw(sql, userId).Scan(&dataScope)
	if err != nil {
		return entity.DataScopeSelf, gerror.Newf("获取用户数据权限范围失败: %v", err)
	}

	// 如果没有角色，返回最小权限
	if dataScope == 0 {
		dataScope = entity.DataScopeSelf
	}

	return dataScope, nil
}

// CheckUsersPermission 批量检查用户权限
func (s *sRole) CheckUsersPermission(ctx context.Context, userIds []int64, permission string) (map[int64]bool, error) {
	if len(userIds) == 0 {
		return make(map[int64]bool), nil
	}

	sql := `SELECT DISTINCT ur.user_id FROM sys_user_roles ur
			JOIN sys_role_menus rm ON ur.role_id = rm.role_id
			JOIN sys_menus m ON rm.menu_id = m.id
			JOIN sys_roles r ON ur.role_id = r.id
			WHERE ur.user_id IN(?) AND m.permission = ? 
			AND r.status = 1 AND m.status = 1 
			AND r.deleted_at IS NULL`

	result, err := g.DB().Raw(sql, userIds, permission).Array()
	if err != nil {
		return nil, gerror.Newf("批量检查用户权限失败: %v", err)
	}

	// 转换结果
	var hasPermissionUserIds []int64
	for _, v := range result {
		if id := gconv.Int64(v); id > 0 {
			hasPermissionUserIds = append(hasPermissionUserIds, id)
		}
	}

	// 构建结果map
	resultMap := make(map[int64]bool)
	hasPermissionMap := make(map[int64]bool)

	for _, userId := range hasPermissionUserIds {
		hasPermissionMap[userId] = true
	}

	for _, userId := range userIds {
		resultMap[userId] = hasPermissionMap[userId]
	}

	return resultMap, nil
}

// FilterUsersByPermission 根据权限过滤用户
func (s *sRole) FilterUsersByPermission(ctx context.Context, userIds []int64, permission string) ([]int64, error) {
	if len(userIds) == 0 {
		return []int64{}, nil
	}

	sql := `SELECT DISTINCT ur.user_id FROM sys_user_roles ur
			JOIN sys_role_menus rm ON ur.role_id = rm.role_id
			JOIN sys_menus m ON rm.menu_id = m.id
			JOIN sys_roles r ON ur.role_id = r.id
			WHERE ur.user_id IN(?) AND m.permission = ? 
			AND r.status = 1 AND m.status = 1 
			AND r.deleted_at IS NULL`

	queryResult, err := g.DB().Raw(sql, userIds, permission).Array()
	if err != nil {
		return nil, gerror.Newf("根据权限过滤用户失败: %v", err)
	}

	// 转换结果
	var result []int64
	for _, v := range queryResult {
		if id := gconv.Int64(v); id > 0 {
			result = append(result, id)
		}
	}

	return result, nil
}
