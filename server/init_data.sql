/*
 多租户管理系统 - 初始化数据
 基于 admin.sql 表结构创建
 
 执行顺序：
 1. 租户数据
 2. 部门数据  
 3. 系统菜单数据
 4. 角色数据
 5. 用户数据
 6. 权限关联数据
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ===========================
-- 1. 租户数据初始化
-- ===========================

-- 插入系统租户（ID=1，系统默认租户）
INSERT INTO `sys_tenants` (`id`, `name`, `code`, `domain`, `status`, `max_users`, `storage_limit`, `expire_at`, `admin_user_id`, `config`, `remark`, `created_by`, `updated_by`, `created_at`, `updated_at`) VALUES
(1, '系统租户', 'system', 'admin.example.com', 1, 999999, 10737418240, NULL, 1, 
JSON_OBJECT(
  'features', JSON_OBJECT(
    'advancedReports', true,
    'apiAccess', true, 
    'customBranding', true
  ),
  'limitations', JSON_OBJECT(
    'maxApiCalls', 100000,
    'maxStorage', 10737418240
  ),
  'settings', JSON_OBJECT(
    'theme', 'light',
    'language', 'zh-CN',
    'timezone', 'Asia/Shanghai'
  )
), '系统默认租户，不可删除', 0, 0, NOW(), NOW());

-- 插入演示租户
INSERT INTO `sys_tenants` (`id`, `name`, `code`, `domain`, `status`, `max_users`, `storage_limit`, `expire_at`, `admin_user_id`, `config`, `remark`, `created_by`, `updated_by`, `created_at`, `updated_at`) VALUES
(2, '演示租户', 'demo', 'demo.example.com', 1, 100, 1073741824, '2025-12-31 23:59:59', 2,
JSON_OBJECT(
  'features', JSON_OBJECT(
    'advancedReports', false,
    'apiAccess', true,
    'customBranding', false
  ),
  'limitations', JSON_OBJECT(
    'maxApiCalls', 10000,
    'maxStorage', 1073741824
  ),
  'settings', JSON_OBJECT(
    'theme', 'dark',
    'language', 'zh-CN', 
    'timezone', 'Asia/Shanghai'
  )
), '演示环境租户', 1, 1, NOW(), NOW());

-- ===========================
-- 2. 部门数据初始化
-- ===========================

-- 系统租户部门
INSERT INTO `sys_department` (`id`, `tenant_id`, `parent_id`, `dept_code`, `dept_name`, `dept_level`, `dept_path`, `sort_order`, `leader`, `phone`, `email`, `status`, `created_at`, `updated_at`, `remark`) VALUES
(1, 1, 0, 'ROOT', '总公司', 1, '/1/', 1, '系统管理员', '400-888-8888', 'admin@system.com', 1, NOW(), NOW(), '系统租户根部门'),
(2, 1, 1, 'TECH', '技术部', 2, '/1/2/', 1, '技术总监', '400-888-8801', 'tech@system.com', 1, NOW(), NOW(), '技术开发部门'),
(3, 1, 1, 'FINANCE', '财务部', 2, '/1/3/', 2, '财务经理', '400-888-8802', 'finance@system.com', 1, NOW(), NOW(), '财务管理部门');

-- 演示租户部门
INSERT INTO `sys_department` (`id`, `tenant_id`, `parent_id`, `dept_code`, `dept_name`, `dept_level`, `dept_path`, `sort_order`, `leader`, `phone`, `email`, `status`, `created_at`, `updated_at`, `remark`) VALUES
(4, 2, 0, 'ROOT', '演示公司', 1, '/4/', 1, '演示管理员', '400-999-9999', 'admin@demo.com', 1, NOW(), NOW(), '演示租户根部门'),
(5, 2, 4, 'SALES', '销售部', 2, '/4/5/', 1, '销售经理', '400-999-9901', 'sales@demo.com', 1, NOW(), NOW(), '销售业务部门');

-- ===========================
-- 3. 系统菜单数据初始化  
-- ===========================

-- 系统菜单（tenant_id=NULL，全局共享）
INSERT INTO `sys_menus` (`id`, `tenant_id`, `parent_id`, `menu_code`, `menu_name`, `menu_type`, `icon`, `path`, `component`, `permission`, `sort_order`, `visible`, `status`, `created_at`, `updated_at`, `remark`) VALUES
-- 一级菜单
(1, NULL, 0, 'dashboard', '仪表盘', 2, 'dashboard', '/dashboard', 'dashboard/index', 'dashboard', 0, 1, 1, NOW(), NOW(), '系统仪表盘'),
(2, NULL, 0, 'system', '系统管理', 1, 'system', '/system', 'Layout', 'system', 1000, 1, 1, NOW(), NOW(), '系统管理目录'),
(3, NULL, 0, 'tenant', '租户管理', 1, 'peoples', '/tenant', 'Layout', 'tenant', 2000, 1, 1, NOW(), NOW(), '租户管理目录'),

-- 系统管理子菜单
(10, NULL, 2, 'user', '用户管理', 2, 'user', '/system/user', 'system/user/index', 'system:user:list', 1, 1, 1, NOW(), NOW(), '用户管理页面'),
(11, NULL, 2, 'role', '角色管理', 2, 'peoples', '/system/role', 'system/role/index', 'system:role:list', 2, 1, 1, NOW(), NOW(), '角色管理页面'),
(12, NULL, 2, 'menu', '菜单管理', 2, 'tree-table', '/system/menu', 'system/menu/index', 'system:menu:list', 3, 1, 1, NOW(), NOW(), '菜单管理页面'),
(13, NULL, 2, 'dept', '部门管理', 2, 'tree', '/system/dept', 'system/dept/index', 'system:dept:list', 4, 1, 1, NOW(), NOW(), '部门管理页面'),

-- 租户管理子菜单
(20, NULL, 3, 'tenant_list', '租户列表', 2, 'list', '/tenant/list', 'tenant/list/index', 'tenant:list', 1, 1, 1, NOW(), NOW(), '租户列表管理'),
(21, NULL, 3, 'tenant_config', '租户配置', 2, 'edit', '/tenant/config', 'tenant/config/index', 'tenant:config', 2, 1, 1, NOW(), NOW(), '租户配置管理'),

-- 用户管理按钮权限
(100, NULL, 10, 'user_add', '新增用户', 3, '', '', '', 'system:user:add', 1, 0, 1, NOW(), NOW(), '新增用户按钮'),
(101, NULL, 10, 'user_edit', '编辑用户', 3, '', '', '', 'system:user:edit', 2, 0, 1, NOW(), NOW(), '编辑用户按钮'),
(102, NULL, 10, 'user_delete', '删除用户', 3, '', '', '', 'system:user:delete', 3, 0, 1, NOW(), NOW(), '删除用户按钮'),
(103, NULL, 10, 'user_reset_pwd', '重置密码', 3, '', '', '', 'system:user:resetPwd', 4, 0, 1, NOW(), NOW(), '重置密码按钮'),

-- 角色管理按钮权限
(110, NULL, 11, 'role_add', '新增角色', 3, '', '', '', 'system:role:add', 1, 0, 1, NOW(), NOW(), '新增角色按钮'),
(111, NULL, 11, 'role_edit', '编辑角色', 3, '', '', '', 'system:role:edit', 2, 0, 1, NOW(), NOW(), '编辑角色按钮'),
(112, NULL, 11, 'role_delete', '删除角色', 3, '', '', '', 'system:role:delete', 3, 0, 1, NOW(), NOW(), '删除角色按钮'),
(113, NULL, 11, 'role_permission', '分配权限', 3, '', '', '', 'system:role:permission', 4, 0, 1, NOW(), NOW(), '分配权限按钮'),

-- 菜单管理按钮权限
(120, NULL, 12, 'menu_add', '新增菜单', 3, '', '', '', 'system:menu:add', 1, 0, 1, NOW(), NOW(), '新增菜单按钮'),
(121, NULL, 12, 'menu_edit', '编辑菜单', 3, '', '', '', 'system:menu:edit', 2, 0, 1, NOW(), NOW(), '编辑菜单按钮'),
(122, NULL, 12, 'menu_delete', '删除菜单', 3, '', '', '', 'system:menu:delete', 3, 0, 1, NOW(), NOW(), '删除菜单按钮'),

-- 部门管理按钮权限
(130, NULL, 13, 'dept_add', '新增部门', 3, '', '', '', 'system:dept:add', 1, 0, 1, NOW(), NOW(), '新增部门按钮'),
(131, NULL, 13, 'dept_edit', '编辑部门', 3, '', '', '', 'system:dept:edit', 2, 0, 1, NOW(), NOW(), '编辑部门按钮'),
(132, NULL, 13, 'dept_delete', '删除部门', 3, '', '', '', 'system:dept:delete', 3, 0, 1, NOW(), NOW(), '删除部门按钮'),

-- 租户管理按钮权限
(200, NULL, 20, 'tenant_add', '新增租户', 3, '', '', '', 'tenant:add', 1, 0, 1, NOW(), NOW(), '新增租户按钮'),
(201, NULL, 20, 'tenant_edit', '编辑租户', 3, '', '', '', 'tenant:edit', 2, 0, 1, NOW(), NOW(), '编辑租户按钮'),
(202, NULL, 20, 'tenant_delete', '删除租户', 3, '', '', '', 'tenant:delete', 3, 0, 1, NOW(), NOW(), '删除租户按钮'),
(203, NULL, 20, 'tenant_status', '状态管理', 3, '', '', '', 'tenant:status', 4, 0, 1, NOW(), NOW(), '租户状态管理');

-- ===========================
-- 4. 角色数据初始化
-- ===========================

-- 系统租户角色
INSERT INTO `sys_roles` (`id`, `tenant_id`, `role_code`, `role_name`, `role_type`, `data_scope`, `status`, `created_at`, `updated_at`, `remark`) VALUES
(1, 1, 'super_admin', '超级管理员', 2, 1, 1, NOW(), NOW(), '系统超级管理员，拥有所有权限'),
(2, 1, 'system_admin', '系统管理员', 2, 2, 1, NOW(), NOW(), '系统管理员，负责系统配置'),
(3, 1, 'tenant_admin', '租户管理员', 2, 2, 1, NOW(), NOW(), '租户管理员，负责租户管理');

-- 演示租户角色
INSERT INTO `sys_roles` (`id`, `tenant_id`, `role_code`, `role_name`, `role_type`, `data_scope`, `status`, `created_at`, `updated_at`, `remark`) VALUES
(4, 2, 'demo_admin', '演示管理员', 2, 1, 1, NOW(), NOW(), '演示租户管理员'),
(5, 2, 'demo_user', '演示用户', 1, 4, 1, NOW(), NOW(), '演示租户普通用户');

-- ===========================
-- 5. 用户数据初始化
-- ===========================

-- 系统租户用户（密码: 123456）
INSERT INTO `sys_users` (`id`, `tenant_id`, `username`, `email`, `phone`, `password`, `salt`, `real_name`, `nickname`, `avatar`, `gender`, `birthday`, `dept_id`, `position`, `status`, `login_ip`, `login_at`, `login_count`, `email_verified_at`, `remark`, `created_by`, `updated_by`, `created_at`, `updated_at`) VALUES
(1, 1, 'admin', 'admin@system.com', '13800138000', 'e10adc3949ba59abbe56e057f20f883e', 'salt123', '系统管理员', '超管', '/avatars/admin.jpg', 1, '1990-01-01', 1, 'CTO', 1, '127.0.0.1', NOW(), 1, NOW(), '系统超级管理员账号', 0, 0, NOW(), NOW()),
(2, 1, 'sysadmin', 'sysadmin@system.com', '13800138001', 'e10adc3949ba59abbe56e057f20f883e', 'salt123', '系统管理', '系管', '/avatars/sysadmin.jpg', 1, '1985-05-15', 2, '系统管理员', 1, NULL, NULL, 0, NOW(), '系统管理员账号', 1, 1, NOW(), NOW()),
(3, 1, 'tenant_admin', 'tadmin@system.com', '13800138002', 'e10adc3949ba59abbe56e057f20f883e', 'salt123', '租户管理', '租管', '/avatars/tenant.jpg', 2, '1988-08-20', 2, '租户管理员', 1, NULL, NULL, 0, NOW(), '租户管理员账号', 1, 1, NOW(), NOW());

-- 演示租户用户（密码: 123456）  
INSERT INTO `sys_users` (`id`, `tenant_id`, `username`, `email`, `phone`, `password`, `salt`, `real_name`, `nickname`, `avatar`, `gender`, `birthday`, `dept_id`, `position`, `status`, `login_ip`, `login_at`, `login_count`, `email_verified_at`, `remark`, `created_by`, `updated_by`, `created_at`, `updated_at`) VALUES
(4, 2, 'demo_admin', 'admin@demo.com', '13900139000', 'e10adc3949ba59abbe56e057f20f883e', 'salt123', '演示管理员', '演管', '/avatars/demo_admin.jpg', 1, '1992-03-10', 4, '总经理', 1, NULL, NULL, 0, NOW(), '演示租户管理员', 1, 1, NOW(), NOW()),
(5, 2, 'demo_user', 'user@demo.com', '13900139001', 'e10adc3949ba59abbe56e057f20f883e', 'salt123', '演示用户', '用户', '/avatars/demo_user.jpg', 2, '1995-07-25', 5, '销售专员', 1, NULL, NULL, 0, NOW(), '演示租户普通用户', 4, 4, NOW(), NOW());

-- 更新租户的管理员用户ID
UPDATE `sys_tenants` SET `admin_user_id` = 1 WHERE `id` = 1;
UPDATE `sys_tenants` SET `admin_user_id` = 4 WHERE `id` = 2;

-- ===========================
-- 6. 用户角色关联数据
-- ===========================

INSERT INTO `sys_user_roles` (`tenant_id`, `user_id`, `role_id`, `is_primary`, `created_at`, `updated_at`) VALUES
-- 系统租户用户角色关联
(1, 1, 1, 1, NOW(), NOW()), -- admin -> 超级管理员
(1, 2, 2, 1, NOW(), NOW()), -- sysadmin -> 系统管理员  
(1, 3, 3, 1, NOW(), NOW()), -- tenant_admin -> 租户管理员

-- 演示租户用户角色关联
(2, 4, 4, 1, NOW(), NOW()), -- demo_admin -> 演示管理员
(2, 5, 5, 1, NOW(), NOW()); -- demo_user -> 演示用户

-- ===========================
-- 7. 角色菜单权限关联数据
-- ===========================

INSERT INTO `sys_role_menus` (`tenant_id`, `role_id`, `menu_id`, `created_at`) VALUES
-- 超级管理员（ID=1）- 拥有所有菜单权限
(1, 1, 1, NOW()), (1, 1, 2, NOW()), (1, 1, 3, NOW()),
(1, 1, 10, NOW()), (1, 1, 11, NOW()), (1, 1, 12, NOW()), (1, 1, 13, NOW()),
(1, 1, 20, NOW()), (1, 1, 21, NOW()),
(1, 1, 100, NOW()), (1, 1, 101, NOW()), (1, 1, 102, NOW()), (1, 1, 103, NOW()),
(1, 1, 110, NOW()), (1, 1, 111, NOW()), (1, 1, 112, NOW()), (1, 1, 113, NOW()),
(1, 1, 120, NOW()), (1, 1, 121, NOW()), (1, 1, 122, NOW()),
(1, 1, 130, NOW()), (1, 1, 131, NOW()), (1, 1, 132, NOW()),
(1, 1, 200, NOW()), (1, 1, 201, NOW()), (1, 1, 202, NOW()), (1, 1, 203, NOW()),

-- 系统管理员（ID=2）- 系统管理相关权限
(1, 2, 1, NOW()), (1, 2, 2, NOW()),
(1, 2, 10, NOW()), (1, 2, 11, NOW()), (1, 2, 12, NOW()), (1, 2, 13, NOW()),
(1, 2, 100, NOW()), (1, 2, 101, NOW()), (1, 2, 102, NOW()), (1, 2, 103, NOW()),
(1, 2, 110, NOW()), (1, 2, 111, NOW()), (1, 2, 112, NOW()), (1, 2, 113, NOW()),
(1, 2, 120, NOW()), (1, 2, 121, NOW()), (1, 2, 122, NOW()),
(1, 2, 130, NOW()), (1, 2, 131, NOW()), (1, 2, 132, NOW()),

-- 租户管理员（ID=3）- 租户管理相关权限
(1, 3, 1, NOW()), (1, 3, 3, NOW()),
(1, 3, 20, NOW()), (1, 3, 21, NOW()),
(1, 3, 200, NOW()), (1, 3, 201, NOW()), (1, 3, 202, NOW()), (1, 3, 203, NOW()),

-- 演示管理员（ID=4）- 演示租户管理权限
(2, 4, 1, NOW()), (2, 4, 2, NOW()),
(2, 4, 10, NOW()), (2, 4, 11, NOW()), (2, 4, 12, NOW()), (2, 4, 13, NOW()),
(2, 4, 100, NOW()), (2, 4, 101, NOW()), (2, 4, 102, NOW()),
(2, 4, 110, NOW()), (2, 4, 111, NOW()), (2, 4, 112, NOW()), (2, 4, 113, NOW()),
(2, 4, 130, NOW()), (2, 4, 131, NOW()),

-- 演示用户（ID=5）- 基础查看权限
(2, 5, 1, NOW()), (2, 5, 2, NOW()),
(2, 5, 10, NOW()), (2, 5, 11, NOW()), (2, 5, 13, NOW());

-- ===========================
-- 8. 系统权限数据初始化
-- ===========================

INSERT INTO `sys_permission` (`tenant_id`, `permission_code`, `permission_name`, `permission_type`, `resource_path`, `method`, `status`, `created_at`, `updated_at`, `remark`) VALUES
-- 系统级API权限（tenant_id=NULL）
(NULL, 'api:user:login', '用户登录', 3, '/api/login', 'POST', 1, NOW(), NOW(), '用户登录接口'),
(NULL, 'api:user:logout', '用户登出', 3, '/api/logout', 'POST', 1, NOW(), NOW(), '用户登出接口'),
(NULL, 'api:user:profile', '获取用户信息', 3, '/api/profile', 'GET', 1, NOW(), NOW(), '获取用户信息接口'),
(NULL, 'api:user:captcha', '获取验证码', 3, '/api/captcha', 'GET', 1, NOW(), NOW(), '获取验证码接口'),

-- 租户管理API权限
(NULL, 'api:tenant:list', '租户列表', 3, '/api/tenant/list', 'GET', 1, NOW(), NOW(), '获取租户列表'),
(NULL, 'api:tenant:create', '创建租户', 3, '/api/tenant/create', 'POST', 1, NOW(), NOW(), '创建租户'),
(NULL, 'api:tenant:update', '更新租户', 3, '/api/tenant/update', 'PUT', 1, NOW(), NOW(), '更新租户'),
(NULL, 'api:tenant:delete', '删除租户', 3, '/api/tenant/delete', 'DELETE', 1, NOW(), NOW(), '删除租户'),

-- 用户管理API权限
(NULL, 'api:user:list', '用户列表', 3, '/api/user/list', 'GET', 1, NOW(), NOW(), '获取用户列表'),
(NULL, 'api:user:create', '创建用户', 3, '/api/user/create', 'POST', 1, NOW(), NOW(), '创建用户'),
(NULL, 'api:user:update', '更新用户', 3, '/api/user/update', 'PUT', 1, NOW(), NOW(), '更新用户'),
(NULL, 'api:user:delete', '删除用户', 3, '/api/user/delete', 'DELETE', 1, NOW(), NOW(), '删除用户'),

-- 角色管理API权限
(NULL, 'api:role:list', '角色列表', 3, '/api/role/list', 'GET', 1, NOW(), NOW(), '获取角色列表'),
(NULL, 'api:role:create', '创建角色', 3, '/api/role', 'POST', 1, NOW(), NOW(), '创建角色'),
(NULL, 'api:role:update', '更新角色', 3, '/api/role/{id}', 'PUT', 1, NOW(), NOW(), '更新角色'),
(NULL, 'api:role:delete', '删除角色', 3, '/api/role/{id}', 'DELETE', 1, NOW(), NOW(), '删除角色'),
(NULL, 'api:role:menus', '角色菜单', 3, '/api/role/{id}/menus', 'GET', 1, NOW(), NOW(), '获取角色菜单'),

-- 菜单管理API权限
(NULL, 'api:menu:list', '菜单列表', 3, '/api/menu/list', 'GET', 1, NOW(), NOW(), '获取菜单列表'),
(NULL, 'api:menu:tree', '菜单树', 3, '/api/menu/tree', 'GET', 1, NOW(), NOW(), '获取菜单树'),
(NULL, 'api:menu:create', '创建菜单', 3, '/api/menu', 'POST', 1, NOW(), NOW(), '创建菜单'),
(NULL, 'api:menu:update', '更新菜单', 3, '/api/menu/{id}', 'PUT', 1, NOW(), NOW(), '更新菜单'),
(NULL, 'api:menu:delete', '删除菜单', 3, '/api/menu/{id}', 'DELETE', 1, NOW(), NOW(), '删除菜单');

-- ===========================
-- 9. 字段权限数据初始化
-- ===========================

INSERT INTO `sys_field_permission` (`tenant_id`, `resource_type`, `field_code`, `field_name`, `field_type`, `mask_type`, `mask_pattern`, `status`, `created_at`, `updated_at`, `remark`) VALUES
-- 用户表字段权限
(1, 'sys_users', 'phone', '手机号码', 'varchar', 1, '***-****-{4}', 1, NOW(), NOW(), '手机号中间遮掩'),
(1, 'sys_users', 'email', '邮箱地址', 'varchar', 1, '{3}***@{domain}', 1, NOW(), NOW(), '邮箱前缀遮掩'),
(1, 'sys_users', 'password', '密码', 'varchar', 2, '********', 1, NOW(), NOW(), '密码完全遮掩'),

-- 演示租户字段权限
(2, 'sys_users', 'phone', '手机号码', 'varchar', 1, '***-****-{4}', 1, NOW(), NOW(), '手机号中间遮掩'),
(2, 'sys_users', 'email', '邮箱地址', 'varchar', 1, '{3}***@{domain}', 1, NOW(), NOW(), '邮箱前缀遮掩');

SET FOREIGN_KEY_CHECKS = 1;

-- ===========================
-- 初始化完成提示
-- ===========================

SELECT 
  '数据初始化完成！' as message,
  '默认系统管理员: admin / 123456' as system_admin,
  '默认演示管理员: demo_admin / 123456' as demo_admin,
  '访问地址: http://localhost:8000' as access_url;

