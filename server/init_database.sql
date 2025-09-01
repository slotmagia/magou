-- 创建数据库
CREATE DATABASE IF NOT EXISTS `admin` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE `admin`;

-- 用户表
CREATE TABLE IF NOT EXISTS `sys_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `salt` varchar(32) NOT NULL COMMENT '密码盐',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `real_name` varchar(50) NOT NULL COMMENT '真实姓名',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `gender` tinyint(1) DEFAULT '0' COMMENT '性别：0=未知 1=男 2=女',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `dept_id` bigint(20) DEFAULT '0' COMMENT '部门ID',
  `position` varchar(100) DEFAULT NULL COMMENT '职位',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态：1=正常 2=锁定 3=禁用',
  `login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `login_count` int(11) DEFAULT '0' COMMENT '登录次数',
  `two_factor_enabled` tinyint(1) DEFAULT '0' COMMENT '是否启用双因子认证',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `created_by` bigint(20) DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '更新者',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 角色表
CREATE TABLE IF NOT EXISTS `roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `data_scope` tinyint(1) DEFAULT '1' COMMENT '数据范围：1=全部 2=自定义 3=部门 4=仅本人',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态：1=正常 2=禁用',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `created_by` bigint(20) DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '更新者',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 菜单表
CREATE TABLE IF NOT EXISTS `menus` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id` bigint(20) DEFAULT '0' COMMENT '父菜单ID',
  `title` varchar(100) NOT NULL COMMENT '菜单标题',
  `name` varchar(100) NOT NULL COMMENT '菜单名称',
  `path` varchar(200) DEFAULT NULL COMMENT '路由路径',
  `component` varchar(200) DEFAULT NULL COMMENT '组件路径',
  `icon` varchar(100) DEFAULT NULL COMMENT '菜单图标',
  `type` tinyint(1) DEFAULT '1' COMMENT '菜单类型：1=目录 2=菜单 3=按钮',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态：1=正常 2=禁用',
  `visible` tinyint(1) DEFAULT '1' COMMENT '是否可见：1=可见 2=隐藏',
  `permission` varchar(100) DEFAULT NULL COMMENT '权限标识',
  `redirect` varchar(200) DEFAULT NULL COMMENT '重定向路径',
  `always_show` tinyint(1) DEFAULT '0' COMMENT '是否总是显示：1=是 0=否',
  `breadcrumb` tinyint(1) DEFAULT '1' COMMENT '是否显示面包屑：1=显示 0=隐藏',
  `active_menu` varchar(200) DEFAULT NULL COMMENT '激活菜单',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `created_by` bigint(20) DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '更新者',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS `user_roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `is_primary` tinyint(1) DEFAULT '0' COMMENT '是否主要角色：1=是 0=否',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`,`role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 角色菜单关联表
CREATE TABLE IF NOT EXISTS `role_menus` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `menu_id` bigint(20) NOT NULL COMMENT '菜单ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_menu` (`role_id`,`menu_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';

-- 插入默认管理员用户
INSERT INTO `users` (`id`, `username`, `password`, `salt`, `email`, `real_name`, `nickname`, `status`, `created_by`) VALUES
(1, 'admin', '7fef6171469e80d32c0559f88b377245', 'admin', 'admin@example.com', '系统管理员', '管理员', 1, 0);

-- 插入默认角色
INSERT INTO `roles` (`id`, `name`, `code`, `data_scope`, `status`, `sort`, `remark`, `created_by`) VALUES
(1, '超级管理员', 'super_admin', 1, 1, 1, '超级管理员，拥有所有权限', 0),
(2, '系统管理员', 'admin', 1, 1, 2, '系统管理员', 0),
(3, '普通用户', 'user', 4, 1, 3, '普通用户', 0);

-- 插入默认菜单
INSERT INTO `menus` (`id`, `parent_id`, `title`, `name`, `path`, `component`, `icon`, `type`, `sort`, `status`, `visible`, `permission`, `redirect`, `always_show`, `breadcrumb`, `remark`, `created_by`) VALUES
(1, 0, '系统管理', 'System', '/system', 'Layout', 'system', 1, 1, 1, 1, 'system', '/system/user', 1, 1, '系统管理目录', 0),
(2, 1, '用户管理', 'User', '/system/user', 'system/user/index', 'user', 2, 1, 1, 1, 'system:user:list', '', 0, 1, '用户管理菜单', 0),
(3, 1, '角色管理', 'Role', '/system/role', 'system/role/index', 'role', 2, 2, 1, 1, 'system:role:list', '', 0, 1, '角色管理菜单', 0),
(4, 1, '菜单管理', 'Menu', '/system/menu', 'system/menu/index', 'menu', 2, 3, 1, 1, 'system:menu:list', '', 0, 1, '菜单管理菜单', 0),
(5, 0, '仪表盘', 'Dashboard', '/dashboard', 'dashboard/index', 'dashboard', 2, 0, 1, 1, 'dashboard', '', 0, 1, '仪表盘', 0);

-- 分配用户角色
INSERT INTO `user_roles` (`user_id`, `role_id`, `is_primary`) VALUES
(1, 1, 1);

-- 分配角色菜单权限
INSERT INTO `role_menus` (`role_id`, `menu_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5),
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5),
(3, 5); 