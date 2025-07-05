-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Waktu pembuatan: 19 Jun 2025 pada 12.58
-- Versi server: 8.0.30
-- Versi PHP: 8.3.20

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `smart_garden`
--

-- --------------------------------------------------------

--
-- Struktur dari tabel `devices`
--

CREATE TABLE `devices` (
  `id` int NOT NULL,
  `device_name` longtext NOT NULL,
  `device_type` varchar(50) DEFAULT 'irrigation',
  `location` longtext,
  `ip_address` longtext,
  `mac_address` varchar(17) DEFAULT NULL,
  `firmware_version` varchar(20) DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `last_seen` timestamp NULL DEFAULT NULL,
  `user_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `devices`
--

INSERT INTO `devices` (`id`, `device_name`, `device_type`, `location`, `ip_address`, `mac_address`, `firmware_version`, `is_active`, `last_seen`, `user_id`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'Smart Irrigation Device 1', 'irrigation', 'Garden Rumah', '192.168.39.89', NULL, NULL, 1, NULL, 1, '2025-06-15 12:29:45', '2025-06-15 12:29:45', NULL),
(4, 'Smart Garden Iot', 'irrigation', 'Depan rumah', '192.168.39.89', NULL, NULL, 1, NULL, 4, '2025-06-18 13:15:54', '2025-06-19 12:38:34', NULL);

-- --------------------------------------------------------

--
-- Struktur dari tabel `sensor_readings`
--

CREATE TABLE `sensor_readings` (
  `id` bigint NOT NULL,
  `device_id` int NOT NULL,
  `temperature` decimal(5,2) DEFAULT NULL COMMENT 'Celsius',
  `humidity` decimal(5,2) DEFAULT NULL COMMENT 'Percentage',
  `temperature_source` enum('sensor','cached') DEFAULT 'sensor',
  `humidity_source` enum('sensor','cached') DEFAULT 'sensor',
  `soil_moisture_raw` int DEFAULT NULL COMMENT 'ADC raw value 0-4095',
  `soil_moisture_percent` decimal(5,2) DEFAULT NULL COMMENT 'Percentage 0-100',
  `water_level_cm` decimal(6,2) DEFAULT NULL COMMENT 'Water height in cm',
  `water_percentage` decimal(5,2) DEFAULT NULL COMMENT 'Tank fill percentage',
  `tank_height_cm` decimal(6,2) DEFAULT '100.00' COMMENT 'Total tank height',
  `pump_status` enum('OFF','MED','HIGH','MAX','NO_WATER') NOT NULL,
  `pump_pwm_value` smallint DEFAULT '0' COMMENT 'PWM value 0-255',
  `pump_percentage` tinyint DEFAULT '0' COMMENT 'Pump power 0-100%',
  `system_status` varchar(50) NOT NULL,
  `logic_explanation` text,
  `wifi_rssi` smallint DEFAULT NULL COMMENT 'WiFi signal strength',
  `free_heap` int DEFAULT NULL COMMENT 'Available memory in bytes',
  `uptime_ms` bigint DEFAULT NULL COMMENT 'Device uptime in milliseconds',
  `device_timestamp` timestamp NULL DEFAULT NULL COMMENT 'Timestamp from device',
  `server_timestamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `sensor_readings`
--

INSERT INTO `sensor_readings` (`id`, `device_id`, `temperature`, `humidity`, `temperature_source`, `humidity_source`, `soil_moisture_raw`, `soil_moisture_percent`, `water_level_cm`, `water_percentage`, `tank_height_cm`, `pump_status`, `pump_pwm_value`, `pump_percentage`, `system_status`, `logic_explanation`, `wifi_rssi`, `free_heap`, `uptime_ms`, `device_timestamp`, `server_timestamp`) VALUES
(1, 4, 29.00, 82.00, 'sensor', 'sensor', 2533, 38.00, 2.00, 13.34, 15.00, 'HIGH', 203, 79, 'PUMPING', 'Mode=FUZZY, Internal=HIGH_FUZZY, Soil=38.0%, Temp=29.0°C, Confidence=0.28', -62, 0, 30015, NULL, '2025-06-19 12:52:51'),
(2, 4, 29.00, 82.00, 'sensor', 'sensor', 2532, 38.00, 2.02, 13.45, 15.00, 'HIGH', 203, 79, 'PUMPING', 'Mode=FUZZY, Internal=HIGH_FUZZY, Soil=38.0%, Temp=29.0°C, Confidence=0.28', -66, 0, 60043, NULL, '2025-06-19 12:53:21'),
(3, 4, 29.00, 82.00, 'sensor', 'sensor', 2537, 38.00, 2.02, 13.45, 15.00, 'HIGH', 203, 79, 'PUMPING', 'Mode=FUZZY, Internal=HIGH_FUZZY, Soil=38.0%, Temp=29.0°C, Confidence=0.28', -61, 0, 90047, NULL, '2025-06-19 12:53:51'),
(4, 4, 29.00, 82.00, 'sensor', 'sensor', 2527, 38.00, 2.02, 13.45, 15.00, 'HIGH', 203, 79, 'PUMPING', 'Mode=FUZZY, Internal=HIGH_FUZZY, Soil=38.0%, Temp=29.0°C, Confidence=0.28', -64, 0, 120103, NULL, '2025-06-19 12:54:21'),
(5, 4, 29.00, 82.00, 'sensor', 'sensor', 2537, 38.00, 2.00, 13.34, 15.00, 'HIGH', 203, 79, 'PUMPING', 'Mode=FUZZY, Internal=HIGH_FUZZY, Soil=38.0%, Temp=29.0°C, Confidence=0.28', -61, 0, 150172, NULL, '2025-06-19 12:54:51'),
(6, 4, 28.90, 82.00, 'sensor', 'sensor', 2527, 38.00, 2.02, 13.45, 15.00, 'HIGH', 203, 79, 'PUMPING', 'Mode=FUZZY, Internal=HIGH_FUZZY, Soil=38.0%, Temp=28.9°C, Confidence=0.29', -63, 0, 180215, NULL, '2025-06-19 12:55:21'),
(7, 4, 28.90, 82.00, 'sensor', 'sensor', 599, 85.00, 2.00, 13.34, 15.00, 'OFF', 0, 0, 'MONITORING', 'Mode=FUZZY, Internal=OFF_FUZZY, Soil=85.0%, Temp=28.9°C, Confidence=0.26', -62, 0, 210221, NULL, '2025-06-19 12:55:51'),
(8, 4, 28.90, 82.00, 'sensor', 'sensor', 511, 87.00, 2.00, 13.34, 15.00, 'OFF', 0, 0, 'MONITORING', 'Mode=FUZZY, Internal=OFF_FUZZY, Soil=87.0%, Temp=28.9°C, Confidence=0.26', -64, 0, 240278, NULL, '2025-06-19 12:56:21');

-- --------------------------------------------------------

--
-- Struktur dari tabel `users`
--

CREATE TABLE `users` (
  `id` int NOT NULL,
  `username` varchar(191) NOT NULL,
  `email` varchar(191) NOT NULL,
  `password` longtext NOT NULL,
  `role` enum('admin','user') DEFAULT 'user',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `users`
--

INSERT INTO `users` (`id`, `username`, `email`, `password`, `role`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'hahaha', 'haha@gmail.com', '$2a$10$eE4/xIBMXaLS2H93cYvjM.c96.XrRpkBEntq5nJDpqZXZegSH8gq.', 'user', '2025-06-15 17:50:57.708', '2025-06-15 17:50:57.708', NULL),
(3, 'adit', 'adit@gmail.com', '$2a$10$SGuVlaHsu6vlKVqsKAie6.pBlEgvnRuRLyfmkYl1Gxw1soXtvbZqq', 'user', '2025-06-15 17:58:25.378', '2025-06-15 17:58:25.378', NULL),
(4, 'admin', 'admin@gmail.com', '$2a$10$xXJa8kFcqgRWDdCxTz7Qbe02QM0FWlCAqbxXa1.n1Om8KJIvUMNjS', 'admin', '2025-06-18 20:10:50.904', '2025-06-18 20:10:50.904', NULL),
(5, 'agus', 'agus@gmail.com', '$2a$10$dHU93fXVtEs0KIj91CLF/uAFkeh9au9UELCOJQrYLKJXtetGhY9Ee', 'user', '2025-06-18 21:04:43.361', '2025-06-18 21:04:43.361', NULL),
(6, 'imam', 'imam@gmail.com', '$2a$10$02wpSchDHuFa31vtrDMBweSC.gVCLLnCsYv/Jd5TXrzY/wF2W//OC', 'user', '2025-06-18 21:29:29.192', '2025-06-18 21:29:29.192', NULL),
(7, 'dodo', 'dodo@gmail.com', '$2a$10$HQ6RpJxLppLGhVl.bcpHfO99S75WhaFUTWDdivhHd7HATIiJUNRz6', 'user', '2025-06-18 22:24:15.823', '2025-06-18 22:24:15.823', NULL),
(9, 'ryan', 'ryan@gmail.com', '$2a$10$KUNxm04QWa1FZGw.mLeHFuwk/4wuudf4aJp5pvyyTuEe1oKGQfYyK', 'user', '2025-06-19 09:21:07.857', '2025-06-19 09:21:07.857', NULL),
(10, 'andi', 'andi@yaho.com', '$2a$10$KXQKD3FAz/Xk91QWKmkOfe4WpJoilPZTnH4gey67l5H9ItOJkcr4W', 'user', '2025-06-19 10:53:28.172', '2025-06-19 10:53:28.172', NULL);

--
-- Indexes for dumped tables
--

--
-- Indeks untuk tabel `devices`
--
ALTER TABLE `devices`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_user_device` (`user_id`,`is_active`),
  ADD KEY `idx_last_seen` (`last_seen`);

--
-- Indeks untuk tabel `sensor_readings`
--
ALTER TABLE `sensor_readings`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_device_time` (`device_id`,`server_timestamp` DESC),
  ADD KEY `idx_device_latest` (`device_id`,`id` DESC),
  ADD KEY `idx_pump_status` (`pump_status`,`server_timestamp`),
  ADD KEY `idx_system_status` (`system_status`,`server_timestamp`);

--
-- Indeks untuk tabel `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uni_users_username` (`username`),
  ADD UNIQUE KEY `uni_users_email` (`email`),
  ADD KEY `idx_users_deleted_at` (`deleted_at`);

--
-- AUTO_INCREMENT untuk tabel yang dibuang
--

--
-- AUTO_INCREMENT untuk tabel `devices`
--
ALTER TABLE `devices`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;

--
-- AUTO_INCREMENT untuk tabel `sensor_readings`
--
ALTER TABLE `sensor_readings`
  MODIFY `id` bigint NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT untuk tabel `users`
--
ALTER TABLE `users`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;

--
-- Ketidakleluasaan untuk tabel pelimpahan (Dumped Tables)
--

--
-- Ketidakleluasaan untuk tabel `devices`
--
ALTER TABLE `devices`
  ADD CONSTRAINT `fk_devices_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Ketidakleluasaan untuk tabel `sensor_readings`
--
ALTER TABLE `sensor_readings`
  ADD CONSTRAINT `sensor_readings_ibfk_1` FOREIGN KEY (`device_id`) REFERENCES `devices` (`id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
