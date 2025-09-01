# 业务模块文档

## 项目概述

本项目(client-app)是一个基于GoFrame v2.9.0框架开发的完整RBAC (Role-Based Access Control) 权限管理系统，提供用户认证、角色管理、菜单权限控制等核心功能。项目采用分层架构设计，通过接口抽象实现高内聚低耦合。

## 核心业务模块

### 1. 用户认证模块 (User Authentication)

#### 1.1 模块概述
用户认证模块负责处理用户的登录、退出、密码管理等核心认证功能，采用JWT Token + RefreshToken机制实现无状态认证。

#### 1.2 主要功能

- **用户登录**: 支持用户名密码登录，集成图形验证码防护，支持JWT Token生成
- **用户退出**: 支持Token黑名单机制，确保安全退出
- **密码管理**: 支持用户修改密码，使用MD5+盐值加密存储
- **Token管理**: 支持AccessToken + RefreshToken双令牌机制
- **验证码服务**: 提供图形验证码生成和验证功能
- **用户信息**: 提供用户资料查询和管理
- **登录追踪**: 记录登录IP、时间、次数等信息

#### 1.3 技术实现

**服务接口定义** (internal/service/user.go - 65行):
```go
type IUser interface {
    // Login 用户登录
    Login(ctx context.Context, in *sysin.UserLoginInp) (res *sysout.LoginTokenModel, err error)
    
    // Logout 用户退出
    Logout(ctx context.Context, in *sysin.UserLogoutInp) error
    
    // GetProfile 获取用户资料
    GetProfile(ctx context.Context, userId int64) (res *sysout.UserModel, err error)
    
    // RefreshToken 刷新访问令牌
    RefreshToken(ctx context.Context, refreshToken string) (res *TokenInfo, err error)
    
    // ChangePassword 修改密码
    ChangePassword(ctx context.Context, userId int64, oldPassword, newPassword string) error
    
    // GenerateCaptcha 生成验证码
    GenerateCaptcha(ctx context.Context) (captchaId, captchaImage string, err error)
    
    // VerifyCaptcha 验证验证码
    VerifyCaptcha(ctx context.Context, captchaId, captcha string) error
    
    // GetUserByUsername 根据用户名获取用户
    GetUserByUsername(ctx context.Context, username string) (user *sysout.UserModel, err error)
    
    // ValidateUser 验证用户密码
    ValidateUser(ctx context.Context, username, password string) (user *sysout.UserModel, err error)
}
```

**业务逻辑实现** (internal/logic/api/user.go - 441行):
实现了完整的用户认证流程，包括:
- 密码验证和用户状态检查
- JWT Token生成和刷新令牌管理
- 用户权限获取和登录信息更新
- 验证码生成验证等安全功能

#### 1.4 安全特性

- **密码加密**: 使用MD5+盐值方式加密存储密码
- **JWT安全**: 支持Token过期、签名验证，24小时有效期
- **验证码防护**: 防止暴力破解登录
- **双令牌机制**: AccessToken + RefreshToken延长会话安全性
- **黑名单机制**: 退出时将Token加入缓存黑名单
- **状态检查**: 支持用户正常、锁定、禁用状态管理

#### 1.5 用户实体设计

**用户实体** (internal/model/entity/users.go - 337行):
```go
type User struct {
    Id                   int64       `json:"id"                    description:"主键ID"`
    Username             string      `json:"username"              description:"用户名"`
    Email                string      `json:"email"                 description:"邮箱地址"`
    Phone                string      `json:"phone"                 description:"手机号码"`
    Password             string      `json:"-"                     description:"密码（不返回给前端）"`
    Salt                 string      `json:"-"                     description:"密码盐值（不返回给前端）"`
    RealName             string      `json:"realName"              description:"真实姓名"`
    Nickname             string      `json:"nickname"              description:"昵称"`
    Avatar               string      `json:"avatar"                description:"头像URL"`
    Gender               int         `json:"gender"                description:"性别"`
    Birthday             *gtime.Time `json:"birthday"              description:"生日"`
    DeptId               int64       `json:"deptId"                description:"部门ID"`
    Position             string      `json:"position"              description:"职位"`
    Status               int         `json:"status"                description:"状态"`
    LoginIp              string      `json:"loginIp"               description:"最后登录IP"`
    LoginAt              *gtime.Time `json:"loginAt"               description:"最后登录时间"`
    LoginCount           int         `json:"loginCount"            description:"登录次数"`
    PasswordResetToken   string      `json:"-"                     description:"密码重置令牌"`
    PasswordResetExpires *gtime.Time `json:"-"                     description:"密码重置过期时间"`
    EmailVerifiedAt      *gtime.Time `json:"emailVerifiedAt"       description:"邮箱验证时间"`
    PhoneVerifiedAt      *gtime.Time `json:"phoneVerifiedAt"       description:"手机验证时间"`
    TwoFactorEnabled     int         `json:"twoFactorEnabled"      description:"是否启用双因子认证"`
    TwoFactorSecret      string      `json:"-"                     description:"双因子认证密钥"`
    Remark               string      `json:"remark"                description:"备注说明"`
    CreatedBy            int64       `json:"createdBy"             description:"创建人ID"`
    UpdatedBy            int64       `json:"updatedBy"             description:"修改人ID"`
    CreatedAt            *gtime.Time `json:"createdAt"             description:"创建时间"`
    UpdatedAt            *gtime.Time `json:"updatedAt"             description:"更新时间"`
    DeletedAt            *gtime.Time `json:"deletedAt,omitempty"   description:"删除时间"`
}
```

