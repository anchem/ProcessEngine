/*
Navicat MySQL Data Transfer

Source Server         : mysql_local
Source Server Version : 50621
Source Host           : localhost:3306
Source Database       : proc_engine

Target Server Type    : MYSQL
Target Server Version : 50621
File Encoding         : 65001

Date: 2016-03-01 10:49:20
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for process_definition
-- ----------------------------
DROP TABLE IF EXISTS `process_definition`;
CREATE TABLE `process_definition` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '流程在系统中的唯一标识',
  `proc_id` varchar(255) NOT NULL COMMENT '流程定义中的ID',
  `proc_name` varchar(255) NOT NULL COMMENT '流程名',
  `proc_desc` varchar(1024) DEFAULT NULL COMMENT '流程描述',
  `proc_file` varchar(512) NOT NULL COMMENT '流程文件名',
  `proc_def` blob COMMENT '序列化后的流程',
  `create_time` int(11) NOT NULL COMMENT '创建时间',
  `multiple_instance` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否多实例{ 0：单实例（默认）；1：多实例 }',
  `is_excutable` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否可执行{ 0：不可执行；1：可执行（默认） }',
  PRIMARY KEY (`id`),
  KEY `proc_id` (`proc_id`)
) ENGINE=InnoDB AUTO_INCREMENT=80 DEFAULT CHARSET=utf8;

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
  CONSTRAINT `fk_proc_id` FOREIGN KEY (`proc_id`) REFERENCES `process_definition` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=72 DEFAULT CHARSET=utf8;
