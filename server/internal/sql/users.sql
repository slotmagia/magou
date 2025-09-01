-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱地址',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号码',
  `password` varchar(255) NOT NULL COMMENT '密码（加密存储）',
  `salt` varchar(32) DEFAULT NULL COMMENT '密码盐值',
  `real_name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
  `gender` tinyint(4) DEFAULT '0' COMMENT '性别：0=未知 1=男 2=女',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `dept_id` bigint(20) unsigned DEFAULT NULL COMMENT '部门ID',
  `position` varchar(100) DEFAULT NULL COMMENT '职位',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1=正常 2=锁定 3=禁用',
  `login_ip` varchar(45) DEFAULT NULL COMMENT '最后登录IP',
  `login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `login_count` int(11) NOT NULL DEFAULT '0' COMMENT '登录次数',
  `password_reset_token` varchar(255) DEFAULT NULL COMMENT '密码重置令牌',
  `password_reset_expires` datetime DEFAULT NULL COMMENT '密码重置过期时间',
  `email_verified_at` datetime DEFAULT NULL COMMENT '邮箱验证时间',
  `phone_verified_at` datetime DEFAULT NULL COMMENT '手机验证时间',
  `two_factor_enabled` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否启用双因子认证：1=是 0=否',
  `two_factor_secret` varchar(255) DEFAULT NULL COMMENT '双因子认证密钥',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注说明',
  `created_by` bigint(20) unsigned DEFAULT NULL COMMENT '创建人ID',
  `updated_by` bigint(20) unsigned DEFAULT NULL COMMENT '修改人ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间（软删除）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  UNIQUE KEY `uk_phone` (`phone`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_login_at` (`login_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 插入默认管理员用户
INSERT INTO `users` (`id`, `username`, `email`, `password`, `salt`, `real_name`, `nickname`, `status`, `dept_id`, `position`, `created_by`, `updated_by`, `created_at`, `updated_at`) VALUES
(1, 'admin', 'admin@example.com', 'e10adc3949ba59abbe56e057f20f883e', 'salt123', '系统管理员', '管理员', 1, 1, '超级管理员', 1, 1, NOW(), NOW());



-- 插入测试用户（财务管理员）
INSERT INTO `users` (`id`, `username`, `email`, `password`, `salt`, `real_name`, `nickname`, `status`, `dept_id`, `position`, `created_by`, `updated_by`, `created_at`, `updated_at`) VALUES
(3, 'finance_admin', 'finance@example.com', 'e10adc3949ba59abbe56e057f20f883e', 'salt123', '李四', '财务管理', 1, 3, '财务管理员', 1, 1, NOW(), NOW());

-- 插入测试用户（运营人员）
INSERT INTO `users` (`id`, `username`, `email`, `password`, `salt`, `real_name`, `nickname`, `status`, `dept_id`, `position`, `created_by`, `updated_by`, `created_at`, `updated_at`) VALUES
(4, 'operator', 'operator@example.com', 'e10adc3949ba59abbe56e057f20f883e', 'salt123', '王五', '运营专员', 1, 4, '运营人员', 1, 1, NOW(), NOW());

-- 插入测试用户（客服人员）
INSERT INTO `users` (`id`, `username`, `email`, `password`, `salt`, `real_name`, `nickname`, `status`, `dept_id`, `position`, `created_by`, `updated_by`, `created_at`, `updated_at`) VALUES
(5, 'customer_service', 'service@example.com', 'e10adc3949ba59abbe56e057f20f883e', 'salt123', '赵六', '客服专员', 1, 5, '客服人员', 1, 1, NOW(), NOW());

-- 补充用户角色关联数据
INSERT INTO `user_roles` (`user_id`, `role_id`, `is_primary`, `assigned_by`, `created_at`, `updated_at`) VALUES

(3, 4, 1, 1, NOW(), NOW()), -- 财务管理员角色
(4, 5, 1, 1, NOW(), NOW()), -- 运营人员角色
(5, 6, 1, 1, NOW(), NOW()); -- 客服人员角色 