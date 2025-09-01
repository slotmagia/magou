-- 创建菜单表
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

-- 插入默认菜单数据
INSERT INTO `menus` (`id`, `parent_id`, `title`, `name`, `path`, `component`, `icon`, `type`, `sort`, `status`, `visible`, `permission`, `redirect`, `always_show`, `breadcrumb`, `active_menu`, `remark`, `created_by`, `updated_by`, `created_at`, `updated_at`) VALUES
(1, 0, '系统管理', 'System', '/system', 'Layout', 'system', 1, 1000, 1, 1, 'system', '/system/menu', 1, 1, NULL, '系统管理目录', 1, 1, NOW(), NOW()),
(2, 1, '菜单管理', 'Menu', '/system/menu', 'system/menu/index', 'tree-table', 2, 1, 1, 1, 'system:menu:list', NULL, 0, 1, NULL, '菜单管理页面', 1, 1, NOW(), NOW()),
(3, 2, '查看菜单', 'MenuView', '', NULL, NULL, 3, 1, 1, 0, 'system:menu:query', NULL, 0, 1, NULL, '查看菜单权限', 1, 1, NOW(), NOW()),
(4, 2, '新增菜单', 'MenuAdd', '', NULL, NULL, 3, 2, 1, 0, 'system:menu:add', NULL, 0, 1, NULL, '新增菜单权限', 1, 1, NOW(), NOW()),
(5, 2, '修改菜单', 'MenuEdit', '', NULL, NULL, 3, 3, 1, 0, 'system:menu:edit', NULL, 0, 1, NULL, '修改菜单权限', 1, 1, NOW(), NOW()),
(6, 2, '删除菜单', 'MenuDelete', '', NULL, NULL, 3, 4, 1, 0, 'system:menu:remove', NULL, 0, 1, NULL, '删除菜单权限', 1, 1, NOW(), NOW()),



(40, 0, '系统工具', 'Tool', '/tool', 'Layout', 'tool', 1, 3000, 1, 1, 'tool', '/tool/log', 1, 1, NULL, '系统工具目录', 1, 1, NOW(), NOW()),
(41, 40, '系统日志', 'Log', '/tool/log', 'tool/log/index', 'documentation', 2, 1, 1, 1, 'tool:log:list', NULL, 0, 1, NULL, '系统日志管理', 1, 1, NOW(), NOW()),
(42, 41, '查看日志', 'LogView', '', NULL, NULL, 3, 1, 1, 0, 'tool:log:query', NULL, 0, 1, NULL, '查看日志权限', 1, 1, NOW(), NOW()),
(43, 41, '删除日志', 'LogDelete', '', NULL, NULL, 3, 2, 1, 0, 'tool:log:remove', NULL, 0, 1, NULL, '删除日志权限', 1, 1, NOW(), NOW()),
(44, 41, '清空日志', 'LogClear', '', NULL, NULL, 3, 3, 1, 0, 'tool:log:clear', NULL, 0, 1, NULL, '清空日志权限', 1, 1, NOW(), NOW()),

(50, 40, '配置管理', 'Config', '/tool/config', 'tool/config/index', 'edit', 2, 2, 1, 1, 'tool:config:list', NULL, 0, 1, NULL, '系统配置管理', 1, 1, NOW(), NOW()),
(51, 50, '查看配置', 'ConfigView', '', NULL, NULL, 3, 1, 1, 0, 'tool:config:query', NULL, 0, 1, NULL, '查看配置权限', 1, 1, NOW(), NOW()),
(52, 50, '修改配置', 'ConfigEdit', '', NULL, NULL, 3, 2, 1, 0, 'tool:config:edit', NULL, 0, 1, NULL, '修改配置权限', 1, 1, NOW(), NOW()); 