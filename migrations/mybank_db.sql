-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Waktu pembuatan: 07 Jun 2026 pada 06.03
-- Versi server: 10.4.32-MariaDB
-- Versi PHP: 8.5.5

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `mybank_db`
--

-- --------------------------------------------------------

--
-- Struktur dari tabel `accounts`
--

CREATE TABLE `accounts` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `account_number` varchar(30) NOT NULL,
  `account_type` varchar(20) NOT NULL,
  `balance` decimal(18,2) NOT NULL DEFAULT 0.00,
  `currency` varchar(10) NOT NULL DEFAULT 'IDR',
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `accounts`
--

INSERT INTO `accounts` (`id`, `user_id`, `account_number`, `account_type`, `balance`, `currency`, `is_active`, `created_at`, `updated_at`) VALUES
(1, 7, '1000348200', 'Saving', 60000.00, 'IDR', 1, '2026-05-31 13:34:03', '2026-06-07 03:59:18'),
(4, 10, '1009919200', 'Saving', 0.00, 'IDR', 1, '2026-05-31 19:16:54', '2026-05-31 19:16:54');

-- --------------------------------------------------------

--
-- Struktur dari tabel `audit_logs`
--

CREATE TABLE `audit_logs` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `action` varchar(100) NOT NULL,
  `detail` text DEFAULT NULL,
  `ip_address` varchar(45) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Stand-in struktur untuk tampilan `category_spending`
-- (Lihat di bawah untuk tampilan aktual)
--
CREATE TABLE `category_spending` (
`merchant_category` varchar(50)
,`total_amount` decimal(40,2)
,`transaction_count` bigint(21)
);

-- --------------------------------------------------------

--
-- Struktur dari tabel `feature_clicks`
--

CREATE TABLE `feature_clicks` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `feature_name` varchar(100) NOT NULL,
  `click_count` bigint(20) DEFAULT 0,
  `last_clicked` datetime NOT NULL DEFAULT current_timestamp(),
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `recommendations`
--

CREATE TABLE `recommendations` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `type` enum('PROMO','FEATURE','PRODUCT') NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `image_url` varchar(500) DEFAULT NULL,
  `reason` text NOT NULL,
  `priority` int(11) DEFAULT 1,
  `is_active` tinyint(1) DEFAULT 1,
  `expires_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `recommendation_clicks`
--

CREATE TABLE `recommendation_clicks` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `recommendation_id` bigint(20) UNSIGNED NOT NULL,
  `clicked_at` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `transactions`
--

CREATE TABLE `transactions` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `account_id` bigint(20) UNSIGNED NOT NULL,
  `reference_number` varchar(50) NOT NULL,
  `type` enum('TRANSFER','PAYMENT','TOPUP','WITHDRAW','DEPOSIT','QRIS') NOT NULL,
  `status` enum('PENDING','SUCCESS','FAILED','CANCELLED') NOT NULL DEFAULT 'PENDING',
  `amount` decimal(18,2) NOT NULL,
  `fee` decimal(18,2) NOT NULL DEFAULT 0.00,
  `balance_before` decimal(18,2) NOT NULL,
  `balance_after` decimal(18,2) NOT NULL,
  `destination_account_number` varchar(30) DEFAULT NULL,
  `destination_bank_code` varchar(20) DEFAULT NULL,
  `destination_name` varchar(100) DEFAULT NULL,
  `merchant_name` varchar(100) DEFAULT NULL,
  `merchant_category` varchar(50) DEFAULT NULL,
  `merchant_location` varchar(100) DEFAULT NULL,
  `description` text DEFAULT NULL,
  `note` text DEFAULT NULL,
  `fail_reason` text DEFAULT NULL,
  `is_recommended` tinyint(1) NOT NULL DEFAULT 0,
  `recommendation_id` bigint(20) UNSIGNED DEFAULT NULL,
  `transacted_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `users`
--

CREATE TABLE `users` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `email` varchar(100) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `password` varchar(255) NOT NULL,
  `gender` varchar(5) DEFAULT NULL,
  `pin` varchar(255) DEFAULT NULL,
  `date_of_birth` date DEFAULT NULL,
  `occupation` varchar(100) DEFAULT NULL,
  `marital_status` varchar(20) DEFAULT NULL,
  `segment` varchar(50) DEFAULT NULL,
  `monthly_income` decimal(15,2) DEFAULT 0.00,
  `monthly_income_range` varchar(30) DEFAULT NULL,
  `ab_group` varchar(20) DEFAULT 'control',
  `is_personalization_enabled` tinyint(1) NOT NULL DEFAULT 1,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `last_login_at` datetime DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `users`
--

