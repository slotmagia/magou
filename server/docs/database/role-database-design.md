# 角色数据库设计文档

## 概述

本文档详细描述了基于开发规范指南构建的角色权限系统数据库设计，包括角色表、角色菜单关联表的结构设计、索引优化、权限控制等内容。该系统采用RBAC (Role-Based Access Control) 模型实现权限管理。

## 数据库表设计

### 1. 角色表（roles）

#### 表结构

```sql
CREATE TABLE IF NOT EXISTS `roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` varchar(200) DEFAULT NULL COMMENT '角色描述',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1=启用 0=禁用',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序号，数字越小越靠前',
  `data_scope` tinyint(4) NOT NULL DEFAULT '1' COMMENT '数据权限范围：1=全部数据 2=部门数据 3=部门及以下数据 4=仅本人数据 5=自定义权限',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注说明',
  `created_by` bigint(20) unsigned DEFAULT NULL COMMENT '创建人ID',
  `updated_by` bigint(20) unsigned DEFAULT NULL COMMENT '修改人ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_name` (`name`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';
```

#### 字段说明

| 字段名      | 类型                | 说明         | 备注                     |
| ----------- | ------------------- | ------------ | ------------------------ |
| id          | bigint(20) unsigned | 主键 ID      | 自增主键                 |
| name        | varchar(50)         | 角色名称     | 显示给用户的角色名称     |
| code        | varchar(50)         | 角色编码     | 系统内部使用的唯一标识   |
| description | varchar(200)        | 角色描述     | 角色的详细说明           |
| status      | tinyint(4)          | 状态         | 1=启用，0=禁用           |
| sort        | int(11)             | 排序号       | 数字越小越靠前           |
| data_scope  | tinyint(4)          | 数据权限范围 | 控制用户可访问的数据范围 |
| remark      | varchar(500)        | 备注说明     | 管理员备注信息           |
| created_by  | bigint(20) unsigned | 创建人 ID    | 创建该角色的用户 ID      |
| updated_by  | bigint(20) unsigned | 修改人 ID    | 最后修改该角色的用户 ID  |
| created_at  | datetime            | 创建时间     | 记录创建时间             |
| updated_at  | datetime            | 更新时间     | 记录最后更新时间         |

### 2. 角色菜单关联表（role_menus）

#### 表结构

```sql
CREATE TABLE IF NOT EXISTS `role_menus` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
  `menu_id` bigint(20) unsigned NOT NULL COMMENT '菜单ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_menu` (`role_id`, `menu_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_menu_id` (`menu_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色菜单关联表';
```

#### 字段说明

| 字段名     | 类型                | 说明     | 备注               |
| ---------- | ------------------- | -------- | ------------------ |
| id         | bigint(20) unsigned | 主键 ID  | 自增主键           |
| role_id    | bigint(20) unsigned | 角色 ID  | 关联 roles 表的 id |
| menu_id    | bigint(20) unsigned | 菜单 ID  | 关联 menus 表的 id |
| created_at | datetime            | 创建时间 | 权限分配时间       |

## 角色类型设计

### 1. 内置角色

系统预置了 6 种基础角色，涵盖权限管理系统的主要管理场景：

#### 超级管理员 (super_admin)

- **权限范围**: 系统所有权限
- **数据权限**: 全部数据
- **主要职责**: 系统最高权限管理者
- **特殊属性**: 不可删除，内置角色

#### 系统管理员 (system_admin)

- **权限范围**: 系统管理、用户管理、角色管理、菜单管理
- **数据权限**: 部门数据
- **主要职责**: 负责系统基础配置和用户管理

#### 财务管理员 (finance_admin)

- **权限范围**: 财务统计、报表导出、数据分析
- **数据权限**: 部门及以下数据
- **主要职责**: 负责财务相关功能管理

#### 运营人员 (operator)

- **权限范围**: 日常操作、基础查询、内容管理
- **数据权限**: 仅本人数据
- **主要职责**: 负责日常业务操作

#### 客服人员 (customer_service)

- **权限范围**: 用户服务、问题处理、基础查询
- **数据权限**: 仅本人数据
- **主要职责**: 负责客户服务和问题处理

#### 审计人员 (auditor)

- **权限范围**: 系统审计、只读权限、日志查看
- **数据权限**: 全部数据
- **主要职责**: 负责系统审计和监督

## 数据权限设计

### 1. 权限范围类型

| 类型值 | 名称           | 说明                           | 使用场景             |
| ------ | -------------- | ------------------------------ | -------------------- |
| 1      | 全部数据       | 可以访问系统所有数据           | 超级管理员、审计人员 |
| 2      | 部门数据       | 只能访问所在部门的数据         | 部门主管             |
| 3      | 部门及以下数据 | 可以访问所在部门及下级部门数据 | 高级主管             |
| 4      | 仅本人数据     | 只能访问自己创建或负责的数据   | 普通员工             |
| 5      | 自定义权限     | 根据具体规则自定义数据范围     | 特殊角色             |

### 2. 权限控制策略

#### 数据过滤规则

```sql
-- 全部数据权限：无需过滤
WHERE 1=1

-- 部门数据权限：
WHERE dept_id = :user_dept_id

-- 部门及以下数据权限：
WHERE dept_id IN (SELECT id FROM departments WHERE path LIKE CONCAT(:user_dept_path, '%'))

-- 仅本人数据权限：
WHERE created_by = :user_id OR assigned_to = :user_id
```

#### 菜单权限控制

- 系统根据角色分配的菜单权限动态生成用户可访问的功能
- 支持细粒度的按钮级权限控制
- 前端根据权限动态显示/隐藏功能按钮

## 角色权限分配

### 1. 预置权限分配

各角色默认拥有的菜单权限：

#### 超级管理员

- 拥有系统所有菜单权限（动态分配）

#### 系统管理员

- 系统管理模块：菜单管理（增删改查）、角色管理（全部权限）
- 系统工具模块：日志管理、配置管理

#### 财务管理员

- 系统管理模块：数据查看权限
- 系统工具模块：日志查看、报表导出

#### 运营人员

- 系统管理模块：基础查看权限
- 系统工具模块：日志查看

#### 客服人员

- 系统管理模块：用户查看权限
- 基础功能模块：问题处理相关功能

#### 审计人员

- 系统工具模块：系统日志（全部权限）、配置查看
- 数据审计模块：所有审计相关功能

## 角色状态管理

### 1. 角色状态常量

```go
const (
    RoleStatusDisabled = 0 // 禁用
    RoleStatusEnabled  = 1 // 启用
)
```

### 2. 内置角色管理

系统内置角色具有以下特性：
- 不可删除（通过IsBuiltIn()方法判断）
- 编码固定不可修改
- 具有特殊权限验证逻辑

内置角色编码列表：
```go
const (
    RoleCodeSuperAdmin     = "super_admin"     // 超级管理员
    RoleCodeSystemAdmin    = "system_admin"    // 系统管理员
    RoleCodeFinanceAdmin   = "finance_admin"   // 财务管理员
    RoleCodeOperator       = "operator"        // 运营人员
    RoleCodeCustomerService = "customer_service" // 客服人员
    RoleCodeAuditor        = "auditor"         // 审计人员
)
```

## 数据库索引设计

### 1. 主要索引

- **主键索引**: `PRIMARY KEY (id)` - 保证数据唯一性
- **唯一索引**: `UNIQUE KEY uk_code (code)` - 保证角色编码唯一
- **普通索引**: 
  - `idx_name (name)` - 优化按名称查询
  - `idx_status (status)` - 优化按状态过滤
  - `idx_sort (sort)` - 优化排序查询
  - `idx_created_at (created_at)` - 优化时间范围查询

### 2. 关联表索引

角色菜单关联表(role_menus)索引：
- **唯一索引**: `uk_role_menu (role_id, menu_id)` - 防止重复关联
- **普通索引**:
  - `idx_role_id (role_id)` - 优化按角色查询菜单
  - `idx_menu_id (menu_id)` - 优化按菜单查询角色

## 使用说明

1. **执行SQL文件**: 运行 `internal/sql/roles.sql` 创建表结构和初始数据
2. **数据模型**: 使用 `internal/model/entity/roles.go` 中的实体模型
3. **输入模型**: 使用 `internal/model/input/sysin/role.go` 中的输入参数模型
4. **输出模型**: 使用 `internal/model/output/sysout/role.go` 中的响应模型

## 扩展功能

### 1. 角色继承

支持角色继承机制，子角色自动继承父角色的权限：

```sql
-- 角色继承表（可选扩展）
CREATE TABLE `role_inheritance` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `parent_role_id` bigint(20) unsigned NOT NULL COMMENT '父角色ID',
  `child_role_id` bigint(20) unsigned NOT NULL COMMENT '子角色ID',
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_parent_child` (`parent_role_id`, `child_role_id`)
);
```

### 2. 临时权限

支持给角色分配临时权限，具有时效性：

```sql
-- 临时权限表（可选扩展）
CREATE TABLE `role_temp_permissions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
  `permission` varchar(200) NOT NULL COMMENT '权限标识',
  `expires_at` datetime NOT NULL COMMENT '过期时间',
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_role_permission` (`role_id`, `permission`),
  KEY `idx_expires_at` (`expires_at`)
);
```

## 特点优势

1. **设计规范**: 严格遵循RBAC标准模型
2. **扩展性强**: 支持自定义角色和权限扩展
3. **性能优秀**: 合理的索引设计保证查询效率
4. **安全可靠**: 多层次的权限控制机制
5. **易于维护**: 清晰的数据结构和完善的文档
