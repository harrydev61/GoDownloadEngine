
CREATE TABLE IF NOT EXISTS `roles`
(
    id                INT                  NOT NULL AUTO_INCREMENT,
    status            TINYINT(3) DEFAULT 1 NOT NULL,
    name              VARCHAR(255)         NOT NULL,
    permission_locked TINYINT(1) DEFAULT 0 NOT NULL,
    description       VARCHAR(500)         NULL,
    created_at        DATETIME   DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME   DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY name (name)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8
    COLLATE = utf8_unicode_ci;

INSERT INTO roles (status, name, permission_locked, description)
VALUES ( 1, 'Super Administrator', 1, 'Super administrator role');
INSERT INTO roles (status, name, permission_locked, description)
VALUES ( 1, 'Administrator', 0, 'administrator role');
INSERT INTO roles ( status, name, permission_locked, description)
VALUES (1, 'User', 0, 'Regular user role');
INSERT INTO roles ( status, name, permission_locked, description)
VALUES ( 1, 'Guest', 0, 'Guest role');


CREATE TABLE IF NOT EXISTS `users`
(
    `user_id`    VARCHAR(36)                                 NOT NULL UNIQUE,
    `first_name` varchar(30) CHARACTER SET utf8mb4           NOT NULL,
    `last_name`  varchar(30) CHARACTER SET utf8mb4           NOT NULL,
    `email`      varchar(255) CHARACTER SET utf8mb4          NOT NULL,
    `phone`      varchar(30)                      DEFAULT NULL,
    `avatar`     text                             DEFAULT NULL,
    `role`       int                              DEFAULT 4  NOT NULL,
    `gender`     enum ('male','female','unknown') DEFAULT 'unknown',
    `dob`        date                             DEFAULT NULL,
    `status`     TINYINT(5)                       DEFAULT -1 NOT NULL,
    `ip`         varchar(255)                                NOT NULL,
    `is_deleted` TINYINT(2)                       DEFAULT 0  NOT NULL,
    `created_at` datetime                         DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime                         DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `email` (`email`) USING BTREE,
    FOREIGN KEY (role) references roles (id)
                                                                            on update cascade on delete cascade
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8
    COLLATE = utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS `auths`
(
    `user_id`     VARCHAR(36)                                  NOT NULL,
    `email`       varchar(255) CHARACTER SET utf8mb4           NOT NULL,
    `salt`        varchar(40) CHARACTER SET utf8mb4  DEFAULT NULL,
    `password`    varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
    `facebook_id` varchar(35) CHARACTER SET utf8mb4  DEFAULT NULL,
    `public_key`  text                                         null,
    `private_key` text                                         null,
    `created_at`  datetime                           DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  datetime                           DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`),
    KEY `email` (`email`) USING BTREE,
    KEY `facebook_id` (`facebook_id`) USING BTREE
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4;



CREATE TABLE IF NOT EXISTS `token_public_keys` (
                                                   `id` BIGINT UNSIGNED AUTO_INCREMENT,
                                                   `public_key` TEXT NOT NULL,
                                                   `created_at`  datetime                           DEFAULT CURRENT_TIMESTAMP,
                                                   `updated_at`  datetime                           DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                                   PRIMARY KEY (`id`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `download_tasks` (
                                                `download_id`       VARCHAR(36)  NOT NULL,
    `user_id`     VARCHAR(36)  NOT NULL,
    `name`     VARCHAR(500)  NOT NULL,
    `description`     VARCHAR(1000),
    `download_type` SMALLINT NOT NULL,
    `url` TEXT NOT NULL,
    `download_status` SMALLINT NOT NULL,
    `is_deleted` TINYINT(2)                       DEFAULT 0  NOT NULL,
    `metadata` TEXT NOT NULL,
    `created_at`  datetime                           DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  datetime                           DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`download_id`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;
