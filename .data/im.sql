/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 50726
Source Host           : localhost:3306
Source Database       : im

Target Server Type    : MYSQL
Target Server Version : 50726
File Encoding         : 65001

Date: 2020-12-29 19:15:44
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
) ENGINE=MyISAM AUTO_INCREMENT=102 DEFAULT CHARSET=utf8mb4 COMMENT='好友关系表';

-- ----------------------------
-- Records of im_friend
-- ----------------------------
INSERT INTO `im_friend` VALUES ('1', '1', '2', '1608395478');
INSERT INTO `im_friend` VALUES ('2', '1', '3', '1608395478');
INSERT INTO `im_friend` VALUES ('3', '1', '4', '1608395478');
INSERT INTO `im_friend` VALUES ('4', '1', '5', '1608395478');
INSERT INTO `im_friend` VALUES ('5', '1', '6', '1608395478');
INSERT INTO `im_friend` VALUES ('6', '1', '7', '1608395478');
INSERT INTO `im_friend` VALUES ('7', '1', '8', '1608395478');
INSERT INTO `im_friend` VALUES ('8', '1', '9', '1608395478');
INSERT INTO `im_friend` VALUES ('9', '1', '10', '1608395478');
INSERT INTO `im_friend` VALUES ('10', '1', '11', '1608395478');
INSERT INTO `im_friend` VALUES ('11', '1', '12', '1608395478');
INSERT INTO `im_friend` VALUES ('101', '2', '4', '1608395499');

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
) ENGINE=MyISAM AUTO_INCREMENT=480 DEFAULT CHARSET=utf8mb4 COMMENT='消息内容';

