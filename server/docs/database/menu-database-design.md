# 菜单数据库设计文档

## 概述

本文档详细描述了基于开发规范指南构建的菜单系统数据库设计，包括表结构、索引设计、数据模型等内容。菜单系统支持无限层级的树形结构，为权限管理系统提供细粒度的功能权限控制。

## 数据库表设计

### 1. 菜单表（menus）

#### 表结构

```sql
CREATE TABLE IF NOT EXISTS `menus` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `parent_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '父菜单ID，0表示顶级菜单',
  `title` varchar(100) NOT NULL COMMENT '菜单标题',
  `name` varchar(100) NOT NULL COMMENT '菜单名称，用于路由name',
  `path` varchar(200) NOT NULL COMMENT '菜单路径',
  `component` varchar(200) DEFAULT NULL COMMENT '组件路径',
  `icon` varchar(100) DEFAULT NULL COMMENT '菜单图标',
  `type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '菜单类型：1=目录 2=菜单 3=按钮',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序号，数字越小越靠前',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1=启用 0=禁用',
  `visible` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否显示：1=显示 0=隐藏',
  `permission` varchar(200) DEFAULT NULL COMMENT '权限标识',
  `redirect` varchar(200) DEFAULT NULL COMMENT '重定向地址',
  `always_show` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否总是显示：1=是 0=否',
  `breadcrumb` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否显示面包屑：1=显示 0=隐藏',
  `active_menu` varchar(200) DEFAULT NULL COMMENT '高亮菜单路径',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注说明',
  `created_by` bigint(20) unsigned DEFAULT NULL COMMENT '创建人ID',
  `updated_by` bigint(20) unsigned DEFAULT NULL COMMENT '修改人ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_name` (`name`),
  KEY `idx_path` (`path`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_type` (`type`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单表';
```

#### 字段说明

| 字段名      | 类型                | 说明           | 备注                         |
| ----------- | ------------------- | -------------- | ---------------------------- |
| id          | bigint(20) unsigned | 主键 ID        | 自增主键                     |
| parent_id   | bigint(20) unsigned | 父菜单 ID      | 0 表示顶级菜单，支持无限层级 |
| title       | varchar(100)        | 菜单标题       | 显示在页面上的菜单名称       |
| name        | varchar(100)        | 菜单名称       | 用于 Vue 路由的 name 属性    |
| path        | varchar(200)        | 菜单路径       | 路由 path，如：/system/menu  |
| component   | varchar(200)        | 组件路径       | Vue 组件文件路径             |
| icon        | varchar(100)        | 菜单图标       | 图标 class 名称              |
| type        | tinyint(4)          | 菜单类型       | 1=目录，2=菜单，3=按钮       |
| sort        | int(11)             | 排序号         | 数字越小越靠前               |
| status      | tinyint(4)          | 状态           | 1=启用，0=禁用               |
| visible     | tinyint(4)          | 是否显示       | 1=显示，0=隐藏               |
| permission  | varchar(200)        | 权限标识       | 权限控制字符串               |
| redirect    | varchar(200)        | 重定向地址     | 访问该路由时重定向的地址     |
| always_show | tinyint(4)          | 是否总是显示   | 1=总是显示，0=根据子菜单决定 |
| breadcrumb  | tinyint(4)          | 是否显示面包屑 | 1=显示，0=隐藏               |
| active_menu | varchar(200)        | 高亮菜单路径   | 指定哪个菜单高亮             |
| remark      | varchar(500)        | 备注说明       | 菜单说明信息                 |
| created_by  | bigint(20) unsigned | 创建人 ID      | 创建该菜单的用户 ID          |
| updated_by  | bigint(20) unsigned | 修改人 ID      | 最后修改该菜单的用户 ID      |
| created_at  | datetime            | 创建时间       | 记录创建时间                 |
| updated_at  | datetime            | 更新时间       | 记录最后更新时间             |

## 菜单类型说明

### 1. 目录（type=1）

- 用作菜单分组，通常不对应具体页面
- 可以包含子菜单
- component 通常为'Layout'或其他布局组件
- 例如：系统管理、系统工具等

### 2. 菜单（type=2）

- 对应具体的页面
- 有实际的 component 路径
- 可以有操作按钮作为子项
- 例如：菜单管理、系统日志等

### 3. 按钮（type=3）

- 表示页面上的操作按钮
- 通常 visible=0（隐藏，不在菜单中显示）
- 主要用于权限控制
- 例如：新增、修改、删除等操作

## 菜单常量定义

系统中定义了以下菜单相关常量：

```go
// 菜单类型常量
const (
    MenuTypeDir    = 1 // 目录
    MenuTypeMenu   = 2 // 菜单
    MenuTypeButton = 3 // 按钮
)

// 菜单状态常量
const (
    MenuStatusDisabled = 0 // 禁用
    MenuStatusEnabled  = 1 // 启用
)

// 菜单可见性常量
const (
    MenuHidden  = 0 // 隐藏
    MenuVisible = 1 // 显示
)

// 菜单总是显示常量
const (
    MenuAlwaysShowNo  = 0 // 否
    MenuAlwaysShowYes = 1 // 是
)

// 面包屑显示常量
const (
    MenuBreadcrumbHidden  = 0 // 隐藏
    MenuBreadcrumbVisible = 1 // 显示
)
```

## 层级结构设计

菜单系统支持无限层级结构：

```
├── 系统管理 (id=1, parent_id=0, type=1)
│   ├── 菜单管理 (id=2, parent_id=1, type=2)
│   │   ├── 查看菜单 (id=3, parent_id=2, type=3)
│   │   ├── 新增菜单 (id=4, parent_id=2, type=3)
│   │   ├── 修改菜单 (id=5, parent_id=2, type=3)
│   │   └── 删除菜单 (id=6, parent_id=2, type=3)
├── 系统工具 (id=40, parent_id=0, type=1)
│   ├── 系统日志 (id=41, parent_id=40, type=2)
│   │   ├── 查看日志 (id=42, parent_id=41, type=3)
│   │   ├── 删除日志 (id=43, parent_id=41, type=3)
│   │   └── 清空日志 (id=44, parent_id=41, type=3)
│   └── 配置管理 (id=50, parent_id=40, type=2)
│       ├── 查看配置 (id=51, parent_id=50, type=3)
│       └── 修改配置 (id=52, parent_id=50, type=3)
```

## 权限设计

权限标识采用模块:功能:操作的格式：

```
system:menu:list     - 系统管理 -> 菜单管理 -> 列表查看
system:menu:add      - 系统管理 -> 菜单管理 -> 新增
system:menu:edit     - 系统管理 -> 菜单管理 -> 修改
system:menu:remove   - 系统管理 -> 菜单管理 -> 删除
tool:log:list        - 系统工具 -> 系统日志 -> 列表查看
tool:log:remove      - 系统工具 -> 系统日志 -> 删除
tool:config:edit     - 系统工具 -> 配置管理 -> 修改
```

## 数据模型设计

### 1. 菜单实体（Menu）

```go
type Menu struct {
    Id         int64       `json:"id"          description:"主键ID"`
    ParentId   int64       `json:"parentId"    description:"父菜单ID，0表示顶级菜单"`
    Title      string      `json:"title"       description:"菜单标题"`
    Name       string      `json:"name"        description:"菜单名称，用于路由name"`
    Path       string      `json:"path"        description:"菜单路径"`
    Component  string      `json:"component"   description:"组件路径"`
    Icon       string      `json:"icon"        description:"菜单图标"`
    Type       int         `json:"type"        description:"菜单类型：1=目录 2=菜单 3=按钮"`
    Sort       int         `json:"sort"        description:"排序号，数字越小越靠前"`
    Status     int         `json:"status"      description:"状态：1=启用 0=禁用"`
    Visible    int         `json:"visible"     description:"是否显示：1=显示 0=隐藏"`
    Permission string      `json:"permission"  description:"权限标识"`
    Redirect   string      `json:"redirect"    description:"重定向地址"`
    AlwaysShow int         `json:"alwaysShow"  description:"是否总是显示：1=是 0=否"`
    Breadcrumb int         `json:"breadcrumb"  description:"是否显示面包屑：1=显示 0=隐藏"`
    ActiveMenu string      `json:"activeMenu"  description:"高亮菜单路径"`
    Remark     string      `json:"remark"      description:"备注说明"`
    CreatedBy  int64       `json:"createdBy"   description:"创建人ID"`
    UpdatedBy  int64       `json:"updatedBy"   description:"修改人ID"`
    CreatedAt  *gtime.Time `json:"createdAt"   description:"创建时间"`
    UpdatedAt  *gtime.Time `json:"updatedAt"   description:"更新时间"`
}
```

### 2. 菜单树结构（MenuTree）

```go
type MenuTree struct {
    Menu
    Children []*MenuTree `json:"children,omitempty" description:"子菜单"`
}
```

## 实体方法

菜单实体提供了以下便利方法：

```go
// 判断菜单类型
func (m *Menu) IsDir() bool       // 是否为目录
func (m *Menu) IsMenu() bool      // 是否为菜单
func (m *Menu) IsButton() bool    // 是否为按钮

// 判断菜单状态
func (m *Menu) IsEnabled() bool   // 是否启用
func (m *Menu) IsVisible() bool   // 是否可见
func (m *Menu) HasPermission() bool // 是否有权限标识

// 获取名称
func (m *Menu) GetTypeName() string   // 获取菜单类型名称
func (m *Menu) GetStatusName() string // 获取状态名称
```

## 索引设计

### 1. 主要索引

- **主键索引**: `PRIMARY KEY (id)` - 保证数据唯一性
- **普通索引**:
  - `idx_parent_id (parent_id)` - 优化层级查询
  - `idx_name (name)` - 优化按名称查询
  - `idx_path (path)` - 优化按路径查询
  - `idx_status (status)` - 优化按状态过滤
  - `idx_sort (sort)` - 优化排序查询
  - `idx_type (type)` - 优化按类型查询
  - `idx_created_at (created_at)` - 优化时间范围查询

### 2. 索引使用场景

```sql
-- 查询顶级菜单（使用idx_parent_id）
SELECT * FROM menus WHERE parent_id = 0 AND status = 1;

-- 按名称搜索（使用idx_name）
SELECT * FROM menus WHERE name LIKE '%menu%';

-- 查询某个路径的菜单（使用idx_path）
SELECT * FROM menus WHERE path = '/system/menu';

-- 按类型查询（使用idx_type）
SELECT * FROM menus WHERE type = 2 AND status = 1;
```

## 预置菜单数据

系统预置了以下菜单结构：

### 1. 系统管理模块

- **系统管理** (目录)
  - **菜单管理** (菜单)
    - 查看菜单 (按钮)
    - 新增菜单 (按钮)
    - 修改菜单 (按钮)
    - 删除菜单 (按钮)

### 2. 系统工具模块

- **系统工具** (目录)
  - **系统日志** (菜单)
    - 查看日志 (按钮)
    - 删除日志 (按钮)
    - 清空日志 (按钮)
  - **配置管理** (菜单)
    - 查看配置 (按钮)
    - 修改配置 (按钮)

## 菜单权限验证

### 1. 权限检查流程

```go
// 检查用户是否有菜单访问权限
func CheckMenuPermission(userId int64, permission string) bool {
    // 1. 获取用户角色
    userRoles := GetUserRoles(userId)
    
    // 2. 获取角色关联的菜单权限
    for _, role := range userRoles {
        menus := GetRoleMenus(role.Id)
        for _, menu := range menus {
            if menu.Permission == permission {
                return true
            }
        }
    }
    
    return false
}
```

### 2. 前端菜单生成

```go
// 根据用户权限生成菜单树
func GenerateUserMenuTree(userId int64) []*MenuTree {
    // 1. 获取用户所有可访问的菜单
    userMenus := GetUserAccessibleMenus(userId)
    
    // 2. 构建菜单树结构
    return BuildMenuTree(userMenus, 0)
}
```

## 使用说明

1. **执行SQL文件**: 运行 `internal/sql/menus.sql` 创建表结构和初始数据
2. **数据模型**: 使用 `internal/model/entity/menus.go` 中的实体模型
3. **输入模型**: 使用 `internal/model/input/sysin/menu.go` 中的输入参数模型
4. **输出模型**: 使用 `internal/model/output/sysout/menu.go` 中的响应模型

## 扩展功能

### 1. 菜单缓存

支持菜单数据缓存，提高查询性能：

```go
// 缓存菜单树
cacheKey := fmt.Sprintf("menu_tree_%d", userId)
gcache.Set(ctx, cacheKey, menuTree, 1*time.Hour)

// 缓存用户权限
permissionKey := fmt.Sprintf("user_permissions_%d", userId)
gcache.Set(ctx, permissionKey, permissions, 30*time.Minute)
```

### 2. 菜单国际化

支持多语言菜单标题：

```sql
-- 菜单国际化表（可选扩展）
CREATE TABLE `menu_i18n` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `menu_id` bigint(20) unsigned NOT NULL COMMENT '菜单ID',
  `language` varchar(10) NOT NULL COMMENT '语言代码',
  `title` varchar(100) NOT NULL COMMENT '菜单标题',
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_menu_lang` (`menu_id`, `language`),
  KEY `idx_menu_id` (`menu_id`)
);
```

### 3. 菜单访问统计

记录菜单访问频率：

```sql
-- 菜单访问统计表（可选扩展）
CREATE TABLE `menu_access_stats` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `menu_id` bigint(20) unsigned NOT NULL COMMENT '菜单ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `access_count` int(11) DEFAULT '1' COMMENT '访问次数',
  `last_access_at` datetime NOT NULL COMMENT '最后访问时间',
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_menu_user` (`menu_id`, `user_id`),
  KEY `idx_menu_id` (`menu_id`),
  KEY `idx_user_id` (`user_id`)
);
```

## 特点优势

1. **结构清晰**: 支持无限层级的树形结构
2. **功能完整**: 涵盖目录、菜单、按钮三种类型
3. **权限精细**: 支持细粒度的权限控制
4. **性能优秀**: 合理的索引设计保证查询效率
5. **扩展性强**: 预留了多种扩展可能性
6. **易于维护**: 完善的文档和规范的代码结构
