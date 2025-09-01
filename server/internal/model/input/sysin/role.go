package sysin

import (
	"context"
	"strings"
)

// RoleListInp 角色列表查询参数
type RoleListInp struct {
	Name      string `json:"name" v:""`      // 角色名称（模糊查询）
	Code      string `json:"code" v:""`      // 角色编码（模糊查询）
	Status    int    `json:"status" v:""`    // 状态：1=启用 0=禁用，-1=全部
	DataScope int    `json:"dataScope" v:""` // 数据权限范围
	Page      int    `json:"page" v:"min:1"` // 页码
	PageSize  int    `json:"pageSize" v:""`  // 每页数量
	OrderBy   string `json:"orderBy" v:""`   // 排序字段
	OrderType string `json:"orderType" v:""` // 排序方式：asc/desc
}

// Filter 过滤输入参数
func (in *RoleListInp) Filter(ctx context.Context) (err error) {
	// 设置默认分页
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 20
	}
	if in.PageSize > 100 {
		in.PageSize = 100 // 限制最大每页数量
	}

	// 设置默认排序
	if in.OrderBy == "" {
		in.OrderBy = "sort"
	}
	if in.OrderType == "" {
		in.OrderType = "asc"
	}

	// 验证排序方式
	in.OrderType = strings.ToLower(in.OrderType)
	if in.OrderType != "asc" && in.OrderType != "desc" {
		in.OrderType = "asc"
	}

	// 去除前后空格
	in.Name = strings.TrimSpace(in.Name)
	in.Code = strings.TrimSpace(in.Code)

	return nil
}

// CreateRoleInp 创建角色参数
type CreateRoleInp struct {
	Name        string  `json:"name" v:"required|length:1,50#角色名称不能为空|角色名称长度不能超过50个字符"`
	Code        string  `json:"code" v:"required|length:1,50#角色编码不能为空|角色编码长度不能超过50个字符"`
	Description string  `json:"description" v:"length:0,200#角色描述长度不能超过200个字符"`
	Status      int     `json:"status" v:"in:0,1#状态必须是0(禁用)或1(启用)"`
	Sort        int     `json:"sort" v:"min:0#排序号不能小于0"`
	DataScope   int     `json:"dataScope" v:"required|in:1,2,3,4,5#数据权限范围不能为空|数据权限范围必须是1-5之间的数字"`
	Remark      string  `json:"remark" v:"length:0,500#备注说明长度不能超过500个字符"`
	MenuIds     []int64 `json:"menuIds" v:""` // 菜单权限ID列表
}

// Filter 过滤输入参数
func (in *CreateRoleInp) Filter(ctx context.Context) (err error) {
	// 去除前后空格
	in.Name = strings.TrimSpace(in.Name)
	in.Code = strings.TrimSpace(in.Code)
	in.Description = strings.TrimSpace(in.Description)
	in.Remark = strings.TrimSpace(in.Remark)

	// 角色编码转小写
	in.Code = strings.ToLower(in.Code)

	// 设置默认值
	if in.Status == 0 {
		in.Status = 1 // 默认启用
	}
	if in.DataScope == 0 {
		in.DataScope = 4 // 默认仅本人数据
	}

	// 去重菜单ID
	if len(in.MenuIds) > 0 {
		menuIdMap := make(map[int64]bool)
		uniqueMenuIds := make([]int64, 0)
		for _, menuId := range in.MenuIds {
			if menuId > 0 && !menuIdMap[menuId] {
				menuIdMap[menuId] = true
				uniqueMenuIds = append(uniqueMenuIds, menuId)
			}
		}
		in.MenuIds = uniqueMenuIds
	}

	return nil
}

// UpdateRoleInp 更新角色参数
type UpdateRoleInp struct {
	Id          int64   `json:"id" v:"required|min:1#角色ID不能为空|角色ID必须大于0"`
	Name        string  `json:"name" v:"required|length:1,50#角色名称不能为空|角色名称长度不能超过50个字符"`
	Code        string  `json:"code" v:"required|length:1,50#角色编码不能为空|角色编码长度不能超过50个字符"`
	Description string  `json:"description" v:"length:0,200#角色描述长度不能超过200个字符"`
	Status      int     `json:"status" v:"in:0,1#状态必须是0(禁用)或1(启用)"`
	Sort        int     `json:"sort" v:"min:0#排序号不能小于0"`
	DataScope   int     `json:"dataScope" v:"required|in:1,2,3,4,5#数据权限范围不能为空|数据权限范围必须是1-5之间的数字"`
	Remark      string  `json:"remark" v:"length:0,500#备注说明长度不能超过500个字符"`
	MenuIds     []int64 `json:"menuIds" v:""` // 菜单权限ID列表
}

