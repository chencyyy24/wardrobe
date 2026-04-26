-- ============================================================
-- 个人衣柜 (Personal Wardrobe) 数据库结构
-- 目标数据库: MySQL 8.0+
-- 字符集: utf8mb4 (支持中文和 emoji)
-- ============================================================

CREATE DATABASE IF NOT EXISTS sql_wardrobe
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE sql_wardrobe;

-- ============================================================
-- 1. 用户表 (预留多用户扩展)
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
  id         INT          PRIMARY KEY AUTO_INCREMENT,
  username   VARCHAR(100) NOT NULL UNIQUE,
  password   VARCHAR(255) NOT NULL COMMENT 'bcrypt 哈希',
  nickname   VARCHAR(100) DEFAULT NULL,
  avatar     VARCHAR(500) DEFAULT NULL,
  created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 2. 大类表 (外套/上衣/裤子/裙子/鞋子/配饰)
-- ============================================================
CREATE TABLE IF NOT EXISTS categories (
  id         INT         PRIMARY KEY AUTO_INCREMENT,
  name       VARCHAR(50) NOT NULL UNIQUE COMMENT '大类名称',
  sort_order INT         NOT NULL DEFAULT 0 COMMENT '排序权重',
  created_at DATETIME   NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 3. 小类表 (如 T恤、衬衫、牛仔裤……)
-- ============================================================
CREATE TABLE IF NOT EXISTS subcategories (
  id          INT         PRIMARY KEY AUTO_INCREMENT,
  category_id INT         NOT NULL COMMENT '所属大类 ID',
  name        VARCHAR(50) NOT NULL COMMENT '小类名称',
  sort_order  INT         NOT NULL DEFAULT 0 COMMENT '排序权重',
  created_at  DATETIME   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
  UNIQUE KEY uk_category_subcategory (category_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 4. 衣物表
-- ============================================================
CREATE TABLE IF NOT EXISTS clothing_items (
  id              CHAR(36)     PRIMARY KEY COMMENT 'UUID',
  user_id         INT          NOT NULL,
  name            VARCHAR(255) NOT NULL COMMENT '衣物名称（用户自定义）',
  category_id     INT          NOT NULL COMMENT '大类 ID',
  subcategory_id  INT          DEFAULT NULL COMMENT '小类 ID',
  original_image  VARCHAR(500) NOT NULL COMMENT '原图路径',
  masked_image    VARCHAR(500) DEFAULT NULL COMMENT '白底抠图路径',
  status          ENUM('pending', 'done', 'failed') NOT NULL DEFAULT 'pending' COMMENT '抠图状态',
  price           DECIMAL(10, 2) DEFAULT NULL COMMENT '购买价格（可选）',
  created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at      DATETIME     DEFAULT NULL COMMENT '软删除时间',
  FOREIGN KEY (user_id)        REFERENCES users(id),
  FOREIGN KEY (category_id)    REFERENCES categories(id),
  FOREIGN KEY (subcategory_id) REFERENCES subcategories(id) ON DELETE SET NULL,
  INDEX idx_user_category (user_id, category_id),
  INDEX idx_status (status),
  INDEX idx_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 5. 搭配表
-- ============================================================
CREATE TABLE IF NOT EXISTS outfits (
  id         CHAR(36)     PRIMARY KEY COMMENT 'UUID',
  user_id    INT          NOT NULL,
  name       VARCHAR(255) DEFAULT NULL COMMENT '搭配名称（自动生成或用户自定义）',
  card_image VARCHAR(500) DEFAULT NULL COMMENT '搭配卡片图路径',
  created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME     DEFAULT NULL COMMENT '软删除时间',
  FOREIGN KEY (user_id) REFERENCES users(id),
  INDEX idx_user_created (user_id, created_at DESC),
  INDEX idx_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 6. 搭配-衣物关联表 (slot 表示穿搭部位)
-- slot 取值: outer | top | bottom | skirt | shoes | accessory
--    外套     上衣   裤子    裙子    鞋子    配饰
-- ============================================================
CREATE TABLE IF NOT EXISTS outfit_items (
  id          INT          PRIMARY KEY AUTO_INCREMENT,
  outfit_id   CHAR(36)     NOT NULL COMMENT '搭配 ID',
  clothing_id CHAR(36)     NOT NULL COMMENT '衣物 ID',
  slot        ENUM('outer', 'top', 'bottom', 'skirt', 'shoes', 'accessory') NOT NULL COMMENT '穿搭部位',
  FOREIGN KEY (outfit_id)   REFERENCES outfits(id) ON DELETE CASCADE,
  FOREIGN KEY (clothing_id) REFERENCES clothing_items(id),
  UNIQUE KEY uk_outfit_slot (outfit_id, slot),
  INDEX idx_outfit (outfit_id),
  INDEX idx_clothing (clothing_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 7. 穿着记录表 (可选，用于统计)
-- ============================================================
CREATE TABLE IF NOT EXISTS wear_records (
  id          INT      PRIMARY KEY AUTO_INCREMENT,
  user_id     INT      NOT NULL,
  clothing_id CHAR(36) DEFAULT NULL COMMENT '可为 NULL（按搭配记录时不指定单品）',
  outfit_id   CHAR(36) DEFAULT NULL COMMENT '可为 NULL（单穿不属任何搭配）',
  wear_date   DATE     NOT NULL COMMENT '穿着日期',
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id)      REFERENCES users(id),
  FOREIGN KEY (clothing_id)  REFERENCES clothing_items(id) ON DELETE SET NULL,
  FOREIGN KEY (outfit_id)    REFERENCES outfits(id) ON DELETE SET NULL,
  INDEX idx_wear_date (wear_date),
  INDEX idx_user_date (user_id, wear_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 初始化数据: 大类 & 小类
-- ============================================================
INSERT INTO categories (name, sort_order) VALUES
  ('外套', 1),
  ('上衣', 2),
  ('裤子', 3),
  ('裙子', 4),
  ('鞋子', 5),
  ('配饰', 6);

INSERT INTO subcategories (category_id, name, sort_order) VALUES
  -- 外套
  (1, '风衣',       1),
  (1, '牛仔夹克',   2),
  (1, '西装',       3),
  (1, '针织开衫',   4),
  -- 上衣
  (2, 'T恤',        1),
  (2, '衬衫',       2),
  (2, '卫衣',       3),
  (2, '毛衣',       4),
  -- 裤子
  (3, '牛仔裤',     1),
  (3, '休闲裤',     2),
  (3, '短裤',       3),
  -- 裙子
  (4, 'JK裙',       1),
  (4, '百褶裙',     2),
  (4, '连衣裙',     3),
  -- 鞋子
  (5, '运动鞋',     1),
  (5, '帆布鞋',     2),
  (5, '靴子',       3),
  -- 配饰
  (6, '帽子',       1),
  (6, '项链',       2),
  (6, '耳环',       3),
  (6, '包包',       4);
