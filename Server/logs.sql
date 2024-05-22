CREATE TABLE `logs`
(
    `id`        BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
    `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `hostname`    varchar(200)  NOT NULL DEFAULT '' COMMENT '主机名称',
    `file`  varchar(200)  NOT NULL DEFAULT '' COMMENT '日志文件路径',
    `log`  varchar(500)  NOT NULL DEFAULT '' COMMENT '日志内容'
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT = '日志表';