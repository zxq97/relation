CREATE TABLE user_follow
(
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `uid`         BIGINT    NOT NULL DEFAULT 0,
    `to_uid`      BIT       NOT NULL DEFAULT 0,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE uniq_uid_touid (`uid`, `to_uid`),
    KEY           idx_uid_createtime (`uid`, `create_time`)
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE user_follower
(
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `uid`         BIGINT    NOT NULL DEFAULT 0,
    `to_uid`      BIT       NOT NULL DEFAULT 0,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE uniq_uid_touid (`uid`, `to_uid`),
    KEY           idx_uid_createtime (`uid`, `create_time`)
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE user_follow_count
(
    `uid`            BIGINT    NOT NULL DEFAULT 0,
    `follow_count`   INT       NOT NULL DEFAULT 0,
    `follower_count` INT       NOT NULL DEFAULT 0,
    `create_time`    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`uid`)
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
