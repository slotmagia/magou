-- 租户表
CREATE TABLE IF NOT EXISTS `tenants` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '租户ID',
  `name` varchar(100) NOT NULL COMMENT '租户名称',
  `code` varchar(50) NOT NULL COMMENT '租户编码',
  `domain` varchar(100) DEFAULT NULL COMMENT '租户域名',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态：1=正常 2=锁定 3=禁用',
  `max_users` int(11) DEFAULT '100' COMMENT '最大用户数',
  `storage_limit` bigint(20) DEFAULT '10737418240' COMMENT '存储限制(字节)',
  `expire_at` datetime DEFAULT NULL COMMENT '过期时间',
  `admin_user_id` bigint(20) DEFAULT NULL COMMENT '租户管理员用户ID',
  `config` json DEFAULT NULL COMMENT '租户配置',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `created_by` bigint(20) DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '更新者',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_name` (`name`),
  KEY `idx_domain` (`domain`),
  KEY `idx_status` (`status`),
  KEY `idx_expire_at` (`expire_at`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- 插入默认租户（系统租户）
INSERT INTO `tenants` (`id`, `name`, `code`, `status`, `max_users`, `admin_user_id`, `remark`, `created_by`) VALUES
(1, '系统租户', 'system', 1, 999999, 1, '系统默认租户，不可删除', 0);
