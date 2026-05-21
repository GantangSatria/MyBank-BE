-- ============================================================
-- MyBank-BE Database Migrations
-- Jalankan script ini di MySQL database mybank_db
-- ============================================================

-- 1. Tambah kolom baru di tabel users (jika tabel sudah ada)
-- ============================================================
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS monthly_income DECIMAL(15,2) DEFAULT 0 AFTER segment,
  ADD COLUMN IF NOT EXISTS ab_group VARCHAR(20) DEFAULT 'control' AFTER monthly_income;


-- 2. Buat tabel transactions
-- ============================================================
CREATE TABLE IF NOT EXISTS transactions (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  account_id BIGINT UNSIGNED NOT NULL,
  reference_number VARCHAR(100) NOT NULL UNIQUE,
  type ENUM('TRANSFER', 'PAYMENT', 'TOPUP', 'WITHDRAW', 'DEPOSIT', 'QRIS') NOT NULL,
  status ENUM('PENDING', 'SUCCESS', 'FAILED', 'CANCELLED') NOT NULL DEFAULT 'PENDING',
  amount DECIMAL(15,2) NOT NULL,
  fee DECIMAL(15,2) DEFAULT 0,
  balance_before DECIMAL(15,2) DEFAULT 0,
  balance_after DECIMAL(15,2) DEFAULT 0,

  destination_account_number VARCHAR(50),
  destination_bank_code VARCHAR(20),
  destination_name VARCHAR(200),

  merchant_name VARCHAR(200),
  merchant_category VARCHAR(100),
  merchant_location VARCHAR(200),

  description VARCHAR(255),
  note VARCHAR(255),
  fail_reason VARCHAR(255),

  is_recommended BOOLEAN DEFAULT FALSE,
  recommendation_id BIGINT UNSIGNED,

  transacted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_transactions_user_id (user_id),
  INDEX idx_transactions_status (status),
  INDEX idx_transactions_type (type),
  INDEX idx_transactions_transacted_at (transacted_at),
  INDEX idx_transactions_merchant_category (merchant_category),
  INDEX idx_transactions_user_date (user_id, transacted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 3. Buat tabel accounts (jika belum ada)
-- ============================================================
CREATE TABLE IF NOT EXISTS accounts (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  account_number VARCHAR(50) NOT NULL UNIQUE,
  account_type VARCHAR(50) NOT NULL DEFAULT 'SAVINGS',
  balance DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'IDR',
  is_active BOOLEAN DEFAULT TRUE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_accounts_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 4. Buat tabel feature_clicks
-- ============================================================
CREATE TABLE IF NOT EXISTS feature_clicks (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  feature_name VARCHAR(100) NOT NULL,
  click_count BIGINT DEFAULT 0,
  last_clicked DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  UNIQUE KEY uk_feature_clicks_user_feature (user_id, feature_name),
  INDEX idx_feature_clicks_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 5. Buat tabel recommendations
-- ============================================================
CREATE TABLE IF NOT EXISTS recommendations (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  type ENUM('PROMO', 'FEATURE', 'PRODUCT') NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  image_url VARCHAR(500),
  reason TEXT NOT NULL,
  priority INT DEFAULT 1,
  is_active BOOLEAN DEFAULT TRUE,
  expires_at DATETIME,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_recommendations_user_id (user_id),
  INDEX idx_recommendations_active (user_id, is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 6. Buat tabel recommendation_clicks
-- ============================================================
CREATE TABLE IF NOT EXISTS recommendation_clicks (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  recommendation_id BIGINT UNSIGNED NOT NULL,
  clicked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  INDEX idx_rec_clicks_user_id (user_id),
  INDEX idx_rec_clicks_rec_id (recommendation_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 7. Buat tabel audit_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  action VARCHAR(100) NOT NULL,
  detail TEXT,
  ip_address VARCHAR(45),
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  INDEX idx_audit_logs_user_id (user_id),
  INDEX idx_audit_logs_action (action),
  INDEX idx_audit_logs_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