支持的用户状态常量:
- `UserStatusNormal = 1`   // 正常
- `UserStatusLocked = 2`   // 锁定
- `UserStatusDisabled = 3` // 禁用

### 2. 角色管理模块 (Role Management)

#### 2.1 模块概述
角色管理模块提供完整的角色CRUD操作，支持角色权限分配、数据权限控制等功能，是RBAC系统的核心组件。

#### 2.2 主要功能

- **角色CRUD**: 支持角色的创建、查询、更新、删除、批量删除操作
- **权限分配**: 支持为角色分配菜单权限，更新角色菜单关联
- **数据权限**: 支持5种数据访问范围配置
- **角色状态**: 支持角色启用/禁用状态管理
- **角色排序**: 支持角色显示顺序配置
- **角色查询**: 支持按名称、编码、状态等条件查询
- **角色统计**: 提供角色统计信息和选项数据
- **用户角色关联**: 支持用户角色分配、移除、主要角色设定

#### 2.3 服务接口定义

**角色服务接口** (internal/service/role.go - 62行):
```go
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
```

#### 2.4 数据权限范围

支持5种数据权限范围：
- **全部数据 (DataScopeAll = 1)**: 可访问系统中所有数据
- **部门数据 (DataScopeDept = 2)**: 只能访问所属部门的数据
- **部门及下级 (DataScopeDeptAndSub = 3)**: 可访问所属部门及下级部门数据
- **仅本人数据 (DataScopeSelf = 4)**: 只能访问自己创建的数据
- **自定义权限 (DataScopeCustom = 5)**: 支持自定义数据访问规则

#### 2.5 角色实体设计

**角色实体** (internal/model/entity/roles.go - 180行):
```go
type Role struct {
    Id          int64       `json:"id"          description:"主键ID"`
    Name        string      `json:"name"        description:"角色名称"`
    Code        string      `json:"code"        description:"角色编码"`
    Description string      `json:"description" description:"角色描述"`
    Status      int         `json:"status"      description:"状态：1=启用 0=禁用"`
    Sort        int         `json:"sort"        description:"排序号，数字越小越靠前"`
    DataScope   int         `json:"dataScope"   description:"数据权限范围"`
    Remark      string      `json:"remark"      description:"备注说明"`
    CreatedBy   int64       `json:"createdBy"   description:"创建人ID"`
    UpdatedBy   int64       `json:"updatedBy"   description:"修改人ID"`
    CreatedAt   *gtime.Time `json:"createdAt"   description:"创建时间"`
    UpdatedAt   *gtime.Time `json:"updatedAt"   description:"更新时间"`
}
```

内置角色编码常量：
- `RoleCodeSuperAdmin = "super_admin"`      // 超级管理员
- `RoleCodeSystemAdmin = "system_admin"`    // 系统管理员
- `RoleCodeFinanceAdmin = "finance_admin"`  // 财务管理员
- `RoleCodeOperator = "operator"`           // 运营人员
- `RoleCodeCustomerService = "customer_service"` // 客服人员
- `RoleCodeAuditor = "auditor"`             // 审计人员

#### 2.6 业务逻辑实现

**角色逻辑实现** (internal/logic/api/role.go - 1086行):
包含完整的角色管理业务逻辑，支持：
- 复杂的角色查询和过滤
- 角色权限分配和验证
- 用户角色关联管理
- 数据权限控制逻辑
- 角色状态和统计功能

