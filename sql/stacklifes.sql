/*
Navicat MySQL Data Transfer

Source Server         : 本地
Source Server Version : 50744
Source Host           : localhost:3306
Source Database       : stacklifes

Target Server Type    : MYSQL
Target Server Version : 50744
File Encoding         : 65001

Date: 2026-05-12 13:03:16
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for sl_admin
-- ----------------------------
DROP TABLE IF EXISTS `sl_admin`;
CREATE TABLE `sl_admin` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户自增id',
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(64) NOT NULL DEFAULT '' COMMENT '加密密码',
  `salt` varchar(64) NOT NULL COMMENT '密码掩码',
  `group_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户组id',
  `name` varchar(64) NOT NULL COMMENT '姓名',
  `update_time` int(11) DEFAULT NULL COMMENT '更新时间',
  `create_time` int(11) DEFAULT '0' COMMENT '创建用户时间戳',
  `expiration_time` int(11) DEFAULT '0' COMMENT '用户过期时间戳',
  `login_num` int(11) NOT NULL DEFAULT '0' COMMENT '累计登陆次数',
  `last_login_time` int(11) DEFAULT '0' COMMENT '最后登录时间',
  `last_login_ip` varchar(64) DEFAULT '' COMMENT '最后登录ip',
  `lang` varchar(255) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '0' COMMENT '用户状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=382 DEFAULT CHARSET=utf8 COMMENT='用户表';

-- ----------------------------
-- Table structure for sl_admin_group
-- ----------------------------
DROP TABLE IF EXISTS `sl_admin_group`;
CREATE TABLE `sl_admin_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户组自增id',
  `mark` varchar(50) NOT NULL COMMENT '用户组标识',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户组名称',
  `description` varchar(256) NOT NULL DEFAULT '' COMMENT '描述',
  `menu_permit` text COMMENT '央服权限数据',
  `menu_modules` text COMMENT '央服权限模块',
  `allow_ip` text COMMENT 'ip白名单',
  `status` int(11) DEFAULT '1',
  `create_time` int(11) DEFAULT '0' COMMENT '创建用户组时间戳',
  `update_time` int(11) DEFAULT NULL COMMENT '更新操作时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=155 DEFAULT CHARSET=utf8 COMMENT='用户组';

-- ----------------------------
-- Table structure for sl_article
-- ----------------------------
DROP TABLE IF EXISTS `sl_article`;
CREATE TABLE `sl_article` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cid` int(11) DEFAULT NULL,
  `did` int(11) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `desc` varchar(255) DEFAULT NULL,
  `keyword` varchar(255) DEFAULT NULL,
  `author` varchar(255) DEFAULT NULL,
  `thumb` varchar(255) DEFAULT NULL,
  `summary` text,
  `content` longtext,
  `is_hot` tinyint(4) DEFAULT NULL,
  `is_new` tinyint(4) DEFAULT NULL,
  `is_recom` tinyint(4) DEFAULT NULL,
  `weight` int(11) DEFAULT NULL,
  `public_time` int(11) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL COMMENT '素材类型，1图片，2视频',
  `month` varchar(255) DEFAULT NULL,
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=316 DEFAULT CHARSET=utf8 COMMENT='文章表';

-- ----------------------------
-- Table structure for sl_article_media
-- ----------------------------
DROP TABLE IF EXISTS `sl_article_media`;
CREATE TABLE `sl_article_media` (
  `id` int(11) DEFAULT NULL,
  `platform_id` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '媒体类型',
  `aid` int(11) DEFAULT NULL COMMENT '文章id',
  `media_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `error_msg` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `publish_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` int(11) DEFAULT NULL COMMENT '发布状态，0，未发布，1，已发布，2发布失败',
  `update_time` int(11) DEFAULT NULL COMMENT '更新时间',
  `create_time` int(11) DEFAULT NULL COMMENT '发布时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for sl_article_tags
-- ----------------------------
DROP TABLE IF EXISTS `sl_article_tags`;
CREATE TABLE `sl_article_tags` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `aid` int(11) DEFAULT NULL,
  `tid` int(11) DEFAULT NULL,
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8 COMMENT='文章标签关联表';

-- ----------------------------
-- Table structure for sl_attach
-- ----------------------------
DROP TABLE IF EXISTS `sl_attach`;
CREATE TABLE `sl_attach` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `link` varchar(255) DEFAULT NULL,
  `path` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `type` tinyint(4) DEFAULT NULL COMMENT '素材类型，1图片，2视频',
  `is_thumb` tinyint(4) DEFAULT '0' COMMENT '是否封面图，1是，0否',
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8 COMMENT='附件表';

-- ----------------------------
-- Table structure for sl_attach_media
-- ----------------------------
DROP TABLE IF EXISTS `sl_attach_media`;
CREATE TABLE `sl_attach_media` (
  `id` int(11) DEFAULT NULL,
  `platform_id` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '媒体类型',
  `attach_id` int(11) DEFAULT NULL COMMENT '文章id',
  `media_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `error_msg` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `media_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` int(11) DEFAULT NULL COMMENT '发布状态，0，未发布，1，已发布，2发布失败',
  `update_time` int(11) DEFAULT NULL COMMENT '更新时间',
  `create_time` int(11) DEFAULT NULL COMMENT '发布时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for sl_category
-- ----------------------------
DROP TABLE IF EXISTS `sl_category`;
CREATE TABLE `sl_category` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '菜单自增id',
  `type` varchar(16) DEFAULT NULL COMMENT '菜单类型，中央服，单服',
  `parent` varchar(255) DEFAULT '0',
  `mark` varchar(60) DEFAULT NULL,
  `author` varchar(32) NOT NULL DEFAULT '' COMMENT '创建人',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '菜单名称',
  `weight` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `create_time` int(11) DEFAULT '0' COMMENT '创建菜单时间戳',
  `update_time` int(11) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`) USING BTREE,
  UNIQUE KEY `mark` (`mark`) USING BTREE,
  KEY `weight` (`weight`) USING BTREE,
  KEY `type_status` (`type`,`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=379 DEFAULT CHARSET=utf8 COMMENT='文字分类表';

-- ----------------------------
-- Table structure for sl_comment
-- ----------------------------
DROP TABLE IF EXISTS `sl_comment`;
CREATE TABLE `sl_comment` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `aid` int(11) DEFAULT NULL,
  `referer` varchar(255) DEFAULT NULL,
  `ip` varchar(50) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `url` varchar(255) DEFAULT NULL,
  `is_reply` int(11) DEFAULT NULL,
  `content` text,
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='文章评论表';

-- ----------------------------
-- Table structure for sl_directory
-- ----------------------------
DROP TABLE IF EXISTS `sl_directory`;
CREATE TABLE `sl_directory` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '菜单自增id',
  `cid` int(11) DEFAULT NULL,
  `type` varchar(16) DEFAULT NULL COMMENT '菜单类型，中央服，单服',
  `parent` varchar(255) DEFAULT '0',
  `mark` varchar(60) DEFAULT NULL,
  `author` varchar(32) NOT NULL DEFAULT '' COMMENT '创建人',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '菜单名称',
  `weight` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `create_time` int(11) DEFAULT '0' COMMENT '创建菜单时间戳',
  `update_time` int(11) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`) USING BTREE,
  UNIQUE KEY `mark` (`mark`) USING BTREE,
  KEY `weight` (`weight`) USING BTREE,
  KEY `type_status` (`type`,`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=405 DEFAULT CHARSET=utf8 COMMENT='文章分类目录表';

-- ----------------------------
-- Table structure for sl_lang
-- ----------------------------
DROP TABLE IF EXISTS `sl_lang`;
CREATE TABLE `sl_lang` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `lang` varchar(16) DEFAULT NULL,
  `default` tinyint(4) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `author` varchar(50) DEFAULT NULL,
  `weight` int(11) DEFAULT NULL,
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `lang` (`lang`) USING BTREE
) ENGINE=MyISAM AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for sl_log
-- ----------------------------
DROP TABLE IF EXISTS `sl_log`;
CREATE TABLE `sl_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '日志表自增id',
  `username` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
  `ip` varchar(32) DEFAULT NULL,
  `content` varchar(500) NOT NULL DEFAULT '' COMMENT '简要描述',
  `create_time` int(11) DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `username` (`username`),
  KEY `create_time` (`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=97 DEFAULT CHARSET=utf8 COMMENT='日志表';

-- ----------------------------
-- Table structure for sl_menu
-- ----------------------------
DROP TABLE IF EXISTS `sl_menu`;
CREATE TABLE `sl_menu` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '菜单自增id',
  `type` int(11) DEFAULT '1',
  `mark` varchar(60) DEFAULT NULL,
  `author` varchar(32) DEFAULT '' COMMENT '创建人',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '菜单名称',
  `url` varchar(256) NOT NULL DEFAULT '' COMMENT '链接',
  `parent` int(11) NOT NULL DEFAULT '0' COMMENT '上级菜单id',
  `icon` varchar(128) DEFAULT '' COMMENT '菜单图标',
  `weight` int(11) DEFAULT '0' COMMENT '排序',
  `status` int(11) DEFAULT '1' COMMENT '状态',
  `create_time` int(11) DEFAULT '0' COMMENT '创建菜单时间戳',
  `update_time` int(11) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`) USING BTREE,
  UNIQUE KEY `mark` (`mark`) USING BTREE,
  KEY `parent` (`parent`) USING BTREE,
  KEY `weight` (`weight`) USING BTREE,
  KEY `type_status` (`type`,`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=375 DEFAULT CHARSET=utf8 COMMENT='菜单';

-- ----------------------------
-- Table structure for sl_message
-- ----------------------------
DROP TABLE IF EXISTS `sl_message`;
CREATE TABLE `sl_message` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip` varchar(50) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `url` varchar(255) DEFAULT NULL,
  `is_reply` tinyint(4) DEFAULT NULL,
  `content` text,
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='留言表';

-- ----------------------------
-- Table structure for sl_permit
-- ----------------------------
DROP TABLE IF EXISTS `sl_permit`;
CREATE TABLE `sl_permit` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '权限表自增id',
  `type` int(11) NOT NULL DEFAULT '1',
  `mark` varchar(255) DEFAULT NULL,
  `parent` int(11) NOT NULL DEFAULT '0',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '权限名称',
  `modules` text COMMENT '模块',
  `author` varchar(255) DEFAULT NULL,
  `status` int(11) DEFAULT '1',
  `weight` int(11) DEFAULT NULL,
  `create_time` int(11) DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=391 DEFAULT CHARSET=utf8 COMMENT='权限表';

-- ----------------------------
-- Table structure for sl_platform
-- ----------------------------
DROP TABLE IF EXISTS `sl_platform`;
CREATE TABLE `sl_platform` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `mark` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '标识',
  `name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '明天',
  `status` tinyint(4) DEFAULT NULL COMMENT '状态',
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for sl_read
-- ----------------------------
DROP TABLE IF EXISTS `sl_read`;
CREATE TABLE `sl_read` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `aid` int(11) DEFAULT NULL,
  `referer` varchar(255) DEFAULT NULL,
  `ip` varchar(50) DEFAULT NULL,
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8695 DEFAULT CHARSET=utf8 COMMENT='文章阅读记录表';

-- ----------------------------
-- Table structure for sl_tags
-- ----------------------------
DROP TABLE IF EXISTS `sl_tags`;
CREATE TABLE `sl_tags` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '菜单自增id',
  `type` varchar(16) DEFAULT NULL COMMENT '菜单类型，中央服，单服',
  `mark` varchar(60) DEFAULT NULL,
  `author` varchar(32) NOT NULL DEFAULT '' COMMENT '创建人',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '菜单名称',
  `weight` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `is_hot` tinyint(1) DEFAULT '1',
  `create_time` int(11) DEFAULT '0' COMMENT '创建菜单时间戳',
  `update_time` int(11) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`) USING BTREE,
  UNIQUE KEY `mark` (`mark`) USING BTREE,
  KEY `weight` (`weight`) USING BTREE,
  KEY `type_status` (`type`,`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=392 DEFAULT CHARSET=utf8 COMMENT='标签表';

-- ----------------------------
-- Table structure for sl_website
-- ----------------------------
DROP TABLE IF EXISTS `sl_website`;
CREATE TABLE `sl_website` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `key` varchar(32) DEFAULT NULL,
  `value` text,
  `staus` int(11) DEFAULT '1',
  `update_time` int(11) DEFAULT NULL,
  `create_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
