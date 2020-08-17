-- phpMyAdmin SQL Dump
-- version 4.9.4
-- https://www.phpmyadmin.net/
--
-- 主机： 127.0.0.1
-- 生成日期： 2020-08-18 01:42:40
-- 服务器版本： 8.0.21
-- PHP 版本： 7.4.3

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";

--
-- 数据库： `vm_manager`
--
CREATE DATABASE IF NOT EXISTS `vm_manager` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `vm_manager`;

-- --------------------------------------------------------

--
-- 表的结构 `vm_disks`
--

CREATE TABLE `vm_disks` (
                            `id` bigint UNSIGNED NOT NULL,
                            `vm_id` bigint UNSIGNED NOT NULL,
                            `disk_path` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '磁盘地址'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- 表的结构 `vm_lists`
--

CREATE TABLE `vm_lists` (
                            `id` bigint UNSIGNED NOT NULL,
                            `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名字',
                            `cpu` int NOT NULL COMMENT 'cpu核心数量',
                            `mem` int NOT NULL COMMENT '内存大小，单位m',
                            `auto_startup` int NOT NULL DEFAULT '0' COMMENT '是否自动启动 0-false 1-true',
                            `status` int NOT NULL DEFAULT '0' COMMENT '当前状态 0-未启动 1-运行中 2-错误',
                            `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- 表的结构 `vm_macs`
--

CREATE TABLE `vm_macs` (
                           `id` bigint UNSIGNED NOT NULL,
                           `vm_id` bigint UNSIGNED NOT NULL COMMENT '虚拟机id',
                           `mac` char(17) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'mac地址'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- 表的结构 `vm_pcis`
--

CREATE TABLE `vm_pcis` (
                           `id` bigint UNSIGNED NOT NULL,
                           `vm_id` bigint UNSIGNED NOT NULL,
                           `pci_id` varchar(32) COLLATE utf8mb4_general_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- 表的结构 `vm_ports`
--

CREATE TABLE `vm_ports` (
                            `id` bigint UNSIGNED NOT NULL,
                            `type` int NOT NULL COMMENT '类型 0-spice 1-vnc 2-monitor',
                            `vm_id` bigint UNSIGNED NOT NULL COMMENT '虚拟机id',
                            `port` int UNSIGNED NOT NULL COMMENT '端口'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- 转储表的索引
--

--
-- 表的索引 `vm_disks`
--
ALTER TABLE `vm_disks`
    ADD PRIMARY KEY (`id`),
    ADD KEY `fk_disk_vm_id` (`vm_id`);

--
-- 表的索引 `vm_lists`
--
ALTER TABLE `vm_lists`
    ADD PRIMARY KEY (`id`);

--
-- 表的索引 `vm_macs`
--
ALTER TABLE `vm_macs`
    ADD PRIMARY KEY (`id`),
    ADD UNIQUE KEY `mac` (`mac`),
    ADD KEY `fk_mac_vm_id` (`vm_id`);

--
-- 表的索引 `vm_pcis`
--
ALTER TABLE `vm_pcis`
    ADD PRIMARY KEY (`id`),
    ADD KEY `fk_pci_vm_id` (`vm_id`);

--
-- 表的索引 `vm_ports`
--
ALTER TABLE `vm_ports`
    ADD PRIMARY KEY (`id`),
    ADD UNIQUE KEY `port` (`port`),
    ADD KEY `fk_port_vm_id` (`vm_id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `vm_disks`
--
ALTER TABLE `vm_disks`
    MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `vm_lists`
--
ALTER TABLE `vm_lists`
    MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `vm_macs`
--
ALTER TABLE `vm_macs`
    MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `vm_pcis`
--
ALTER TABLE `vm_pcis`
    MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `vm_ports`
--
ALTER TABLE `vm_ports`
    MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 限制导出的表
--

--
-- 限制表 `vm_disks`
--
ALTER TABLE `vm_disks`
    ADD CONSTRAINT `fk_disk_vm_id` FOREIGN KEY (`vm_id`) REFERENCES `vm_lists` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

--
-- 限制表 `vm_macs`
--
ALTER TABLE `vm_macs`
    ADD CONSTRAINT `fk_mac_vm_id` FOREIGN KEY (`vm_id`) REFERENCES `vm_lists` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

--
-- 限制表 `vm_pcis`
--
ALTER TABLE `vm_pcis`
    ADD CONSTRAINT `fk_pci_vm_id` FOREIGN KEY (`vm_id`) REFERENCES `vm_lists` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

--
-- 限制表 `vm_ports`
--
ALTER TABLE `vm_ports`
    ADD CONSTRAINT `fk_port_vm_id` FOREIGN KEY (`vm_id`) REFERENCES `vm_lists` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
COMMIT;