### 3. 菜单权限模块 (Menu Permission)

#### 3.1 模块概述
菜单权限模块管理系统的菜单结构和权限标识，支持层级菜单管理、权限细粒度控制等功能。

#### 3.2 主要功能

- **菜单管理**: 支持层级菜单的创建、编辑、删除
- **权限标识**: 为每个菜单配置权限标识符
- **菜单类型**: 支持目录、菜单、按钮三种类型
- **显示控制**: 支持菜单显示/隐藏状态管理
- **路由配置**: 支持前端路由路径配置
- **图标管理**: 支持菜单图标配置
- **面包屑控制**: 支持面包屑显示配置
- **层级结构**: 支持无限层级菜单树结构

#### 3.3 菜单实体设计

**菜单实体** (internal/model/entity/menus.go - 124行):
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

菜单类型常量：
- `MenuTypeDir = 1`    // 目录
- `MenuTypeMenu = 2`   // 菜单  
- `MenuTypeButton = 3` // 按钮

#### 3.4 权限验证机制

通过中间件实现权限验证：

```go
// 权限验证中间件
func (s *sMiddleware) ApiAuth(r *ghttp.Request) {
    // 1. 获取Token
    token := s.getTokenFromRequest(r)
    
    // 2. 验证Token
    payload, err := simple.ParseJWTToken(token, secretKey)
    if err != nil {
        response.JsonExit(r, consts.ErrorTokenInvalid)
    }
    
    // 3. 验证用户状态
    user := s.getCurrentUser(r.Context(), payload.UserId)
    if !user.CanLogin() {
        response.JsonExit(r, consts.ErrorUserDisabled)
    }
    
    // 4. 验证权限
    if !s.hasPermission(r.Context(), payload.UserId, r.URL.Path) {
        response.JsonExit(r, consts.ErrorPermissionDenied)
    }
    
    // 5. 设置用户上下文
    r.SetCtxVar(consts.CtxUserId, payload.UserId)
    r.Middleware.Next()
}
```

### 4. 中间件服务模块

#### 4.1 中间件接口

**中间件服务** (internal/service/middleware.go):
提供HTTP请求处理中间件，包括：
- 用户认证中间件(ApiAuth)
- CORS跨域处理
- 请求日志记录
- 响应处理中间件
- 上下文处理中间件

#### 4.2 中间件实现

实现位于 `internal/logic/middleware/` 目录，提供完整的中间件功能支持。

## 数据库设计

### 核心表结构

#### 1. 用户表 (users)
```sql
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱地址',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号码',
  `password` varchar(64) NOT NULL COMMENT '密码',
  `salt` varchar(32) NOT NULL COMMENT '密码盐值',
  `real_name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
  `gender` tinyint DEFAULT '0' COMMENT '性别：0=未知 1=男 2=女',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `position` varchar(50) DEFAULT NULL COMMENT '职位',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1=正常 2=锁定 3=禁用',
  `login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
  `login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `login_count` int DEFAULT '0' COMMENT '登录次数',
  `password_reset_token` varchar(64) DEFAULT NULL COMMENT '密码重置令牌',
  `password_reset_expires` datetime DEFAULT NULL COMMENT '密码重置过期时间',
  `email_verified_at` datetime DEFAULT NULL COMMENT '邮箱验证时间',
  `phone_verified_at` datetime DEFAULT NULL COMMENT '手机验证时间',
  `two_factor_enabled` tinyint DEFAULT '0' COMMENT '是否启用双因子认证',
  `two_factor_secret` varchar(64) DEFAULT NULL COMMENT '双因子认证密钥',
  `remark` text COMMENT '备注说明',
  `created_by` bigint DEFAULT NULL COMMENT '创建人ID',
  `updated_by` bigint DEFAULT NULL COMMENT '修改人ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

#### 2. 角色表 (roles)
```sql
CREATE TABLE `roles` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` varchar(255) DEFAULT NULL COMMENT '角色描述',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1=启用 0=禁用',
  `sort` int DEFAULT '0' COMMENT '排序号，数字越小越靠前',
  `data_scope` tinyint DEFAULT '1' COMMENT '数据权限范围：1=全部 2=部门 3=部门及下级 4=仅本人 5=自定义',
  `remark` text COMMENT '备注说明',
  `created_by` bigint DEFAULT NULL COMMENT '创建人ID',
  `updated_by` bigint DEFAULT NULL COMMENT '修改人ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';
