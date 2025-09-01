package sysin

import (
	"context"
	"strings"
)

// MenuListInp 菜单列表查询参数
type MenuListInp struct {
	Title     string `json:"title" v:""`     // 菜单标题（模糊查询）
	Name      string `json:"name" v:""`      // 菜单名称（模糊查询）
	Status    int    `json:"status" v:""`    // 状态：1=启用 0=禁用，-1=全部
	Type      int    `json:"type" v:""`      // 菜单类型：1=目录 2=菜单 3=按钮，0=全部
	ParentId  int64  `json:"parentId" v:""`  // 父菜单ID
	Page      int    `json:"page" v:"min:1"` // 页码
	PageSize  int    `json:"pageSize" v:""`  // 每页数量
	OrderBy   string `json:"orderBy" v:""`   // 排序字段
	OrderType string `json:"orderType" v:""` // 排序方式：asc/desc
}

// Filter 过滤输入参数
func (in *MenuListInp) Filter(ctx context.Context) (err error) {
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

	return nil
}

// CreateMenuInp 创建菜单参数
type CreateMenuInp struct {
	ParentId   int64  `json:"parentId" v:"min:0#父菜单ID不能小于0"`
	Title      string `json:"title" v:"required|length:1,100#菜单标题不能为空|菜单标题长度不能超过100个字符"`
	Name       string `json:"name" v:"required|length:1,100#菜单名称不能为空|菜单名称长度不能超过100个字符"`
	Path       string `json:"path" v:"required|length:1,200#菜单路径不能为空|菜单路径长度不能超过200个字符"`
	Component  string `json:"component" v:"length:0,200#组件路径长度不能超过200个字符"`
	Icon       string `json:"icon" v:"length:0,100#菜单图标长度不能超过100个字符"`
	Type       int    `json:"type" v:"required|in:1,2,3#菜单类型不能为空|菜单类型必须是1(目录)、2(菜单)或3(按钮)"`
	Sort       int    `json:"sort" v:"min:0#排序号不能小于0"`
	Status     int    `json:"status" v:"in:0,1#状态必须是0(禁用)或1(启用)"`
	Visible    int    `json:"visible" v:"in:0,1#是否显示必须是0(隐藏)或1(显示)"`
	Permission string `json:"permission" v:"length:0,200#权限标识长度不能超过200个字符"`
	Redirect   string `json:"redirect" v:"length:0,200#重定向地址长度不能超过200个字符"`
	AlwaysShow int    `json:"alwaysShow" v:"in:0,1#是否总是显示必须是0(否)或1(是)"`
	Breadcrumb int    `json:"breadcrumb" v:"in:0,1#是否显示面包屑必须是0(隐藏)或1(显示)"`
	ActiveMenu string `json:"activeMenu" v:"length:0,200#高亮菜单路径长度不能超过200个字符"`
	Remark     string `json:"remark" v:"length:0,500#备注说明长度不能超过500个字符"`
}

// Filter 过滤输入参数
func (in *CreateMenuInp) Filter(ctx context.Context) (err error) {
	// 去除前后空格
	in.Title = strings.TrimSpace(in.Title)
	in.Name = strings.TrimSpace(in.Name)
	in.Path = strings.TrimSpace(in.Path)
	in.Component = strings.TrimSpace(in.Component)
	in.Icon = strings.TrimSpace(in.Icon)
	in.Permission = strings.TrimSpace(in.Permission)
	in.Redirect = strings.TrimSpace(in.Redirect)
	in.ActiveMenu = strings.TrimSpace(in.ActiveMenu)
	in.Remark = strings.TrimSpace(in.Remark)

	// 设置默认值
	if in.Status == 0 {
		in.Status = 1 // 默认启用
	}
	if in.Visible == 0 {
		in.Visible = 1 // 默认显示
	}
	if in.Breadcrumb == 0 {
		in.Breadcrumb = 1 // 默认显示面包屑
	}

	// 如果是按钮类型，设置为隐藏
	if in.Type == 3 {
		in.Visible = 0
	}

	// 路径规范化
	if in.Path != "" && !strings.HasPrefix(in.Path, "/") && in.Type != 3 {
		in.Path = "/" + in.Path
	}

	return nil
}

