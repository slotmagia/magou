-- 创建角色表
CREATE TABLE IF NOT EXISTS `roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` varchar(200) DEFAULT NULL COMMENT '角色描述',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1=启用 0=禁用',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序号，数字越小越靠前',
  `data_scope` tinyint(4) NOT NULL DEFAULT '1' COMMENT '数据权限范围：1=全部数据 2=部门数据 3=部门及以下数据 4=仅本人数据',
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

-- 插入默认角色数据
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`, `sort`, `data_scope`, `remark`, `created_by`, `updated_by`, `created_at`, `updated_at`) VALUES
(1, '超级管理员', 'super_admin', '拥有系统所有权限的超级管理员角色', 1, 1, 1, '系统内置角色，不可删除', 1, 1, NOW(), NOW()),
(2, '系统管理员', 'system_admin', '系统管理员，拥有系统管理相关权限', 1, 2, 2, '负责系统基础配置和用户管理', 1, 1, NOW(), NOW()),

(4, '财务管理员', 'finance_admin', '财务管理员，拥有财务统计和报表权限', 1, 4, 3, '负责财务数据查看和导出', 1, 1, NOW(), NOW()),
(5, '运营人员', 'operator', '日常运营人员，拥有基础操作权限', 1, 5, 4, '负责日常业务操作', 1, 1, NOW(), NOW()),
(6, '客服人员', 'customer_service', '客服人员，拥有订单查看和处理权限', 1, 6, 4, '负责客户服务和订单处理', 1, 1, NOW(), NOW()),
(7, '审计人员', 'auditor', '审计人员，拥有只读权限', 1, 7, 1, '负责系统审计和监督', 1, 1, NOW(), NOW()); 