-- ----------------------------
-- Records of im_message
-- ----------------------------
INSERT INTO `im_message` VALUES ('1', '1', '你好哈！854', '1608449028', '1', '2');
INSERT INTO `im_message` VALUES ('3', '1', '你好哈！625', '1608449028', '1', '4');
INSERT INTO `im_message` VALUES ('4', '1', '你好哈！103', '1608449028', '1', '5');
INSERT INTO `im_message` VALUES ('5', '1', '你好哈！596', '1608449028', '1', '6');
INSERT INTO `im_message` VALUES ('6', '1', '你好哈！596', '1608449028', '1', '7');
INSERT INTO `im_message` VALUES ('7', '1', '你好哈！815', '1608449028', '1', '8');
INSERT INTO `im_message` VALUES ('8', '1', '你好哈！815', '1608449028', '1', '9');
INSERT INTO `im_message` VALUES ('9', '1', '你好哈！498', '1608449028', '1', '10');
INSERT INTO `im_message` VALUES ('10', '1', '你好哈！722', '1608449028', '1', '11');
INSERT INTO `im_message` VALUES ('11', '1', '你好哈！420', '1608449067', '1', '2');
INSERT INTO `im_message` VALUES ('32', '1', '我真的来了', '1608459066', '1', '2');
INSERT INTO `im_message` VALUES ('13', '1', '你好哈！852', '1608449067', '1', '4');
INSERT INTO `im_message` VALUES ('14', '1', '你好哈！289', '1608449067', '1', '5');
INSERT INTO `im_message` VALUES ('15', '1', '你好哈！289', '1608449067', '1', '6');
INSERT INTO `im_message` VALUES ('16', '1', '你好哈！483', '1608449067', '1', '7');
INSERT INTO `im_message` VALUES ('17', '1', '你好哈！483', '1608449067', '1', '8');
INSERT INTO `im_message` VALUES ('18', '1', '你好哈！778', '1608449067', '1', '9');
INSERT INTO `im_message` VALUES ('19', '1', '你好哈！171', '1608449067', '1', '10');
INSERT INTO `im_message` VALUES ('20', '1', '你好哈！843', '1608449067', '1', '11');
INSERT INTO `im_message` VALUES ('21', '1', '你好哈！744', '1608449113', '2', '1');
INSERT INTO `im_message` VALUES ('31', '1', '我来了', '1608459032', '1', '2');
INSERT INTO `im_message` VALUES ('23', '1', '你好哈！967', '1608449113', '4', '1');
INSERT INTO `im_message` VALUES ('24', '1', '你好哈！480', '1608449113', '5', '1');
INSERT INTO `im_message` VALUES ('25', '1', '你好哈！400', '1608449113', '6', '1');
INSERT INTO `im_message` VALUES ('26', '1', '你好哈！400', '1608449113', '7', '1');
INSERT INTO `im_message` VALUES ('27', '1', '你好哈！667', '1608449113', '8', '1');
INSERT INTO `im_message` VALUES ('28', '1', '你好哈！667', '1608449113', '9', '1');
INSERT INTO `im_message` VALUES ('29', '1', '你好哈！1301577173864', '1577173864', '12', '1');
INSERT INTO `im_message` VALUES ('30', '1', '你好哈！1301608796264', '1608796264', '11', '1');
INSERT INTO `im_message` VALUES ('33', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('34', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('35', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('36', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('37', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('38', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('39', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('40', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('41', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('42', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('43', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('44', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('45', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('46', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('47', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('48', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('49', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('50', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('51', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('52', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('53', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('54', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('55', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('56', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('57', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('58', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('59', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('60', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('61', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('62', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('63', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('64', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('65', '1', '是的', '1609089054', '4', '1');
INSERT INTO `im_message` VALUES ('66', '1', '哈哈哈', '1609089278', '4', '1');
INSERT INTO `im_message` VALUES ('67', '1', '来客', '1609089280', '4', '1');
INSERT INTO `im_message` VALUES ('68', '1', '浪费计算机', '1609089281', '4', '1');
INSERT INTO `im_message` VALUES ('69', '1', '分解落实', '1609089282', '4', '1');
INSERT INTO `im_message` VALUES ('70', '1', '杰弗里斯', '1609089284', '4', '1');
INSERT INTO `im_message` VALUES ('71', '1', '分解落实', '1609089285', '4', '1');
INSERT INTO `im_message` VALUES ('72', '1', '分解落实', '1609089285', '4', '1');
INSERT INTO `im_message` VALUES ('73', '1', '嘻嘻', '1609089512', '4', '1');
INSERT INTO `im_message` VALUES ('74', '1', '是啊', '1609089517', '4', '1');
INSERT INTO `im_message` VALUES ('75', '1', '房价低', '1609089518', '4', '1');
INSERT INTO `im_message` VALUES ('76', '1', '大姐夫is', '1609089520', '4', '1');
INSERT INTO `im_message` VALUES ('77', '1', '大家粉丝', '1609089521', '4', '1');
INSERT INTO `im_message` VALUES ('78', '1', '福建地税', '1609089522', '4', '1');
INSERT INTO `im_message` VALUES ('79', '1', '打飞机死哦', '1609089523', '4', '1');
INSERT INTO `im_message` VALUES ('80', '1', '福佳斯', '1609089524', '4', '1');
INSERT INTO `im_message` VALUES ('81', '1', '飞机师', '1609089525', '4', '1');
INSERT INTO `im_message` VALUES ('82', '1', '飞机师', '1609089526', '4', '1');
INSERT INTO `im_message` VALUES ('83', '1', '福佳斯', '1609089527', '4', '1');
INSERT INTO `im_message` VALUES ('84', '1', '福建省第', '1609089529', '4', '1');
INSERT INTO `im_message` VALUES ('85', '1', '福佳斯', '1609089530', '4', '1');
INSERT INTO `im_message` VALUES ('86', '1', '福佳斯', '1609089531', '4', '1');
INSERT INTO `im_message` VALUES ('87', '1', '房价低', '1609089532', '4', '1');
INSERT INTO `im_message` VALUES ('88', '1', '发动机是', '1609089533', '4', '1');
INSERT INTO `im_message` VALUES ('89', '1', '房价低', '1609089534', '4', '1');
INSERT INTO `im_message` VALUES ('90', '1', '福佳斯', '1609089536', '4', '1');
INSERT INTO `im_message` VALUES ('91', '1', '发动机为', '1609089537', '4', '1');
INSERT INTO `im_message` VALUES ('92', '1', '福建地税', '1609089538', '4', '1');
INSERT INTO `im_message` VALUES ('93', '1', '飞机师', '1609089539', '4', '1');
INSERT INTO `im_message` VALUES ('94', '1', '福佳斯', '1609089540', '4', '1');
INSERT INTO `im_message` VALUES ('95', '1', '福佳斯', '1609089541', '4', '1');
INSERT INTO `im_message` VALUES ('96', '1', '福佳斯', '1609089542', '4', '1');
INSERT INTO `im_message` VALUES ('97', '1', '飞机师', '1609089543', '4', '1');
INSERT INTO `im_message` VALUES ('98', '1', '飞机是', '1609089545', '4', '1');
INSERT INTO `im_message` VALUES ('99', '1', '飞机', '1609089546', '4', '1');
INSERT INTO `im_message` VALUES ('100', '1', '飞机', '1609089547', '4', '1');
INSERT INTO `im_message` VALUES ('101', '1', '管控', '1609089549', '4', '1');
INSERT INTO `im_message` VALUES ('102', '1', '费劲口舌', '1609089549', '4', '1');
INSERT INTO `im_message` VALUES ('103', '1', '估计', '1609089550', '4', '1');
INSERT INTO `im_message` VALUES ('104', '1', '飞机', '1609089551', '4', '1');
INSERT INTO `im_message` VALUES ('105', '1', '送排', '1609089553', '4', '1');
INSERT INTO `im_message` VALUES ('106', '1', '飞机师', '1609089554', '4', '1');
INSERT INTO `im_message` VALUES ('107', '1', '飞机', '1609089556', '4', '1');
INSERT INTO `im_message` VALUES ('108', '1', 'fd', '1609089616', '1', '4');
INSERT INTO `im_message` VALUES ('109', '1', '好的', '1609089628', '1', '4');
INSERT INTO `im_message` VALUES ('110', '1', '历史', '1609089631', '1', '4');
INSERT INTO `im_message` VALUES ('111', '1', '废旧塑料', '1609089633', '1', '4');
INSERT INTO `im_message` VALUES ('112', '1', '来了', '1609089635', '1', '4');
INSERT INTO `im_message` VALUES ('113', '1', '飞机师', '1609089637', '1', '4');
INSERT INTO `im_message` VALUES ('114', '1', '列回复', '1609089639', '1', '4');
INSERT INTO `im_message` VALUES ('115', '1', '附件四欧锦', '1609089641', '1', '4');
INSERT INTO `im_message` VALUES ('116', '1', '类似回复', '1609089642', '1', '4');
INSERT INTO `im_message` VALUES ('117', '1', '废旧塑料', '1609089644', '1', '4');
INSERT INTO `im_message` VALUES ('118', '1', '房间里思', '1609089645', '1', '4');
INSERT INTO `im_message` VALUES ('119', '1', '废旧塑料', '1609089647', '1', '4');
INSERT INTO `im_message` VALUES ('120', '1', '老司机费', '1609089648', '1', '4');
INSERT INTO `im_message` VALUES ('121', '1', '飞机搜了搜垃圾分类', '1609089651', '1', '4');
INSERT INTO `im_message` VALUES ('122', '1', '解放路上季度', '1609089653', '1', '4');
INSERT INTO `im_message` VALUES ('123', '1', '房间里思了', '1609089655', '1', '4');
INSERT INTO `im_message` VALUES ('124', '1', '解放了时间', '1609089656', '1', '4');
INSERT INTO `im_message` VALUES ('125', '1', '放假了时间了', '1609089658', '1', '4');
INSERT INTO `im_message` VALUES ('126', '1', '杰弗里斯屁股', '1609089660', '1', '4');
INSERT INTO `im_message` VALUES ('127', '1', '飞机类似的减肥呢', '1609089661', '1', '4');
INSERT INTO `im_message` VALUES ('128', '1', '飞机螺丝刀', '1609089663', '1', '4');
INSERT INTO `im_message` VALUES ('129', '1', '飞机刺痒元', '1609089668', '1', '4');
INSERT INTO `im_message` VALUES ('130', '1', '咖喱鸡', '1609090285', '2', '1');
INSERT INTO `im_message` VALUES ('131', '1', 'look', '1609090308', '2', '1');
INSERT INTO `im_message` VALUES ('132', '1', '空间', '1609090317', '2', '1');
INSERT INTO `im_message` VALUES ('133', '1', '太空旅客', '1609090321', '2', '1');
INSERT INTO `im_message` VALUES ('134', '1', '班尼路', '1609090327', '2', '1');
INSERT INTO `im_message` VALUES ('135', '1', '低发热', '1609090736', '1', '2');
INSERT INTO `im_message` VALUES ('136', '1', '1\n', '1609090755', '1', '2');
INSERT INTO `im_message` VALUES ('137', '1', '2', '1609090761', '1', '2');
INSERT INTO `im_message` VALUES ('138', '1', '是的', '1609090805', '2', '1');
INSERT INTO `im_message` VALUES ('139', '1', '还好', '1609090810', '2', '1');
INSERT INTO `im_message` VALUES ('140', '1', '来了', '1609090815', '2', '1');
INSERT INTO `im_message` VALUES ('141', '1', '⊙ω⊙', '1609090858', '2', '1');
INSERT INTO `im_message` VALUES ('142', '1', '⊙ω⊙⊙ω⊙', '1609090864', '2', '1');
INSERT INTO `im_message` VALUES ('143', '1', '[微笑]', '1609090944', '1', '2');
INSERT INTO `im_message` VALUES ('144', '1', '身边', '1609090966', '2', '1');
INSERT INTO `im_message` VALUES ('145', '1', '来了', '1609090970', '2', '1');
INSERT INTO `im_message` VALUES ('146', '1', '我饿啦', '1609090974', '2', '1');
INSERT INTO `im_message` VALUES ('147', '1', '都卡', '1609090978', '2', '1');
INSERT INTO `im_message` VALUES ('148', '1', '是的', '1609090983', '2', '1');
INSERT INTO `im_message` VALUES ('149', '1', '原来在哦了', '1609090991', '2', '1');
INSERT INTO `im_message` VALUES ('150', '1', '是哦', '1609090996', '2', '1');
INSERT INTO `im_message` VALUES ('151', '1', '啊哈', '1609091009', '2', '1');
INSERT INTO `im_message` VALUES ('152', '1', '来了', '1609091012', '2', '1');
INSERT INTO `im_message` VALUES ('153', '1', '中国', '1609091016', '2', '1');
INSERT INTO `im_message` VALUES ('154', '1', '是的', '1609091079', '1', '4');
INSERT INTO `im_message` VALUES ('155', '1', '舒服', '1609091083', '1', '4');
INSERT INTO `im_message` VALUES ('156', '1', '回复度过', '1609091096', '1', '4');
INSERT INTO `im_message` VALUES ('157', '1', '我觉得', '1609091106', '2', '1');
INSERT INTO `im_message` VALUES ('158', '1', 'vnjc', '1609091112', '1', '2');
INSERT INTO `im_message` VALUES ('159', '1', '凤凰都会', '1609091115', '1', '2');
INSERT INTO `im_message` VALUES ('160', '1', '胡海峰', '1609091118', '1', '2');
INSERT INTO `im_message` VALUES ('161', '1', '和辅导班', '1609091120', '1', '2');
INSERT INTO `im_message` VALUES ('162', '1', '宽带缴费', '1609091121', '1', '2');
INSERT INTO `im_message` VALUES ('163', '1', '回复电话', '1609091123', '1', '2');
INSERT INTO `im_message` VALUES ('164', '1', '客服电话', '1609091125', '1', '2');
INSERT INTO `im_message` VALUES ('165', '1', ' 发货的bhfdu', '1609091129', '1', '2');
INSERT INTO `im_message` VALUES ('166', '1', '飞虎队发布', '1609091131', '1', '2');
INSERT INTO `im_message` VALUES ('167', '1', '还复读更方便', '1609091133', '1', '2');
INSERT INTO `im_message` VALUES ('168', '1', '付大哥高鸿股份', '1609091136', '1', '2');
INSERT INTO `im_message` VALUES ('169', '1', '互粉都会', '1609091138', '1', '2');
INSERT INTO `im_message` VALUES ('170', '1', '和风度和', '1609091140', '1', '2');
INSERT INTO `im_message` VALUES ('171', '1', '库高鸿股份', '1609091142', '1', '2');
INSERT INTO `im_message` VALUES ('172', '1', '符合度还不对黄飞鸿', '1609091145', '1', '2');
INSERT INTO `im_message` VALUES ('173', '1', '和辅导班', '1609091147', '1', '2');
INSERT INTO `im_message` VALUES ('174', '1', '还复读好好发挥', '1609091149', '1', '2');
INSERT INTO `im_message` VALUES ('175', '1', '符合度很高的股份', '1609091151', '1', '2');
INSERT INTO `im_message` VALUES ('176', '1', '符合度发广告', '1609091153', '1', '2');
INSERT INTO `im_message` VALUES ('177', '1', '符合度非官方个', '1609091155', '1', '2');
INSERT INTO `im_message` VALUES ('178', '1', '飞虎队苟富贵', '1609091157', '1', '2');
INSERT INTO `im_message` VALUES ('179', '1', '福个个都是', '1609091159', '1', '2');
INSERT INTO `im_message` VALUES ('180', '1', '烦得很付大哥', '1609091161', '1', '2');
INSERT INTO `im_message` VALUES ('181', '1', '符合度过', '1609091162', '1', '2');
INSERT INTO `im_message` VALUES ('182', '1', 'logo', '1609091189', '2', '1');
INSERT INTO `im_message` VALUES ('183', '1', '我卡了', '1609091197', '2', '1');
INSERT INTO `im_message` VALUES ('184', '1', '恩爱经理', '1609091202', '2', '1');
INSERT INTO `im_message` VALUES ('185', '1', 'Emily', '1609091205', '2', '1');
INSERT INTO `im_message` VALUES ('186', '1', '我来了', '1609091208', '2', '1');
INSERT INTO `im_message` VALUES ('187', '1', '肆无忌惮', '1609091213', '2', '1');
INSERT INTO `im_message` VALUES ('188', '1', '好的呢', '1609091217', '2', '1');
INSERT INTO `im_message` VALUES ('189', '1', '呵呵呵', '1609091220', '2', '1');
INSERT INTO `im_message` VALUES ('190', '1', '敏婷民莫咯ing莫咯后哦咯JOJO你OKlook哦嘞哦咯了龙哦嘞你虐了龙哦了咯空间了龙哦了红楼梦龙口哦咯咯哦LOL木有哦JOJO空间哦咯了LOL哦咯龙走咯够了某民乐啊咯哈多米杜宾辛苦哦咯巨虫咯哈DJ你那头LOL你呢咯心累你妹哦预订龙女抽我够了OK你为你OK民疼哦了轰隆隆', '1609091260', '2', '1');
INSERT INTO `im_message` VALUES ('191', '1', '是的', '1609091689', '1', '2');
INSERT INTO `im_message` VALUES ('192', '1', '啦', '1609091767', '2', '1');
INSERT INTO `im_message` VALUES ('193', '1', '了啊', '1609092299', '2', '1');
INSERT INTO `im_message` VALUES ('194', '1', '卡口', '1609092390', '2', '1');
INSERT INTO `im_message` VALUES ('195', '1', '用', '1609092595', '2', '1');
INSERT INTO `im_message` VALUES ('196', '1', '哈哈', '1609093105', '1', '2');
INSERT INTO `im_message` VALUES ('197', '1', '就好', '1609093135', '1', '2');
INSERT INTO `im_message` VALUES ('198', '1', 'jfjjf', '1609093147', '2', '1');
INSERT INTO `im_message` VALUES ('199', '1', 'hhhdhf', '1609093154', '2', '1');
INSERT INTO `im_message` VALUES ('200', '1', 'nujfjd', '1609093173', '2', '1');
INSERT INTO `im_message` VALUES ('201', '1', 'hudj', '1609093183', '2', '1');
INSERT INTO `im_message` VALUES ('202', '1', '✋', '1609093198', '2', '1');
INSERT INTO `im_message` VALUES ('203', '1', 'jhdjd', '1609093438', '2', '1');
INSERT INTO `im_message` VALUES ('204', '1', '✊', '1609093753', '2', '1');
INSERT INTO `im_message` VALUES ('205', '1', '✊', '1609093811', '2', '1');
INSERT INTO `im_message` VALUES ('206', '1', 'tf', '1609095084', '2', '1');
INSERT INTO `im_message` VALUES ('207', '1', '哈哈', '1609095379', '2', '1');
INSERT INTO `im_message` VALUES ('208', '1', '[\\u1f600]', '1609095396', '2', '1');
INSERT INTO `im_message` VALUES ('209', '1', '[\\u1f916]', '1609095405', '2', '1');
INSERT INTO `im_message` VALUES ('210', '1', '哈哈哈', '1609095889', '2', '1');
INSERT INTO `im_message` VALUES ('211', '1', '东', '1609095902', '2', '1');
INSERT INTO `im_message` VALUES ('212', '1', '嗯好', '1609095916', '2', '1');
INSERT INTO `im_message` VALUES ('213', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('214', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('215', '1', '啦', '1609096624', '2', '1');
INSERT INTO `im_message` VALUES ('216', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('217', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('218', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('219', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('220', '1', '卡啊', '1609096667', '2', '1');
INSERT INTO `im_message` VALUES ('221', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('222', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('223', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('224', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('225', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('226', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('227', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('228', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('229', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('230', '0', '', '0', '0', '0');
INSERT INTO `im_message` VALUES ('231', '1', '哈喽', '1609097035', '4', '1');
INSERT INTO `im_message` VALUES ('232', '1', '谷歌', '1609097412', '4', '1');
INSERT INTO `im_message` VALUES ('233', '1', '回复第', '1609097434', '4', '1');
INSERT INTO `im_message` VALUES ('234', '1', 'vbfv', '1609097552', '4', '1');
INSERT INTO `im_message` VALUES ('235', '1', '发动机帝豪', '1609097570', '4', '1');
INSERT INTO `im_message` VALUES ('236', '1', '关于', '1609097806', '4', '1');
INSERT INTO `im_message` VALUES ('237', '1', '复活电击', '1609097830', '4', '1');
INSERT INTO `im_message` VALUES ('238', '1', '☠️[\\u1f480]', '1609097850', '4', '1');
INSERT INTO `im_message` VALUES ('239', '1', '[\\u1f600]', '1609097884', '4', '1');
INSERT INTO `im_message` VALUES ('240', '1', '[\\u1f603]', '1609097892', '4', '1');
INSERT INTO `im_message` VALUES ('241', '1', '[\\u1f601]', '1609097897', '4', '1');
INSERT INTO `im_message` VALUES ('242', '1', '[\\u1f600]', '1609097925', '4', '1');
INSERT INTO `im_message` VALUES ('243', '1', '[\\u1f600]', '1609097941', '4', '1');
INSERT INTO `im_message` VALUES ('244', '1', '[\\u1f635]', '1609097946', '4', '1');
INSERT INTO `im_message` VALUES ('245', '1', '[\\u1f48b]', '1609097956', '4', '1');
INSERT INTO `im_message` VALUES ('246', '1', '[\\u1f600]', '1609097960', '4', '1');
INSERT INTO `im_message` VALUES ('247', '1', '[\\u1f440]', '1609097964', '4', '1');
INSERT INTO `im_message` VALUES ('248', '1', '[\\u1f600]您好啊[\\u1f600]', '1609097999', '4', '1');
INSERT INTO `im_message` VALUES ('249', '1', '[\\u1f91c]', '1609098415', '4', '1');
INSERT INTO `im_message` VALUES ('250', '1', '[\\u1f9e0]', '1609098420', '4', '1');
INSERT INTO `im_message` VALUES ('251', '1', '[\\u1f600][\\u1f61a][\\u1f612][\\u1f628][\\u1f9e0]', '1609098506', '4', '1');
INSERT INTO `im_message` VALUES ('252', '1', '胡', '1609099099', '4', '1');
INSERT INTO `im_message` VALUES ('253', '1', '突然', '1609099106', '4', '1');
INSERT INTO `im_message` VALUES ('254', '1', '符合度', '1609099124', '4', '1');
INSERT INTO `im_message` VALUES ('255', '1', '积分ID', '1609099136', '4', '1');
INSERT INTO `im_message` VALUES ('256', '1', '积分ID', '1609099140', '4', '1');
INSERT INTO `im_message` VALUES ('257', '1', '放到', '1609099154', '4', '1');
INSERT INTO `im_message` VALUES ('258', '1', '积分', '1609099372', '4', '1');
INSERT INTO `im_message` VALUES ('259', '1', '见覅', '1609099397', '4', '1');
INSERT INTO `im_message` VALUES ('260', '1', '发过的', '1609099570', '4', '1');
INSERT INTO `im_message` VALUES ('261', '1', '发多少', '1609099572', '4', '1');
INSERT INTO `im_message` VALUES ('262', '1', '发多少', '1609099573', '4', '1');
INSERT INTO `im_message` VALUES ('263', '1', '放到', '1609099573', '4', '1');
INSERT INTO `im_message` VALUES ('264', '1', '发多少', '1609099574', '4', '1');
INSERT INTO `im_message` VALUES ('265', '1', '放到', '1609099575', '4', '1');
INSERT INTO `im_message` VALUES ('266', '1', '放到', '1609099576', '4', '1');
INSERT INTO `im_message` VALUES ('267', '1', '放到', '1609099577', '4', '1');
INSERT INTO `im_message` VALUES ('268', '1', '发送到', '1609099578', '4', '1');
INSERT INTO `im_message` VALUES ('269', '1', '发多少', '1609099579', '4', '1');
INSERT INTO `im_message` VALUES ('270', '1', '方式', '1609099580', '4', '1');
INSERT INTO `im_message` VALUES ('271', '1', '方式', '1609099581', '4', '1');
INSERT INTO `im_message` VALUES ('272', '1', '都是', '1609099582', '4', '1');
INSERT INTO `im_message` VALUES ('273', '1', '方式', '1609099583', '4', '1');
INSERT INTO `im_message` VALUES ('274', '1', '方式', '1609099584', '4', '1');
INSERT INTO `im_message` VALUES ('275', '1', '方式', '1609099584', '4', '1');
INSERT INTO `im_message` VALUES ('276', '1', '就是', '1609099840', '4', '1');
INSERT INTO `im_message` VALUES ('277', '1', '哈哈', '1609099847', '4', '1');
INSERT INTO `im_message` VALUES ('278', '1', '房价低', '1609099861', '4', '1');
INSERT INTO `im_message` VALUES ('279', '1', '54', '1609099976', '4', '1');
INSERT INTO `im_message` VALUES ('280', '1', '565', '1609100334', '4', '1');
INSERT INTO `im_message` VALUES ('281', '1', '56656', '1609100385', '4', '1');
INSERT INTO `im_message` VALUES ('282', '1', '级', '1609100636', '4', '1');
INSERT INTO `im_message` VALUES ('283', '1', '你好', '1609100673', '4', '1');
INSERT INTO `im_message` VALUES ('284', '1', '[\\u1f600]', '1609100809', '4', '1');
INSERT INTO `im_message` VALUES ('285', '1', '中国', '1609100822', '4', '1');
INSERT INTO `im_message` VALUES ('286', '1', '[\\u1f9e0]', '1609100871', '4', '1');
INSERT INTO `im_message` VALUES ('287', '1', '5665', '1609101183', '4', '1');
INSERT INTO `im_message` VALUES ('288', '1', '5', '1609101273', '4', '1');
INSERT INTO `im_message` VALUES ('289', '1', '69', '1609101359', '4', '1');
INSERT INTO `im_message` VALUES ('290', '1', '969', '1609101401', '4', '1');
INSERT INTO `im_message` VALUES ('291', '1', '9889', '1609101592', '4', '1');
INSERT INTO `im_message` VALUES ('292', '1', '6556', '1609101691', '4', '1');
INSERT INTO `im_message` VALUES ('293', '1', '4545', '1609101755', '4', '1');
INSERT INTO `im_message` VALUES ('294', '1', '5656', '1609101766', '4', '1');
INSERT INTO `im_message` VALUES ('295', '1', '5665', '1609101787', '4', '1');
INSERT INTO `im_message` VALUES ('296', '1', '忠告', '1609102325', '4', '1');
INSERT INTO `im_message` VALUES ('297', '1', '555', '1609102377', '4', '1');
INSERT INTO `im_message` VALUES ('298', '1', '233223', '1609102406', '4', '1');
INSERT INTO `im_message` VALUES ('299', '1', '4545', '1609102622', '4', '1');
INSERT INTO `im_message` VALUES ('300', '1', '445', '1609102907', '4', '1');
INSERT INTO `im_message` VALUES ('301', '1', '[\\u1f630]', '1609102977', '4', '1');
INSERT INTO `im_message` VALUES ('302', '1', '21212', '1609103120', '4', '1');
INSERT INTO `im_message` VALUES ('303', '1', '454545', '1609103208', '4', '1');
INSERT INTO `im_message` VALUES ('304', '1', '好大夫', '1609103349', '4', '1');
INSERT INTO `im_message` VALUES ('305', '1', '545', '1609103400', '4', '1');
INSERT INTO `im_message` VALUES ('306', '1', '21212', '1609103810', '4', '1');
INSERT INTO `im_message` VALUES ('307', '1', '5454545', '1609103825', '4', '1');
INSERT INTO `im_message` VALUES ('308', '1', '2112', '1609103868', '4', '1');
INSERT INTO `im_message` VALUES ('309', '1', '65', '1609106802', '4', '1');
INSERT INTO `im_message` VALUES ('310', '1', '6565', '1609107408', '4', '1');
INSERT INTO `im_message` VALUES ('311', '1', '54545', '1609107458', '4', '1');
INSERT INTO `im_message` VALUES ('312', '1', '5', '1609107558', '4', '1');
INSERT INTO `im_message` VALUES ('313', '1', '级', '1609108121', '4', '1');
INSERT INTO `im_message` VALUES ('314', '1', '放到分户式会', '1609108156', '4', '1');
INSERT INTO `im_message` VALUES ('315', '1', '飞机爱搜好看恢复好会', '1609108159', '4', '1');
INSERT INTO `im_message` VALUES ('316', '1', '付好款哈哈法医会', '1609108161', '4', '1');
INSERT INTO `im_message` VALUES ('317', '1', '粉红色的啊哈客户爱U欧尼', '1609108164', '4', '1');
INSERT INTO `im_message` VALUES ('318', '1', '5656', '1609108805', '4', '1');
INSERT INTO `im_message` VALUES ('319', '1', '8787878', '1609108840', '4', '1');
INSERT INTO `im_message` VALUES ('320', '1', '5656656', '1609108985', '4', '1');
INSERT INTO `im_message` VALUES ('321', '1', '你好', '1609109089', '4', '1');
INSERT INTO `im_message` VALUES ('322', '1', '解放路第四按计划', '1609109091', '4', '1');
INSERT INTO `im_message` VALUES ('323', '1', '[\\u1f44a][\\u1f4a9]', '1609109108', '4', '1');
INSERT INTO `im_message` VALUES ('324', '1', '[\\u1f600][\\u1f603]', '1609109114', '4', '1');
INSERT INTO `im_message` VALUES ('325', '1', '发多少', '1609109116', '4', '1');
INSERT INTO `im_message` VALUES ('326', '1', '大家傲娇', '1609109117', '4', '1');
INSERT INTO `im_message` VALUES ('327', '1', '卡夫卡哈哈发货安徽合肥哈哈哈', '1609110334', '4', '1');
INSERT INTO `im_message` VALUES ('328', '1', '超模小睵上合法哦是目露苛上', '1609110338', '4', '1');
INSERT INTO `im_message` VALUES ('329', '1', '腹肌撕裂虎', '1609110340', '4', '1');
INSERT INTO `im_message` VALUES ('330', '1', '里湖上', '1609110956', '4', '1');
INSERT INTO `im_message` VALUES ('331', '1', '飞机螯合剂上', '1609110969', '4', '1');
INSERT INTO `im_message` VALUES ('332', '1', 'fjoiahih', '1609110974', '4', '1');
INSERT INTO `im_message` VALUES ('333', '1', '[\\u1f480][\\u1f9b7]', '1609110980', '4', '1');
INSERT INTO `im_message` VALUES ('334', '1', '1', '1609111485', '4', '1');
INSERT INTO `im_message` VALUES ('335', '1', '2', '1609111487', '4', '1');
INSERT INTO `im_message` VALUES ('336', '1', '3', '1609111489', '4', '1');
INSERT INTO `im_message` VALUES ('337', '1', '4', '1609111490', '4', '1');
INSERT INTO `im_message` VALUES ('338', '1', '5', '1609111492', '4', '1');
INSERT INTO `im_message` VALUES ('339', '1', '6', '1609111495', '4', '1');
INSERT INTO `im_message` VALUES ('340', '1', '7', '1609111496', '4', '1');
INSERT INTO `im_message` VALUES ('341', '1', '8', '1609111498', '4', '1');
INSERT INTO `im_message` VALUES ('342', '1', '9', '1609111500', '4', '1');
INSERT INTO `im_message` VALUES ('343', '1', '10', '1609111502', '4', '1');
INSERT INTO `im_message` VALUES ('344', '1', '96554', '1609111803', '4', '1');
INSERT INTO `im_message` VALUES ('345', '1', '96998854', '1609111807', '4', '1');
INSERT INTO `im_message` VALUES ('346', '1', '8989454844515', '1609111813', '4', '1');
INSERT INTO `im_message` VALUES ('347', '1', '9874525', '1609111817', '4', '1');
INSERT INTO `im_message` VALUES ('348', '1', '585554', '1609111822', '4', '1');
INSERT INTO `im_message` VALUES ('349', '1', '[\\u1f612][\\u1f612]', '1609111829', '4', '1');
INSERT INTO `im_message` VALUES ('350', '1', '1', '1609112568', '4', '1');
INSERT INTO `im_message` VALUES ('351', '1', '2', '1609112572', '4', '1');
INSERT INTO `im_message` VALUES ('352', '1', '3', '1609112574', '4', '1');
INSERT INTO `im_message` VALUES ('353', '1', '4', '1609112578', '4', '1');
INSERT INTO `im_message` VALUES ('354', '1', '5', '1609112580', '4', '1');
INSERT INTO `im_message` VALUES ('355', '1', '6', '1609112581', '4', '1');
INSERT INTO `im_message` VALUES ('356', '1', '7', '1609112584', '4', '1');
INSERT INTO `im_message` VALUES ('357', '1', '8', '1609112585', '4', '1');
INSERT INTO `im_message` VALUES ('358', '1', '9', '1609112586', '4', '1');
INSERT INTO `im_message` VALUES ('359', '1', '10', '1609112588', '4', '1');
INSERT INTO `im_message` VALUES ('360', '1', '11', '1609112589', '4', '1');
INSERT INTO `im_message` VALUES ('361', '1', '12', '1609112590', '4', '1');
INSERT INTO `im_message` VALUES ('362', '1', '13', '1609112591', '4', '1');
INSERT INTO `im_message` VALUES ('363', '1', '14', '1609112592', '4', '1');
INSERT INTO `im_message` VALUES ('364', '1', '15', '1609112593', '4', '1');
INSERT INTO `im_message` VALUES ('365', '1', '16', '1609112595', '4', '1');
INSERT INTO `im_message` VALUES ('366', '1', '17', '1609112596', '4', '1');
INSERT INTO `im_message` VALUES ('367', '1', '18', '1609112597', '4', '1');
INSERT INTO `im_message` VALUES ('368', '1', '19', '1609112599', '4', '1');
INSERT INTO `im_message` VALUES ('369', '1', '20', '1609112601', '4', '1');
INSERT INTO `im_message` VALUES ('370', '1', '21', '1609112603', '4', '1');
INSERT INTO `im_message` VALUES ('371', '1', '22', '1609112604', '4', '1');
INSERT INTO `im_message` VALUES ('372', '1', '23', '1609112606', '4', '1');
INSERT INTO `im_message` VALUES ('373', '1', '24', '1609112607', '4', '1');
INSERT INTO `im_message` VALUES ('374', '1', '25', '1609112608', '4', '1');
INSERT INTO `im_message` VALUES ('375', '1', '1', '1609112681', '4', '1');
INSERT INTO `im_message` VALUES ('376', '1', '2', '1609112683', '4', '1');
INSERT INTO `im_message` VALUES ('377', '1', '3', '1609112684', '4', '1');
INSERT INTO `im_message` VALUES ('378', '1', '4', '1609112685', '4', '1');
INSERT INTO `im_message` VALUES ('379', '1', '5', '1609112686', '4', '1');
INSERT INTO `im_message` VALUES ('380', '1', '6', '1609112687', '4', '1');
INSERT INTO `im_message` VALUES ('381', '1', '7', '1609112688', '4', '1');
INSERT INTO `im_message` VALUES ('382', '1', '8', '1609112689', '4', '1');
INSERT INTO `im_message` VALUES ('383', '1', '9', '1609112690', '4', '1');
INSERT INTO `im_message` VALUES ('384', '1', '10', '1609112692', '4', '1');
INSERT INTO `im_message` VALUES ('385', '1', '11', '1609112693', '4', '1');
INSERT INTO `im_message` VALUES ('386', '1', '12', '1609112694', '4', '1');
INSERT INTO `im_message` VALUES ('387', '1', '13', '1609112695', '4', '1');
INSERT INTO `im_message` VALUES ('388', '1', '14', '1609112696', '4', '1');
INSERT INTO `im_message` VALUES ('389', '1', '15', '1609112698', '4', '1');
INSERT INTO `im_message` VALUES ('390', '1', '16', '1609112699', '4', '1');
INSERT INTO `im_message` VALUES ('391', '1', '17', '1609112700', '4', '1');
INSERT INTO `im_message` VALUES ('392', '1', '18', '1609112701', '4', '1');
INSERT INTO `im_message` VALUES ('393', '1', '19', '1609112703', '4', '1');
INSERT INTO `im_message` VALUES ('394', '1', '20', '1609112704', '4', '1');
INSERT INTO `im_message` VALUES ('395', '1', '回话', '1609113042', '4', '1');
INSERT INTO `im_message` VALUES ('396', '1', 'hin', '1609113045', '4', '1');
INSERT INTO `im_message` VALUES ('397', '1', '非机动', '1609113047', '4', '1');
INSERT INTO `im_message` VALUES ('398', '1', '打击', '1609113049', '4', '1');
INSERT INTO `im_message` VALUES ('399', '1', '咯偶家', '1609113051', '4', '1');
INSERT INTO `im_message` VALUES ('400', '1', '见附件', '1609113053', '4', '1');
INSERT INTO `im_message` VALUES ('401', '1', '飞机', '1609113054', '4', '1');
INSERT INTO `im_message` VALUES ('402', '1', '[\\u1f63e][\\u1f63e]', '1609113091', '4', '1');
INSERT INTO `im_message` VALUES ('403', '1', '522', '1609113341', '4', '1');
INSERT INTO `im_message` VALUES ('404', '1', '66', '1609113360', '4', '1');
INSERT INTO `im_message` VALUES ('405', '1', '77', '1609113365', '4', '1');
INSERT INTO `im_message` VALUES ('406', '1', '88', '1609113369', '4', '1');
INSERT INTO `im_message` VALUES ('407', '1', '100', '1609113558', '4', '1');
INSERT INTO `im_message` VALUES ('408', '1', '101', '1609113560', '4', '1');
INSERT INTO `im_message` VALUES ('409', '1', '102', '1609113562', '4', '1');
INSERT INTO `im_message` VALUES ('410', '1', '[\\u1f9e0][\\u1f9b7]', '1609113627', '1', '4');
INSERT INTO `im_message` VALUES ('411', '1', '[\\u1f9b7][\\u1f9b4][\\u1f440]', '1609116216', '1', '4');
INSERT INTO `im_message` VALUES ('412', '1', '<img src=\"/web/assets/image/emoji/3.jpeg\" alt=\"\">', '1609121840', '1', '4');
INSERT INTO `im_message` VALUES ('413', '1', '感觉覅', '1609122197', '1', '4');
INSERT INTO `im_message` VALUES ('414', '1', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609122200', '1', '4');
INSERT INTO `im_message` VALUES ('415', '1', '&#60;&#105;&#109;&#103;&#32;&#115;&#116;&#121;&#108;&#101;&#61;&#34;&#119;&#105;&#100;&#116;&#104;&#58;&#56;&#48;&#112;&#120;&#59;&#104;&#101;&#105;&#103;&#104;&#116;&#58;&#56;&#48;&#112;&#120;&#34;&#32;&#115;&#114;&#99;&#61;&#34;&#47;&#119;&#101;&#98;&#47;&#97;&#115;&#115;&#101;&#116;&#115;&#47;&#105;&#109;&#97;&#103;&#101;&#47;&#101;&#109;&#111;&#106;&#105;&#47;&#51;&#46;&#106;&#112;&#101;&#103;&#34;&#62;', '1609122702', '1', '4');
INSERT INTO `im_message` VALUES ('416', '1', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609122730', '1', '4');
INSERT INTO `im_message` VALUES ('417', '1', '5656', '1609122996', '4', '1');
INSERT INTO `im_message` VALUES ('418', '1', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609123765', '4', '1');
INSERT INTO `im_message` VALUES ('419', '1', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609123769', '4', '1');
INSERT INTO `im_message` VALUES ('420', '1', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609123772', '4', '1');
INSERT INTO `im_message` VALUES ('421', '1', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609123775', '4', '1');
INSERT INTO `im_message` VALUES ('422', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609124086', '4', '1');
INSERT INTO `im_message` VALUES ('423', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609124110', '4', '1');
INSERT INTO `im_message` VALUES ('424', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609125453', '1', '4');
INSERT INTO `im_message` VALUES ('425', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609125457', '1', '4');
INSERT INTO `im_message` VALUES ('426', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609125535', '1', '4');
INSERT INTO `im_message` VALUES ('427', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609125849', '1', '4');
INSERT INTO `im_message` VALUES ('428', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609125912', '1', '4');
INSERT INTO `im_message` VALUES ('429', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609126323', '1', '4');
INSERT INTO `im_message` VALUES ('430', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609126390', '1', '4');
INSERT INTO `im_message` VALUES ('431', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609126535', '1', '4');
INSERT INTO `im_message` VALUES ('432', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609126548', '1', '4');
INSERT INTO `im_message` VALUES ('433', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609126592', '1', '4');
INSERT INTO `im_message` VALUES ('434', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609126596', '1', '4');
INSERT INTO `im_message` VALUES ('435', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609127015', '1', '4');
INSERT INTO `im_message` VALUES ('436', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609127237', '1', '4');
INSERT INTO `im_message` VALUES ('437', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609127246', '1', '4');
INSERT INTO `im_message` VALUES ('438', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609127254', '1', '4');
INSERT INTO `im_message` VALUES ('439', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609127411', '1', '4');
INSERT INTO `im_message` VALUES ('440', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609127415', '4', '1');
INSERT INTO `im_message` VALUES ('441', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609127520', '4', '1');
INSERT INTO `im_message` VALUES ('442', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609127526', '4', '1');
INSERT INTO `im_message` VALUES ('443', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609127532', '1', '4');
INSERT INTO `im_message` VALUES ('444', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609127542', '1', '4');
INSERT INTO `im_message` VALUES ('445', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609127560', '1', '4');
INSERT INTO `im_message` VALUES ('446', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609127565', '4', '1');
INSERT INTO `im_message` VALUES ('447', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609127582', '4', '1');
INSERT INTO `im_message` VALUES ('448', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609127602', '4', '1');
INSERT INTO `im_message` VALUES ('449', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609127622', '4', '1');
INSERT INTO `im_message` VALUES ('450', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609127817', '1', '4');
INSERT INTO `im_message` VALUES ('451', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609127825', '4', '1');
INSERT INTO `im_message` VALUES ('452', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609127832', '4', '1');
INSERT INTO `im_message` VALUES ('453', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609127835', '1', '4');
INSERT INTO `im_message` VALUES ('454', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609127938', '1', '4');
INSERT INTO `im_message` VALUES ('455', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609127942', '4', '1');
INSERT INTO `im_message` VALUES ('456', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609127945', '4', '1');
INSERT INTO `im_message` VALUES ('457', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609127949', '1', '4');
INSERT INTO `im_message` VALUES ('458', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609127951', '4', '1');
INSERT INTO `im_message` VALUES ('459', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609173385', '1', '4');
INSERT INTO `im_message` VALUES ('460', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609173388', '1', '4');
INSERT INTO `im_message` VALUES ('461', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609173428', '4', '1');
INSERT INTO `im_message` VALUES ('462', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609173442', '4', '1');
INSERT INTO `im_message` VALUES ('463', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609173444', '4', '1');
INSERT INTO `im_message` VALUES ('464', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/1.jpeg\">', '1609173445', '4', '1');
INSERT INTO `im_message` VALUES ('465', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609173447', '4', '1');
INSERT INTO `im_message` VALUES ('466', '1', '[\\u1f48b]', '1609173450', '4', '1');
INSERT INTO `im_message` VALUES ('467', '1', '[\\u1f64c]', '1609173453', '4', '1');
INSERT INTO `im_message` VALUES ('468', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609173455', '4', '1');
INSERT INTO `im_message` VALUES ('469', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609173459', '4', '1');
INSERT INTO `im_message` VALUES ('470', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609173471', '4', '1');
INSERT INTO `im_message` VALUES ('471', '1', '[\\u1f61a][\\u1f9b7]', '1609175749', '4', '1');
INSERT INTO `im_message` VALUES ('472', '1', '[\\u1f619][\\u1f921]', '1609178007', '1', '11');
INSERT INTO `im_message` VALUES ('473', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/3.jpeg\">', '1609178009', '1', '11');
INSERT INTO `im_message` VALUES ('474', '2', '<img style=\"width:80px;height:80px\" src=\"/web/assets/image/emoji/2.jpeg\">', '1609178014', '1', '11');
INSERT INTO `im_message` VALUES ('475', '1', '哈哈哈', '1609178018', '1', '11');
INSERT INTO `im_message` VALUES ('476', '1', '哈哈哈', '1609178088', '1', '4');
INSERT INTO `im_message` VALUES ('477', '1', '哈哈哈', '1609178091', '1', '4');
INSERT INTO `im_message` VALUES ('478', '1', 'hah', '1609240467', '2', '1');
INSERT INTO `im_message` VALUES ('479', '1', '上野树里', '1609240509', '1', '2');

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
) ENGINE=MyISAM AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- ----------------------------
-- Records of im_user
-- ----------------------------
INSERT INTO `im_user` VALUES ('1', 'zzzy', '/web/static/uploads/images/timg.jpg', '0', '0', 'ha123', '$2a$10$2dY7VOSi9XmG9cWBEJ.jJ.elaDnCtks4ZDWh9xZ5FUiL1NBnlPR1m');
INSERT INTO `im_user` VALUES ('2', 'W', '/web/static/uploads/images/timg0.jpg', '1608368144', '0', 'w', '$2a$10$2dY7VOSi9XmG9cWBEJ.jJ.elaDnCtks4ZDWh9xZ5FUiL1NBnlPR1m');
INSERT INTO `im_user` VALUES ('3', '朱院长', '/web/static/uploads/images/timg.jpg', '1608368145', '0', 'zhu', '$2a$10$2dY7VOSi9XmG9cWBEJ.jJ.elaDnCtks4ZDWh9xZ5FUiL1NBnlPR1m');
INSERT INTO `im_user` VALUES ('4', '天下无敌', '/web/static/uploads/images/timg0.jpg', '0', '0', 'for505', '$2a$10$2dY7VOSi9XmG9cWBEJ.jJ.elaDnCtks4ZDWh9xZ5FUiL1NBnlPR1m');
INSERT INTO `im_user` VALUES ('5', 'cc', '/web/static/uploads/images/timg.jpg', '0', '0', '', '');
INSERT INTO `im_user` VALUES ('6', '莱莱', '/web/static/uploads/images/timg.jpg', '0', '0', '', '');
INSERT INTO `im_user` VALUES ('7', '牛二虎', '/web/static/uploads/images/timg.jpg', '1608373253', '0', '', '');
INSERT INTO `im_user` VALUES ('8', '张一一', '/web/static/uploads/images/timg.jpg', '1608373253', '0', '', '');
INSERT INTO `im_user` VALUES ('9', '胡小小', '/web/static/uploads/images/timg.jpg', '1608373253', '0', '', '');
INSERT INTO `im_user` VALUES ('10', '二狗', '/web/static/uploads/images/timg.jpg', '1608373253', '0', '', '');
INSERT INTO `im_user` VALUES ('11', 'haizi', '/web/static/uploads/images/timg.jpg', '1608373253', '0', '', '');
INSERT INTO `im_user` VALUES ('12', '出水芙蓉', '/web/static/uploads/images/timg0.jpg', '0', '0', '', '');
