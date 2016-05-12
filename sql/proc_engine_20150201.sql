/*
Navicat MySQL Data Transfer

Source Server         : mysql_local
Source Server Version : 50626
Source Host           : localhost:3306
Source Database       : proc_engine

Target Server Type    : MYSQL
Target Server Version : 50626
File Encoding         : 65001

Date: 2016-02-01 14:11:01
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for process_definition
-- ----------------------------
DROP TABLE IF EXISTS `process_definition`;
CREATE TABLE `process_definition` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `proc_id` varchar(255) NOT NULL,
  `proc_name` varchar(255) NOT NULL,
  `proc_file` varchar(255) NOT NULL,
  `proc_def` blob,
  PRIMARY KEY (`id`),
  KEY `proc_id` (`proc_id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for rlt_user_proc
-- ----------------------------
DROP TABLE IF EXISTS `rlt_user_proc`;
CREATE TABLE `rlt_user_proc` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `proc_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_proc_id` (`proc_id`),
  KEY `fk_user_id` (`user_id`),
  CONSTRAINT `fk_proc_id` FOREIGN KEY (`proc_id`) REFERENCES `process_definition` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `user_info` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_info
-- ----------------------------
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