// Filter 过滤输入参数
func (in *UpdateRoleInp) Filter(ctx context.Context) (err error) {
	// 去除前后空格
	in.Name = strings.TrimSpace(in.Name)
	in.Code = strings.TrimSpace(in.Code)
	in.Description = strings.TrimSpace(in.Description)
	in.Remark = strings.TrimSpace(in.Remark)

	// 角色编码转小写
	in.Code = strings.ToLower(in.Code)

	// 去重菜单ID
	if len(in.MenuIds) > 0 {
		menuIdMap := make(map[int64]bool)
		uniqueMenuIds := make([]int64, 0)
		for _, menuId := range in.MenuIds {
			if menuId > 0 && !menuIdMap[menuId] {
				menuIdMap[menuId] = true
				uniqueMenuIds = append(uniqueMenuIds, menuId)
			}
		}
		in.MenuIds = uniqueMenuIds
	}

	return nil
}

// DeleteRoleInp 删除角色参数
type DeleteRoleInp struct {
	Id int64 `json:"id" v:"required|min:1#角色ID不能为空|角色ID必须大于0"`
}

// Filter 过滤输入参数
func (in *DeleteRoleInp) Filter(ctx context.Context) (err error) {
	return nil
}

// RoleDetailInp 角色详情查询参数
type RoleDetailInp struct {
	Id int64 `json:"id" v:"required|min:1#角色ID不能为空|角色ID必须大于0"`
}

// Filter 过滤输入参数
func (in *RoleDetailInp) Filter(ctx context.Context) (err error) {
	return nil
}

// UpdateRoleStatusInp 更新角色状态参数
type UpdateRoleStatusInp struct {
	Id     int64 `json:"id" v:"required|min:1#角色ID不能为空|角色ID必须大于0"`
	Status int   `json:"status" v:"required|in:0,1#状态不能为空|状态必须是0(禁用)或1(启用)"`
}

// Filter 过滤输入参数
func (in *UpdateRoleStatusInp) Filter(ctx context.Context) (err error) {
	return nil
}

// RoleMenuInp 角色菜单权限查询参数
type RoleMenuInp struct {
	RoleId int64 `json:"roleId" v:"required|min:1#角色ID不能为空|角色ID必须大于0"`
}

// Filter 过滤输入参数
func (in *RoleMenuInp) Filter(ctx context.Context) (err error) {
	return nil
}

// UpdateRoleMenuInp 更新角色菜单权限参数
type UpdateRoleMenuInp struct {
	RoleId  int64   `json:"roleId" v:"required|min:1#角色ID不能为空|角色ID必须大于0"`
	MenuIds []int64 `json:"menuIds" v:""` // 菜单权限ID列表
}

// Filter 过滤输入参数
func (in *UpdateRoleMenuInp) Filter(ctx context.Context) (err error) {
	// 去重菜单ID
	if len(in.MenuIds) > 0 {
		menuIdMap := make(map[int64]bool)
		uniqueMenuIds := make([]int64, 0)
		for _, menuId := range in.MenuIds {
			if menuId > 0 && !menuIdMap[menuId] {
				menuIdMap[menuId] = true
				uniqueMenuIds = append(uniqueMenuIds, menuId)
			}
		}
		in.MenuIds = uniqueMenuIds
	}

	return nil
}

// RoleOptionInp 角色选项查询参数（用于下拉框等）
type RoleOptionInp struct {
	Status int `json:"status" v:""` // 状态过滤：1=启用 0=禁用，-1=全部
}

// Filter 过滤输入参数
func (in *RoleOptionInp) Filter(ctx context.Context) (err error) {
	return nil
}

// BatchDeleteRoleInp 批量删除角色参数
type BatchDeleteRoleInp struct {
	Ids []int64 `json:"ids" v:"required|min-length:1#角色ID列表不能为空|至少选择一个角色"`
}

// Filter 过滤输入参数
func (in *BatchDeleteRoleInp) Filter(ctx context.Context) (err error) {
	// 去重ID并过滤无效值
	if len(in.Ids) > 0 {
		idMap := make(map[int64]bool)
		uniqueIds := make([]int64, 0)
		for _, id := range in.Ids {
			if id > 0 && !idMap[id] {
				idMap[id] = true
				uniqueIds = append(uniqueIds, id)
			}
		}
		in.Ids = uniqueIds
	}

	return nil
}

// CopyRoleInp 复制角色参数
type CopyRoleInp struct {
	Id   int64  `json:"id" v:"required|min:1#角色ID不能为空|角色ID必须大于0"`
	Name string `json:"name" v:"required|length:1,50#新角色名称不能为空|角色名称长度不能超过50个字符"`
	Code string `json:"code" v:"required|length:1,50#新角色编码不能为空|角色编码长度不能超过50个字符"`
}

// Filter 过滤输入参数
func (in *CopyRoleInp) Filter(ctx context.Context) (err error) {
	// 去除前后空格
	in.Name = strings.TrimSpace(in.Name)
	in.Code = strings.TrimSpace(in.Code)

	// 角色编码转小写
	in.Code = strings.ToLower(in.Code)

	return nil
}

// RolePermissionInp 角色权限详情查询参数
type RolePermissionInp struct {
	RoleId int64 `json:"roleId" v:"required|min:1#角色ID不能为空|角色ID必须大于0"`
}

// Filter 过滤输入参数
func (in *RolePermissionInp) Filter(ctx context.Context) (err error) {
	return nil
}