INSERT INTO `users` (`id`, `name`, `email`, `phone`, `password`, `gender`, `pin`, `date_of_birth`, `occupation`, `marital_status`, `segment`, `monthly_income`, `monthly_income_range`, `ab_group`, `is_personalization_enabled`, `is_active`, `last_login_at`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, NULL, 'satria1@gmail.com', NULL, '$2a$12$c4N.zXWkVi1xt7VhFtIJkuDwAnoHl9YNOPk75d0CrNQLG7h6aYPkG', NULL, NULL, NULL, NULL, NULL, NULL, 0.00, NULL, 'control', 1, 1, NULL, '2026-05-06 16:15:53', '2026-05-06 16:15:53', NULL),
(2, NULL, 'satriates1@gmail.com', NULL, '$2a$12$A5CYqePNBhQbLXjfx9DsWOQmWRAL7OioCq3XjACn14gbhEnm2noGG', NULL, NULL, NULL, NULL, NULL, NULL, 0.00, NULL, 'control', 1, 1, NULL, '2026-05-31 12:26:36', '2026-05-31 12:26:36', NULL),
(7, 'Satriates2', 'satriates2@gmail.com', '081234567890', '$2a$12$r1FyHbBhjWME1GdgIJFUgu9YsKfYK..v07Fz3qe.TYP9PdCB3PnrO', NULL, '$2a$12$ziUI77dYTjFXQv1CrVEj9eqvuXV2WDUeVpB8PwQ.ibVXZX078ve4C', '1999-08-17', 'Software Engineer', NULL, NULL, 0.00, NULL, 'control', 1, 1, NULL, '2026-05-31 13:34:03', '2026-06-07 03:27:40', NULL),
(10, 'Satriates3', 'satriates3@gmail.com', '081234567891', '$2a$12$KvxcWUdNDkVG5Hi8.Si7O.YbNQqMel2uUAvpJw6XhpLv4DjKwrBlm', NULL, NULL, '1999-08-18', 'PNS', NULL, NULL, 0.00, NULL, 'control', 1, 1, NULL, '2026-05-31 19:16:54', '2026-05-31 19:16:54', NULL);

-- --------------------------------------------------------

--
-- Struktur untuk view `category_spending`
--
DROP TABLE IF EXISTS `category_spending`;

CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY DEFINER VIEW `category_spending`  AS SELECT `transactions`.`merchant_category` AS `merchant_category`, sum(`transactions`.`amount`) AS `total_amount`, count(0) AS `transaction_count` FROM `transactions` WHERE `transactions`.`status` = 'SUCCESS' GROUP BY `transactions`.`merchant_category` ;

--
-- Indexes for dumped tables
--

--
-- Indeks untuk tabel `accounts`
--
ALTER TABLE `accounts`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `account_number` (`account_number`),
  ADD KEY `fk_accounts_user` (`user_id`);

--
-- Indeks untuk tabel `audit_logs`
--
ALTER TABLE `audit_logs`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_audit_logs_user_id` (`user_id`),
  ADD KEY `idx_audit_logs_action` (`action`),
  ADD KEY `idx_audit_logs_created_at` (`created_at`);

--
-- Indeks untuk tabel `feature_clicks`
--
ALTER TABLE `feature_clicks`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_feature_clicks_user_feature` (`user_id`,`feature_name`),
  ADD KEY `idx_feature_clicks_user_id` (`user_id`);

--
-- Indeks untuk tabel `recommendations`
--
ALTER TABLE `recommendations`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_recommendations_user_id` (`user_id`),
  ADD KEY `idx_recommendations_active` (`user_id`,`is_active`);

--
-- Indeks untuk tabel `recommendation_clicks`
--
ALTER TABLE `recommendation_clicks`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_rec_clicks_user_id` (`user_id`),
  ADD KEY `idx_rec_clicks_rec_id` (`recommendation_id`);

--
-- Indeks untuk tabel `transactions`
--
ALTER TABLE `transactions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `reference_number` (`reference_number`),
  ADD KEY `idx_transactions_user_id` (`user_id`),
  ADD KEY `idx_transactions_account_id` (`account_id`),
  ADD KEY `idx_transactions_status` (`status`),
  ADD KEY `idx_transactions_type` (`type`),
  ADD KEY `idx_transactions_category` (`merchant_category`);

--
-- Indeks untuk tabel `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD UNIQUE KEY `phone` (`phone`),
  ADD KEY `idx_users_email` (`email`),
  ADD KEY `idx_users_phone` (`phone`);

--
-- AUTO_INCREMENT untuk tabel yang dibuang
--

--
-- AUTO_INCREMENT untuk tabel `accounts`
--
ALTER TABLE `accounts`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT untuk tabel `audit_logs`
--
ALTER TABLE `audit_logs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `feature_clicks`
--
ALTER TABLE `feature_clicks`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `recommendations`
--
ALTER TABLE `recommendations`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `recommendation_clicks`
--
ALTER TABLE `recommendation_clicks`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `transactions`
--
ALTER TABLE `transactions`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `users`
--
ALTER TABLE `users`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;

--
-- Ketidakleluasaan untuk tabel pelimpahan (Dumped Tables)
--

--
-- Ketidakleluasaan untuk tabel `accounts`
--
ALTER TABLE `accounts`
  ADD CONSTRAINT `fk_accounts_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Ketidakleluasaan untuk tabel `transactions`
--
ALTER TABLE `transactions`
  ADD CONSTRAINT `fk_transactions_account` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_transactions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
