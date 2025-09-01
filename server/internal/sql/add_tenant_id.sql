-- 为现有表添加 tenant_id 字段

-- 用户表添加租户ID
ALTER TABLE `sys_users` ADD COLUMN `tenant_id` bigint(20) unsigned DEFAULT 1 COMMENT '租户ID' AFTER `id`;
ALTER TABLE `sys_users` ADD INDEX `idx_tenant_id` (`tenant_id`);
ALTER TABLE `sys_users` ADD INDEX `idx_tenant_username` (`tenant_id`, `username`);

-- 角色表添加租户ID
ALTER TABLE `sys_roles` ADD COLUMN `tenant_id` bigint(20) unsigned DEFAULT 1 COMMENT '租户ID' AFTER `id`;
ALTER TABLE `sys_roles` ADD INDEX `idx_tenant_id` (`tenant_id`);
ALTER TABLE `sys_roles` ADD INDEX `idx_tenant_code` (`tenant_id`, `code`);

-- 菜单表保持全局共享，不添加租户ID
-- 菜单是系统级别的，所有租户共享同一套菜单结构
-- 租户权限通过角色-菜单关联表控制

-- 用户角色关联表添加租户ID
ALTER TABLE `sys_user_roles` ADD COLUMN `tenant_id` bigint(20) unsigned DEFAULT 1 COMMENT '租户ID' AFTER `id`;
ALTER TABLE `sys_user_roles` ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 角色菜单关联表添加租户ID
ALTER TABLE `sys_role_menus` ADD COLUMN `tenant_id` bigint(20) unsigned DEFAULT 1 COMMENT '租户ID' AFTER `id`;
ALTER TABLE `sys_role_menus` ADD INDEX `idx_tenant_id` (`tenant_id`);

-- 修改现有唯一索引，加入租户隔离
-- 用户表：同一租户内用户名唯一
ALTER TABLE `sys_users` DROP INDEX `uk_username`;
ALTER TABLE `sys_users` ADD UNIQUE KEY `uk_tenant_username` (`tenant_id`, `username`);

-- 角色表：同一租户内角色编码唯一
ALTER TABLE `sys_roles` DROP INDEX `uk_code`;
ALTER TABLE `sys_roles` ADD UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`);

-- 用户角色关联表：同一租户内用户角色关联唯一
ALTER TABLE `sys_user_roles` DROP INDEX `uk_user_role`;
ALTER TABLE `sys_user_roles` ADD UNIQUE KEY `uk_tenant_user_role` (`tenant_id`, `user_id`, `role_id`);

-- 角色菜单关联表：同一租户内角色菜单关联唯一
ALTER TABLE `sys_role_menus` DROP INDEX `uk_role_menu`;
ALTER TABLE `sys_role_menus` ADD UNIQUE KEY `uk_tenant_role_menu` (`tenant_id`, `role_id`, `menu_id`);