// UpdateMenuInp 更新菜单参数
type UpdateMenuInp struct {
	Id         int64  `json:"id" v:"required|min:1#菜单ID不能为空|菜单ID必须大于0"`
	ParentId   int64  `json:"parentId" v:"min:0#父菜单ID不能小于0"`
	Title      string `json:"title" v:"required|length:1,100#菜单标题不能为空|菜单标题长度不能超过100个字符"`
	Name       string `json:"name" v:"required|length:1,100#菜单名称不能为空|菜单名称长度不能超过100个字符"`
	Path       string `json:"path" v:"required|length:1,200#菜单路径不能为空|菜单路径长度不能超过200个字符"`
	Component  string `json:"component" v:"length:0,200#组件路径长度不能超过200个字符"`
	Icon       string `json:"icon" v:"length:0,100#菜单图标长度不能超过100个字符"`
	Type       int    `json:"type" v:"required|in:1,2,3#菜单类型不能为空|菜单类型必须是1(目录)、2(菜单)或3(按钮)"`
	Sort       int    `json:"sort" v:"min:0#排序号不能小于0"`
	Status     int    `json:"status" v:"in:0,1#状态必须是0(禁用)或1(启用)"`
	Visible    int    `json:"visible" v:"in:0,1#是否显示必须是0(隐藏)或1(显示)"`
	Permission string `json:"permission" v:"length:0,200#权限标识长度不能超过200个字符"`
	Redirect   string `json:"redirect" v:"length:0,200#重定向地址长度不能超过200个字符"`
	AlwaysShow int    `json:"alwaysShow" v:"in:0,1#是否总是显示必须是0(否)或1(是)"`
	Breadcrumb int    `json:"breadcrumb" v:"in:0,1#是否显示面包屑必须是0(隐藏)或1(显示)"`
	ActiveMenu string `json:"activeMenu" v:"length:0,200#高亮菜单路径长度不能超过200个字符"`
	Remark     string `json:"remark" v:"length:0,500#备注说明长度不能超过500个字符"`
}

// Filter 过滤输入参数
func (in *UpdateMenuInp) Filter(ctx context.Context) (err error) {
	// 去除前后空格
	in.Title = strings.TrimSpace(in.Title)
	in.Name = strings.TrimSpace(in.Name)
	in.Path = strings.TrimSpace(in.Path)
	in.Component = strings.TrimSpace(in.Component)
	in.Icon = strings.TrimSpace(in.Icon)
	in.Permission = strings.TrimSpace(in.Permission)
	in.Redirect = strings.TrimSpace(in.Redirect)
	in.ActiveMenu = strings.TrimSpace(in.ActiveMenu)
	in.Remark = strings.TrimSpace(in.Remark)

	// 如果是按钮类型，设置为隐藏
	if in.Type == 3 {
		in.Visible = 0
	}

	// 路径规范化
	if in.Path != "" && !strings.HasPrefix(in.Path, "/") && in.Type != 3 {
		in.Path = "/" + in.Path
	}

	return nil
}

// DeleteMenuInp 删除菜单参数
type DeleteMenuInp struct {
	Id int64 `json:"id" v:"required|min:1#菜单ID不能为空|菜单ID必须大于0"`
}

// Filter 过滤输入参数
func (in *DeleteMenuInp) Filter(ctx context.Context) (err error) {
	return nil
}

// MenuDetailInp 菜单详情查询参数
type MenuDetailInp struct {
	Id int64 `json:"id" v:"required|min:1#菜单ID不能为空|菜单ID必须大于0"`
}

// Filter 过滤输入参数
func (in *MenuDetailInp) Filter(ctx context.Context) (err error) {
	return nil
}

// MenuTreeInp 菜单树查询参数
type MenuTreeInp struct {
	Status int `json:"status" v:""` // 状态过滤：1=启用 0=禁用，-1=全部
	Type   int `json:"type" v:""`   // 类型过滤：1=目录 2=菜单 3=按钮，0=全部
}

// Filter 过滤输入参数
func (in *MenuTreeInp) Filter(ctx context.Context) (err error) {
	return nil
}

// UpdateMenuStatusInp 更新菜单状态参数
type UpdateMenuStatusInp struct {
	Id     int64 `json:"id" v:"required|min:1#菜单ID不能为空|菜单ID必须大于0"`
	Status int   `json:"status" v:"required|in:0,1#状态不能为空|状态必须是0(禁用)或1(启用)"`
}

// Filter 过滤输入参数
func (in *UpdateMenuStatusInp) Filter(ctx context.Context) (err error) {
	return nil
}

// BatchDeleteMenuInp 批量删除菜单参数
type BatchDeleteMenuInp struct {
	Ids []int64 `json:"ids" v:"required|min-length:1#菜单ID列表不能为空|至少选择一个菜单"`
}

// Filter 过滤输入参数
func (in *BatchDeleteMenuInp) Filter(ctx context.Context) (err error) {
	// 去重并过滤无效ID
	var validIds []int64
	idMap := make(map[int64]bool)
	for _, id := range in.Ids {
		if id > 0 && !idMap[id] {
			validIds = append(validIds, id)
			idMap[id] = true
		}
	}
	in.Ids = validIds
	return nil
}

// MenuOptionInp 菜单选项查询参数
type MenuOptionInp struct {
	Type       int   `json:"type" v:""`       // 菜单类型过滤：1=目录 2=菜单 3=按钮，0=全部
	Status     int   `json:"status" v:""`     // 状态过滤：1=启用 0=禁用，-1=全部
	ExcludeId  int64 `json:"excludeId" v:""`  // 排除的菜单ID（用于编辑时排除自己）
	ParentOnly bool  `json:"parentOnly" v:""` // 是否只返回可作为父菜单的选项（目录和菜单）
}

// Filter 过滤输入参数
func (in *MenuOptionInp) Filter(ctx context.Context) (err error) {
	return nil
}
