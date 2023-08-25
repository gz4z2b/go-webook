create database webook;

CREATE TABLE `t_user_profile` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Id',
  `user_id` int unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `nickname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '昵称',
  `birthday` bigint unsigned NOT NULL DEFAULT '0' COMMENT '生日',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '个人简介',
  `createtime` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatetime` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deletetime` bigint unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_userid` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户详情';

CREATE TABLE `t_user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `email` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `createtime` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatetime` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deletetime` bigint unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户';