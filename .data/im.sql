/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 50726
Source Host           : localhost:3306
Source Database       : blog_temp

Target Server Type    : MYSQL
Target Server Version : 50726
File Encoding         : 65001

Date: 2020-12-28 23:55:56
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `im_friend`
-- ----------------------------
DROP TABLE IF EXISTS `im_friend`;
CREATE TABLE `im_friend` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` bigint(20) DEFAULT '0' COMMENT '我的用户ID',
  `friend_id` bigint(20) DEFAULT '0' COMMENT '朋友的用户ID',
  `add_time` int(10) DEFAULT '0' COMMENT '加为好友时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COMMENT='好友关系表';

-- ----------------------------
-- Records of im_friend
-- ----------------------------

-- ----------------------------
-- Table structure for `im_message`
-- ----------------------------
DROP TABLE IF EXISTS `im_message`;
CREATE TABLE `im_message` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `content_type` tinyint(4) DEFAULT '1' COMMENT '消息类型 1文本 2动画表情',
  `content` text COMMENT '消息内容',
  `create_time` int(10) DEFAULT '0' COMMENT '创建时间',
  `user_id` bigint(20) DEFAULT '0' COMMENT '用户ID',
  `friend_id` bigint(20) DEFAULT '0' COMMENT '朋友ID',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COMMENT='消息内容';

-- ----------------------------
-- Records of im_message
-- ----------------------------

-- ----------------------------
-- Table structure for `im_user`
-- ----------------------------
DROP TABLE IF EXISTS `im_user`;
CREATE TABLE `im_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(200) DEFAULT NULL COMMENT '头像',
  `create_time` int(10) DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) DEFAULT '0' COMMENT '更新时间',
  `username` varchar(50) DEFAULT NULL COMMENT '用户名',
  `password` varchar(200) DEFAULT NULL COMMENT '密码',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- ----------------------------
-- Records of im_user
-- ----------------------------
INSERT INTO `im_user` VALUES ('1', 'A', '/web/static/uploads/images/timg.jpg', '0', '0', 'www', '$2a$10$2dY7VOSi9XmG9cWBEJ.jJ.elaDnCtks4ZDWh9xZ5FUiL1NBnlPR1m');
INSERT INTO `im_user` VALUES ('2', 'W', '/web/static/uploads/images/timg.jpg', '1608368144', '0', 'w', '$2a$10$2dY7VOSi9XmG9cWBEJ.jJ.elaDnCtks4ZDWh9xZ5FUiL1NBnlPR1m');
INSERT INTO `im_user` VALUES ('3', 'F', '/web/static/uploads/images/timg0.jpg', '0', '0', 'haha', '$2a$10$2dY7VOSi9XmG9cWBEJ.jJ.elaDnCtks4ZDWh9xZ5FUiL1NBnlPR1m');
