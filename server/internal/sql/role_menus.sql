-- 创建角色菜单关联表
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

-- 为超级管理员分配所有菜单权限
INSERT INTO `role_menus` (`role_id`, `menu_id`, `created_at`)
SELECT 1, id, NOW() FROM menus WHERE status = 1;

-- 为系统管理员分配系统管理权限
INSERT INTO `role_menus` (`role_id`, `menu_id`, `created_at`) VALUES
-- 系统管理模块
(2, 1, NOW()),   -- 系统管理目录
(2, 2, NOW()),   -- 菜单管理
(2, 3, NOW()),   -- 查看菜单
(2, 4, NOW()),   -- 新增菜单
(2, 5, NOW()),   -- 修改菜单
(2, 6, NOW()),   -- 删除菜单
-- 系统工具模块
(2, 40, NOW()),  -- 系统工具目录
(2, 41, NOW()),  -- 系统日志
(2, 42, NOW()),  -- 查看日志
(2, 43, NOW()),  -- 删除日志
(2, 44, NOW()),  -- 清空日志
(2, 50, NOW()),  -- 配置管理
(2, 51, NOW()),  -- 查看配置
(2, 52, NOW());  -- 修改配置



-- 为财务管理员分配财务统计权限
INSERT INTO `role_menus` (`role_id`, `menu_id`, `created_at`) VALUES
-- 系统日志（只读）
(4, 40, NOW()),  -- 系统工具目录
(4, 41, NOW()),  -- 系统日志
(4, 42, NOW());  -- 查看日志

-- 为运营人员分配基础操作权限
INSERT INTO `role_menus` (`role_id`, `menu_id`, `created_at`) VALUES
-- 系统日志（只读）
(5, 40, NOW()),  -- 系统工具目录
(5, 41, NOW()),  -- 系统日志
(5, 42, NOW());  -- 查看日志

-- 为客服人员分配基础权限
INSERT INTO `role_menus` (`role_id`, `menu_id`, `created_at`) VALUES
-- 系统日志（只读）
(6, 40, NOW()),  -- 系统工具目录
(6, 41, NOW()),  -- 系统日志
(6, 42, NOW());  -- 查看日志

-- 为审计人员分配只读权限
INSERT INTO `role_menus` (`role_id`, `menu_id`, `created_at`) VALUES
-- 系统管理模块（只读）
(7, 1, NOW()),   -- 系统管理目录
(7, 2, NOW()),   -- 菜单管理
(7, 3, NOW()),   -- 查看菜单
-- 系统工具模块（只读）
(7, 40, NOW()),  -- 系统工具目录
(7, 41, NOW()),  -- 系统日志
(7, 42, NOW()),  -- 查看日志
(7, 50, NOW()),  -- 配置管理
(7, 51, NOW());  -- 查看配置