-- MySQL dump 10.16  Distrib 10.2.15-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: monitor
-- ------------------------------------------------------
-- Server version	10.2.15-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `active_probe`
--

DROP TABLE IF EXISTS `active_probe`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `active_probe` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `plugin_name` varchar(30) NOT NULL,
  `host_name` varchar(50) NOT NULL,
  `host_ip` varchar(15) NOT NULL,
  `interval` int(11) NOT NULL DEFAULT 0,
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `UNI_ActiveProbe_plugin_name_host_name` (`plugin_name`,`host_name`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `active_probe_config`
--

DROP TABLE IF EXISTS `active_probe_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `active_probe_config` (
  `id` int(255) NOT NULL AUTO_INCREMENT,
  `active_probe_id` int(11) NOT NULL,
  `target` varchar(255) NOT NULL,
  `arg1` varchar(255) NOT NULL DEFAULT '',
  `arg2` varchar(255) NOT NULL DEFAULT '',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQUE_ActiveProbeConfig_id_target` (`active_probe_id`,`target`) USING BTREE,
  KEY `IDX_ActiveProbeConfig_id` (`id`) USING BTREE,
  CONSTRAINT `active_probe_config_ibfk_1` FOREIGN KEY (`active_probe_id`) REFERENCES `active_probe` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `alarm_judge`
--

DROP TABLE IF EXISTS `alarm_judge`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `alarm_judge` (
  `alarm_name` varchar(50) NOT NULL COMMENT '需同influxdb的measurement',
  `alarmele` varchar(50) NOT NULL COMMENT 'influxdb的typeinstance+type',
  `ajtype` enum('le','ne','ge') NOT NULL COMMENT '判断类型，<=、!=、>=',
  `level1` int(11) DEFAULT NULL,
  `level2` int(11) DEFAULT NULL,
  `level3` int(11) DEFAULT NULL,
  PRIMARY KEY (`alarm_name`,`alarmele`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `alarm_link`
--

DROP TABLE IF EXISTS `alarm_link`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `alarm_link` (
  `alarm_name` varchar(50) NOT NULL,
  `list` text NOT NULL DEFAULT '',
  `type` enum('team','staff','other') NOT NULL DEFAULT 'other',
  `channel` int(3) NOT NULL DEFAULT 0 COMMENT '采用二进制表示\r\n0仅更新状态;\r\nB001 邮件;\r\nB010 企业微信;\r\n如同时发邮件企业微信填3',
  PRIMARY KEY (`alarm_name`) USING BTREE,
  KEY `channel` (`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `alarm_queue`
--

DROP TABLE IF EXISTS `alarm_queue`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `alarm_queue` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `host_name` varchar(15) NOT NULL,
  `alarm_name` varchar(100) NOT NULL,
  `alarmele` varchar(100) NOT NULL,
  `value` double(16,2) NOT NULL,
  `message` varchar(255) NOT NULL DEFAULT '',
  `handle_man` varchar(25) NOT NULL DEFAULT '',
  `stat` int(1) NOT NULL DEFAULT 0 COMMENT '0初始,1已发告警，2已处理中,3关闭',
  `level` enum('level1','level2','level3') NOT NULL DEFAULT 'level1',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `IDX_AlarmQueue_stat` (`stat`)
) ENGINE=InnoDB AUTO_INCREMENT=65243 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `plugin`
--

DROP TABLE IF EXISTS `plugin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `plugin` (
  `plugin_name` varchar(32) NOT NULL,
  `plugin_type` varchar(16) NOT NULL DEFAULT '' COMMENT 'python;shell..',
  `file_name` varchar(128) NOT NULL,
  `comment` varchar(256) NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`plugin_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `plugin_config`
--

DROP TABLE IF EXISTS `plugin_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `plugin_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `host_ip` varchar(255) NOT NULL,
  `host_name` varchar(255) NOT NULL,
  `plugin_name` varchar(255) NOT NULL,
  `interval` int(11) NOT NULL DEFAULT 0,
  `timeout` int(11) NOT NULL DEFAULT 3,
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQUE_PluginConfig_host_name_plugin_name` (`plugin_name`,`host_name`) USING BTREE,
  CONSTRAINT `plugin_config_ibfk_1` FOREIGN KEY (`plugin_name`) REFERENCES `plugin` (`plugin_name`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-07-02 18:24:14
