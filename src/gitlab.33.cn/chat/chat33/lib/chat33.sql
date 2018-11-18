/*
SQLyog Ultimate v12.09 (64 bit)
MySQL - 5.6.37 : Database - chat33
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`chat33` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `chat33`;

/*Table structure for table `apply` */

DROP TABLE IF EXISTS `apply`;

CREATE TABLE `apply` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` tinyint(255) DEFAULT NULL COMMENT '1群 2好友',
  `apply_user` int(255) DEFAULT NULL,
  `target` int(255) DEFAULT NULL,
  `apply_reason` varchar(255) DEFAULT NULL,
  `state` tinyint(255) DEFAULT NULL COMMENT '1待处理   2拒绝   3同意',
  `remark` varchar(255) DEFAULT NULL,
  `datetime` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_type_apply_user_target` (`type`,`apply_user`,`target`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4;

/*Table structure for table `chat_log` */

DROP TABLE IF EXISTS `chat_log`;

CREATE TABLE `chat_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sender_id` varchar(255) DEFAULT NULL COMMENT '发送者id',
  `receive_id` varchar(255) DEFAULT NULL COMMENT '聊天室id',
  `msg_type` int(10) unsigned DEFAULT NULL COMMENT '消息类型',
  `content` varchar(7168) DEFAULT NULL COMMENT '内容',
  `log_type` int(255) DEFAULT NULL,
  `send_time` bigint(13) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=684 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT;

/*Table structure for table `coin` */

DROP TABLE IF EXISTS `coin`;

CREATE TABLE `coin` (
  `coin_id` int(255) NOT NULL,
  `coin_name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`coin_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT;

-- ----------------------------
-- Records of coin
-- ----------------------------
INSERT INTO `coin` (coin_id,coin_name) VALUES (1,'BTY');
INSERT INTO `coin` (coin_id,coin_name) VALUES (2,'YCC');

/*Table structure for table `friends` */

DROP TABLE IF EXISTS `friends`;

CREATE TABLE `friends` (
  `user_id` int(11) NOT NULL,
  `friend_id` int(11) NOT NULL,
  `remark` varchar(45) DEFAULT NULL COMMENT '好友备注',
  `add_time` bigint(20) NOT NULL DEFAULT '0',
  `DND` tinyint(4) NOT NULL DEFAULT '2' COMMENT '是否消息免打扰 1免打扰  2关闭',
  `top` int(3) NOT NULL DEFAULT '2' COMMENT '好友置顶  1置顶 2不置顶',
  `type` int(3) DEFAULT '1' COMMENT '1 普通  2 常用',
  `is_delete` int(3) NOT NULL DEFAULT '1' COMMENT '1 未删除  2 已删除',
  PRIMARY KEY (`user_id`,`friend_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

/*Table structure for table `group` */

DROP TABLE IF EXISTS `group`;

CREATE TABLE `group` (
  `group_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `group_name` varchar(45) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `status` tinyint(4) DEFAULT NULL COMMENT '开启、关闭',
  `description` varchar(255) DEFAULT NULL,
  `open_time` bigint(3) DEFAULT NULL,
  `close_time` bigint(3) DEFAULT NULL,
  `create_time` bigint(3) DEFAULT NULL,
  PRIMARY KEY (`group_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT;

/*Table structure for table `login_log` */

DROP TABLE IF EXISTS `login_log`;

CREATE TABLE `login_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `app_id` int(11) DEFAULT NULL COMMENT '未用',
  `login_time` bigint(20) NOT NULL,
  `device` varchar(30) DEFAULT NULL COMMENT '登录设备',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2010 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

/*Table structure for table `private_chat_log` */

DROP TABLE IF EXISTS `private_chat_log`;

CREATE TABLE `private_chat_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sender_id` varchar(255) DEFAULT NULL,
  `receive_id` varchar(255) DEFAULT NULL,
  `msg_type` int(10) unsigned DEFAULT NULL,
  `content` varchar(7168) DEFAULT NULL,
  `status` int(11) DEFAULT NULL COMMENT '1 已读  2未读',
  `send_time` bigint(13) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=206 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT;

/*Table structure for table `red_packet_log` */

DROP TABLE IF EXISTS `red_packet_log`;

CREATE TABLE `red_packet_log` (
  `packet_id` varchar(64) NOT NULL,
  `ctype` tinyint(255) NOT NULL COMMENT '红包发到群/用户/聊天室',
  `user_id` varchar(64) NOT NULL COMMENT '发红包的用户id',
  `to_id` varchar(64) NOT NULL,
  `coin` int(11) NOT NULL,
  `size` int(11) NOT NULL COMMENT '几个红包',
  `amount` decimal(11,0) NOT NULL COMMENT '总金额',
  `remark` varchar(64) DEFAULT NULL COMMENT '红包备注',
  `created_at` bigint(13) NOT NULL COMMENT '时间',
  `type` int(11) NOT NULL COMMENT '新人、拼手气红包',
  PRIMARY KEY (`packet_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT COMMENT='发红包的记录';

/*Table structure for table `room` */

DROP TABLE IF EXISTS `room`;

CREATE TABLE `room` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `mark_id` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `master_id` int(11) DEFAULT NULL,
  `create_time` bigint(20) DEFAULT NULL,
  `can_add_friend` tinyint(1) DEFAULT '1' COMMENT '1可添加好友 2不可添加好友',
  `join_permission` tinyint(1) DEFAULT '2' COMMENT '1需要审批 2不需要审批  3禁止加入',
  `admin_muted` tinyint(1) DEFAULT NULL,
  `master_muted` tinyint(255) DEFAULT NULL,
  `is_delete` tinyint(255) DEFAULT '1' COMMENT '1 未删除 2 删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4;

/*Table structure for table `room_join_request` */

DROP TABLE IF EXISTS `room_join_request`;

CREATE TABLE `room_join_request` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `room_id` int(11) DEFAULT NULL,
  `apply_user` int(255) DEFAULT NULL,
  `apply_reason` varchar(255) DEFAULT NULL,
  `inviter` int(255) DEFAULT NULL,
  `status` int(255) DEFAULT NULL,
  `time` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

/*Table structure for table `room_msg_content` */

DROP TABLE IF EXISTS `room_msg_content`;

CREATE TABLE `room_msg_content` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `room_id` int(11) DEFAULT NULL,
  `sender_id` int(11) DEFAULT NULL,
  `msg_type` int(255) DEFAULT NULL COMMENT '0：系统消息，1:文字，2:音频，3：图片，4：红包，5：视频',
  `content` varchar(255) DEFAULT NULL,
  `datetime` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=210 DEFAULT CHARSET=utf8mb4;

/*Table structure for table `room_msg_receive` */

DROP TABLE IF EXISTS `room_msg_receive`;

CREATE TABLE `room_msg_receive` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `room_msg_id` int(11) DEFAULT NULL,
  `receive_id` int(11) DEFAULT NULL,
  `state` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=157 DEFAULT CHARSET=utf8mb4;

/*Table structure for table `room_user` */

DROP TABLE IF EXISTS `room_user`;

CREATE TABLE `room_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `room_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `user_nickname` varchar(255) DEFAULT NULL,
  `level` int(11) DEFAULT NULL COMMENT '1:普通用户 2：管理员 3：群主',
  `no_disturbing` tinyint(1) DEFAULT NULL COMMENT '1：开启了免打扰，2：关闭',
  `common_use` tinyint(255) DEFAULT NULL COMMENT '1 普通 ；2 常用',
  `room_top` tinyint(255) DEFAULT NULL COMMENT '1 置顶 2不置顶 ',
  `create_time` bigint(20) DEFAULT NULL,
  `is_delete` tinyint(255) DEFAULT NULL COMMENT '1 未删除 2 删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_room_id_user_id` (`room_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4;

/*Table structure for table `user` */

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `user_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mark_id` varchar(255) DEFAULT NULL COMMENT '号码',
  `uid` int(11) NOT NULL COMMENT '找币uid',
  `app_id` int(11) DEFAULT NULL COMMENT '不用',
  `username` varchar(255) DEFAULT NULL,
  `account` varchar(255) DEFAULT NULL COMMENT '手机号或者邮箱  找币那边的',
  `user_level` int(11) DEFAULT NULL COMMENT '0 游客 1 普通用户 2 客服 3 管理员  数据库里不存游客  然后接口返回的话 客服和管理员都是2',
  `verified` int(11) DEFAULT NULL COMMENT '是否实名  找币',
  `description` varchar(255) DEFAULT NULL COMMENT '管理员给的备注',
  `avatar` varchar(255) DEFAULT NULL,
  `sex` tinyint(4) DEFAULT NULL,
  `phone` varchar(45) DEFAULT NULL,
  `email` varchar(45) DEFAULT NULL,
  `com_id` int(11) DEFAULT NULL COMMENT '公司id 未用到',
  `getui_cid` varchar(255) DEFAULT NULL,
  `position` varchar(20) DEFAULT NULL COMMENT '职位',
  PRIMARY KEY (`user_id`) USING BTREE,
  UNIQUE KEY `uid_UNIQUE` (`uid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=COMPACT;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