```

#### 3. 菜单表 (menus)
```sql
CREATE TABLE `menus` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `parent_id` bigint DEFAULT '0' COMMENT '父菜单ID，0表示顶级菜单',
  `title` varchar(50) NOT NULL COMMENT '菜单标题',
  `name` varchar(50) DEFAULT NULL COMMENT '菜单名称，用于路由name',
  `path` varchar(255) DEFAULT NULL COMMENT '菜单路径',
  `component` varchar(255) DEFAULT NULL COMMENT '组件路径',
  `icon` varchar(50) DEFAULT NULL COMMENT '菜单图标',
  `type` tinyint NOT NULL DEFAULT '1' COMMENT '菜单类型：1=目录 2=菜单 3=按钮',
  `sort` int DEFAULT '0' COMMENT '排序号，数字越小越靠前',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1=启用 0=禁用',
  `visible` tinyint DEFAULT '1' COMMENT '是否显示：1=显示 0=隐藏',
  `permission` varchar(100) DEFAULT NULL COMMENT '权限标识',
  `redirect` varchar(255) DEFAULT NULL COMMENT '重定向地址',
  `always_show` tinyint DEFAULT '0' COMMENT '是否总是显示：1=是 0=否',
  `breadcrumb` tinyint DEFAULT '1' COMMENT '是否显示面包屑：1=显示 0=隐藏',
  `active_menu` varchar(255) DEFAULT NULL COMMENT '高亮菜单路径',
  `remark` text COMMENT '备注说明',
  `created_by` bigint DEFAULT NULL COMMENT '创建人ID',
  `updated_by` bigint DEFAULT NULL COMMENT '修改人ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单表';
```

#### 4. 用户角色关联表 (user_roles)
```sql
CREATE TABLE `user_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `is_primary` tinyint DEFAULT '0' COMMENT '是否主要角色：1=是 0=否',
  `assigned_by` bigint DEFAULT NULL COMMENT '分配人ID',
  `expires_at` datetime DEFAULT NULL COMMENT '过期时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`,`role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';
```

#### 5. 角色菜单关联表 (role_menus)
```sql
CREATE TABLE `role_menus` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `menu_id` bigint NOT NULL COMMENT '菜单ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_menu` (`role_id`,`menu_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色菜单关联表';
```

## 技术特性

### 1. 架构特性
- **分层架构**: Router -> Controller -> Service -> Logic -> Model/ORM
- **依赖注入**: 通过服务注册机制实现松耦合
- **接口抽象**: 业务逻辑通过接口定义，便于测试和扩展

### 2. 安全特性  
- **JWT认证**: 无状态Token认证，支持过期和刷新
- **密码安全**: MD5+盐值加密，支持密码强度验证
- **权限控制**: 基于RBAC的细粒度权限控制
- **数据权限**: 支持5种数据访问范围控制
- **双因子认证**: 支持2FA增强安全性

### 3. 性能特性
- **缓存机制**: 使用GoFrame内置缓存提升性能
- **连接池**: 数据库连接池管理
- **索引优化**: 数据库表结构包含必要索引

### 4. 扩展特性
- **多数据库支持**: 通过GoFrame ORM支持多种数据库
- **链路追踪**: 集成Jaeger进行分布式链路追踪
- **日志系统**: 完整的日志记录和管理
- **配置管理**: 灵活的配置文件管理

## 项目统计信息

### 代码行数统计
- **用户逻辑实现**: 441行 (internal/logic/api/user.go)
- **角色逻辑实现**: 1086行 (internal/logic/api/role.go)
- **用户实体定义**: 337行 (internal/model/entity/users.go)
- **角色实体定义**: 180行 (internal/model/entity/roles.go)
- **菜单实体定义**: 124行 (internal/model/entity/menus.go)
- **用户服务接口**: 65行 (internal/service/user.go)
- **角色服务接口**: 62行 (internal/service/role.go)

### 功能模块统计
- **核心业务模块**: 3个 (用户认证、角色管理、菜单权限)
- **服务接口**: 5个 (User、Role、Middleware、Hook、View)
- **数据实体**: 3个主要实体 (User、Role、Menu)
- **数据表**: 5个核心表 + 关联表
- **权限类型**: 3种菜单类型，5种数据权限范围
- **用户状态**: 3种状态 (正常、锁定、禁用)

这个RBAC权限管理系统提供了完整的用户认证、角色管理、菜单权限控制功能，通过合理的架构设计和丰富的安全特性，能够满足各种企业级应用的权限管理需求。 