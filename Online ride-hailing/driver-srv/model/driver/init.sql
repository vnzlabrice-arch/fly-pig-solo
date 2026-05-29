SET NAMES utf8mb4;

-- 司机用户表
CREATE TABLE IF NOT EXISTS `driver_users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `phone` varchar(11) NOT NULL,
  `nickname` varchar(32) DEFAULT '',
  `real_name` varchar(32) DEFAULT '',
  `avatar_url` varchar(255) DEFAULT '',
  `gender` tinyint DEFAULT 0,
  `status` tinyint DEFAULT 1,
  `service_score` decimal(3,2) DEFAULT 5.00,
  `order_count` bigint DEFAULT 0,
  `last_login_time` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机用户';

-- 司机资料详情
CREATE TABLE IF NOT EXISTS `driver_profiles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `driver_id` bigint NOT NULL,
  `license_no` varchar(32) DEFAULT '',
  `license_expire` date DEFAULT NULL,
  `drive_years` int DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_driver_id` (`driver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机资料';

-- 车辆信息
CREATE TABLE IF NOT EXISTS `driver_vehicles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `driver_id` bigint NOT NULL,
  `plate_no` varchar(16) NOT NULL,
  `brand` varchar(32) DEFAULT '',
  `model` varchar(32) DEFAULT '',
  `color` varchar(16) DEFAULT '',
  `status` tinyint DEFAULT 1,
  PRIMARY KEY (`id`),
  KEY `idx_driver_id` (`driver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='车辆信息';

-- 司机认证审核
CREATE TABLE IF NOT EXISTS `driver_auths` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `driver_id` bigint NOT NULL,
  `auth_type` tinyint NOT NULL,
  `pic_url` varchar(255) DEFAULT '',
  `audit_status` tinyint DEFAULT 0,
  `audit_time` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_driver_id` (`driver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utfmb4 COMMENT='司机认证';

-- 司机在线状态 & 位置
CREATE TABLE IF NOT EXISTS `driver_online_status` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `driver_id` bigint NOT NULL,
  `online_status` tinyint DEFAULT 0,
  `accept_order` tinyint DEFAULT 0,
  `lng` decimal(10,6) DEFAULT 0.000000,
  `lat` decimal(10,6) DEFAULT 0.000000,
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_driver_id` (`driver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机在线状态';

-- 司机订单推荐（待接单）
CREATE TABLE IF NOT EXISTS `driver_order_recommends` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `driver_id` bigint NOT NULL,
  `order_id` varchar(64) NOT NULL,
  `distance` int DEFAULT 0,
  `estimate_price` decimal(10,2) DEFAULT 0.00,
  `status` tinyint DEFAULT 1,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_driver_id` (`driver_id`),
  KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机订单推荐';

-- 司机收入明细
CREATE TABLE IF NOT EXISTS `driver_incomes` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `driver_id` bigint NOT NULL,
  `order_id` varchar(64) NOT NULL,
  `total_fee` decimal(10,2) DEFAULT 0.00,
  `platform_fee` decimal(10,2) DEFAULT 0.00,
  `actual_income` decimal(10,2) DEFAULT 0.00,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_driver_id` (`driver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机收入';

-- 司机钱包
CREATE TABLE IF NOT EXISTS `driver_wallets` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `driver_id` bigint NOT NULL,
  `balance` decimal(12,2) DEFAULT 0.00,
  `withdrawable` decimal(12,2) DEFAULT 0.00,
  `frozen` decimal(12,2) DEFAULT 0.00,
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_driver_id` (`driver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机钱包';

-- 司机提现记录
CREATE TABLE IF NOT EXISTS `driver_withdraws` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `driver_id` bigint NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `bank_card` varchar(32) DEFAULT '',
  `status` tinyint DEFAULT 1,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `arrive_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_driver_id` (`driver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机提现';

-- 司机违规记录
CREATE TABLE IF NOT EXISTS `driver_violations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `driver_id` bigint NOT NULL,
  `order_id` varchar(64) DEFAULT '',
  `violation_type` tinyint NOT NULL,
  `score_deduct` int DEFAULT 0,
  `reason` varchar(255) DEFAULT '',
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_driver_id` (`driver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机违规';

-- 司机服务区域
CREATE TABLE IF NOT EXISTS `driver_service_areas` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `driver_id` bigint NOT NULL,
  `city_code` varchar(16) NOT NULL,
  `status` tinyint DEFAULT 1,
  PRIMARY KEY (`id`),
  KEY `idx_driver_id` (`driver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机服务区域';

-- 司机消息
CREATE TABLE IF NOT EXISTS `driver_messages` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `driver_id` bigint NOT NULL,
  `title` varchar(64) DEFAULT '',
  `content` varchar(512) DEFAULT '',
  `is_read` tinyint DEFAULT 0,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_driver_id` (`driver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机消息';
