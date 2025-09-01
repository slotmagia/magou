-- 创建用户角色关联表
CREATE TABLE IF NOT EXISTS `user_roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
  `is_primary` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否主要角色：1=是 0=否',
  `assigned_by` bigint(20) unsigned DEFAULT NULL COMMENT '分配人ID',
  `expires_at` datetime DEFAULT NULL COMMENT '过期时间，NULL表示永不过期',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_is_primary` (`is_primary`),
  KEY `idx_expires_at` (`expires_at`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 创建默认系统用户的角色分配（用于测试和初始化）
-- 注意：这里假设users表已存在，实际使用时需要先创建用户记录
INSERT INTO `user_roles` (`user_id`, `role_id`, `is_primary`, `assigned_by`, `created_at`, `updated_at`) VALUES
(1, 1, 1, 1, NOW(), NOW()); -- 给用户ID为1的用户分配超级管理员角色作为主要角色 