-- ----------------------------
-- Table structure for blog_article
-- ----------------------------
DROP TABLE IF EXISTS `blog_article`;
CREATE TABLE `blog_article`
(
    `id`              int unsigned NOT NULL AUTO_INCREMENT,
    `tag_id`          int(10) unsigned DEFAULT '0' COMMENT '标签ID',
    `title`           varchar(100)     DEFAULT '' COMMENT '文章标题',
    `desc`            varchar(255)     DEFAULT '' COMMENT '简述',
    `content`         text COMMENT '内容',
    `cover_image_url` varchar(255)     DEFAULT '' COMMENT '封面图片地址',
    `created_by`      varchar(128)     DEFAULT '' COMMENT '创建人',
    `modified_by`     varchar(128)     DEFAULT '' COMMENT '修改人',
    `deleted`         boolean NOT NULL DEFAULT false COMMENT '是否删除',
    `created_at`      int(10) unsigned DEFAULT '0' COMMENT '创建时间',
    `updated_at`      int(10) unsigned DEFAULT '0' COMMENT '修改时间',
    `deleted_at`      int(10) unsigned DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY               IDX_TAG_ID (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章管理';

-- ----------------------------
-- Table structure for blog_tag
-- ----------------------------
DROP TABLE IF EXISTS `blog_tag`;
CREATE TABLE `blog_tag`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(100)     DEFAULT '' COMMENT '标签名称',
    `created_by`  varchar(128)     DEFAULT '' COMMENT '创建人',
    `modified_by` varchar(128)     DEFAULT '' COMMENT '修改人',
    `disabled`    boolean NOT NULL DEFAULT false COMMENT '是否禁用',
    `deleted`     boolean NOT NULL DEFAULT false COMMENT '是否删除',
    `created_at`  int(10) unsigned DEFAULT '0' COMMENT '创建时间',
    `updated_at`  int(10) unsigned DEFAULT '0' COMMENT '修改时间',
    `deleted_at`  int(10) unsigned DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章标签管理';

-- ----------------------------
-- Table structure for rbac
-- ----------------------------
DROP TABLE IF EXISTS `blog_user`;
CREATE TABLE `blog_user`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT,
    `username`   varchar(128) NOT NULL DEFAULT '' COMMENT '账号',
    `password`   varchar(128) NOT NULL DEFAULT '' COMMENT '密码',
    `phone`      varchar(256)          DEFAULT '' COMMENT '手机号',
    `email`      varchar(256)          DEFAULT '' COMMENT '邮箱',
    `locked`     boolean      NOT NULL DEFAULT false COMMENT '是否锁定',
    `deleted`    boolean      NOT NULL DEFAULT false COMMENT '是否删除',
    `created_at` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
    `updated_at` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
    `deleted_at` int(10) unsigned DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY          IDX_USERNAME (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `blog_group`;
CREATE TABLE `blog_group`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(128) NOT NULL DEFAULT '' COMMENT '用户组名称',
    `description` varchar(256) NOT NULL DEFAULT '' COMMENT '用户组描述',
    `locked`      boolean      NOT NULL DEFAULT false COMMENT '是否锁定',
    `deleted`     boolean      NOT NULL DEFAULT false COMMENT '是否删除',
    `created_at`  int(10) unsigned DEFAULT '0' COMMENT '创建时间',
    `updated_at`  int(10) unsigned DEFAULT '0' COMMENT '修改时间',
    `deleted_at`  int(10) unsigned DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `blog_role`;
CREATE TABLE `blog_role`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(128) NOT NULL DEFAULT '' COMMENT '角色名称',
    `description` varchar(256) NOT NULL DEFAULT '' COMMENT '角色描述',
    `locked`      boolean      NOT NULL DEFAULT false COMMENT '是否锁定',
    `deleted`     boolean      NOT NULL DEFAULT false COMMENT '是否删除',
    `created_at`  int(10) unsigned DEFAULT '0' COMMENT '创建时间',
    `updated_at`  int(10) unsigned DEFAULT '0' COMMENT '修改时间',
    `deleted_at`  int(10) unsigned DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `blog_resource`;
CREATE TABLE `blog_resource`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
    `url`         varchar(256) NOT NULL DEFAULT '' COMMENT '资源 url',
    `name`        varchar(128) NOT NULL DEFAULT '' COMMENT '资源名称',
    `description` varchar(256) NOT NULL DEFAULT '' COMMENT '资源描述',
    `deleted`     boolean      NOT NULL DEFAULT false COMMENT '是否删除',
    `created_at`  int(10) unsigned DEFAULT '0' COMMENT '创建时间',
    `updated_at`  int(10) unsigned DEFAULT '0' COMMENT '修改时间',
    `deleted_at`  int(10) unsigned DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `blog_operation`;
CREATE TABLE `blog_operation`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(128) NOT NULL DEFAULT '' COMMENT '操作名称',
    `description` varchar(256) NOT NULL DEFAULT '' COMMENT '操作描述',
    `deleted`     boolean      NOT NULL DEFAULT false COMMENT '是否删除',
    `created_at`  int(10) unsigned DEFAULT '0' COMMENT '创建时间',
    `updated_at`  int(10) unsigned DEFAULT '0' COMMENT '修改时间',
    `deleted_at`  int(10) unsigned DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `blog_permission`;
CREATE TABLE `blog_permission`
(
    `id`           int(10) unsigned NOT NULL AUTO_INCREMENT,
    `resource_id`  int(10) NOT NULL DEFAULT '0' COMMENT '资源 ID',
    `operation_id` int(10) NOT NULL DEFAULT '0' COMMENT '操作 ID',
    `deleted`      boolean NOT NULL DEFAULT false COMMENT '是否删除',
    `created_at`   int(10) unsigned DEFAULT '0' COMMENT '创建时间',
    `updated_at`   int(10) unsigned DEFAULT '0' COMMENT '修改时间',
    `deleted_at`   int(10) unsigned DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `blog_group_user`;
CREATE TABLE `blog_group_user`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT,
    `group_id`   int(10) NOT NULL DEFAULT '0' COMMENT '用户组 ID',
    `user_id`    int     NOT NULL DEFAULT '0' COMMENT '用户 ID',
    `deleted`    boolean NOT NULL DEFAULT false COMMENT '是否删除',
    `created_at` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
    `updated_at` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
    `deleted_at` int(10) unsigned DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `blog_group_role`;
CREATE TABLE `blog_group_role`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
    `group_id`   int(10) NOT NULL DEFAULT '0' COMMENT '用户组 ID',
    `role_id`    int(10) NOT NULL DEFAULT '0' COMMENT '角色 ID',
    `deleted`    boolean NOT NULL DEFAULT false COMMENT '是否删除',
    `created_at` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
    `updated_at` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
    `deleted_at` int(10) unsigned DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `blog_role_permission`;
CREATE TABLE `blog_role_permission`
(
    `id`            int(10) unsigned NOT NULL AUTO_INCREMENT,
    `role_id`       int(10) NOT NULL DEFAULT '0' COMMENT '角色 ID',
    `permission_id` int(10) NOT NULL DEFAULT '0' COMMENT '权限 ID',
    `deleted`       boolean NOT NULL DEFAULT false COMMENT '是否删除',
    `created_at`    int(10) unsigned DEFAULT '0' COMMENT '创建时间',
    `updated_at`    int(10) unsigned DEFAULT '0' COMMENT '修改时间',
    `deleted_at`    int(10) unsigned DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `blog_user` (`id`, `username`, `password`, `phone`)
VALUES ('1', 'admin', 'Admin@111', '15666666666');

INSERT INTO `blog_group` (`id`, `name`, `description`)
VALUES ('1', 'provider', 'default admin group');

INSERT INTO `blog_role` (`id`, `name`, `description`)
VALUES ('1', 'admin', 'default admin role');

INSERT INTO `blog_resource` (`id`, `url`, `name`, `description`)
VALUES ('1', '/api/v1//tags', 'tags', 'blog tags'),
       ('2', '/api/v1//users', 'users', 'blog users'),
       ('3', '/api/v1//images', 'images', 'upload images');

INSERT INTO `blog_operation` (`id`, `name`, `description`)
VALUES ('1', 'POST', 'post method'),
       ('2', 'GET', 'get method'),
       ('3', 'PUT', 'put method'),
       ('4', 'DELETE', 'delete method');

INSERT INTO `blog_permission` (`id`, `resource_id`, `operation_id`)
VALUES ('1', '1'),
       ('1', '2'),
       ('1', '3'),
       ('2', '1'),
       ('2', '2'),
       ('2', '3'),
       ('3', '1'),
       ('3', '2'),
       ('3', '3');

INSERT INTO `blog_group_user` (`id`, `group_id`, `user_id`)
VALUES ('1', '1', '1');

INSERT INTO `blog_group_role` (`id`, `group_id`, `role_id`)
VALUES ('1', '1', '1');

INSERT INTO `blog_role_permission` (`role_id`, `permission_id`)
VALUES ('1', '1'),
       ('1', '2'),
       ('1', '3'),
       ('1', '5'),
       ('1', '6'),
       ('1', '7'),
       ('1', '9'),
       ('1', '10'),
       ('1', '11');