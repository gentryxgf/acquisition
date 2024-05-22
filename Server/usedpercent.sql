CREATE TABLE `used_percent`
(
    `id`        BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
    `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `metric`    varchar(200)  NOT NULL DEFAULT '' COMMENT '指标名称',
    `endpoint`  varchar(200)  NOT NULL DEFAULT '' COMMENT '当前主机名称',
    `timestamp` int UNSIGNED NOT NULL DEFAULT '0' COMMENT '采集数据时的时间',
    `step`      int UNSIGNED NOT NULL DEFAULT '0' COMMENT '指标的采集周期',
    `value`     float UNSIGNED NOT NULL DEFAULT '0.0' COMMENT '利用率值%',
    `extend`    varchar(200)  NOT NULL DEFAULT '' COMMENT '指标名称'
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT = '利用率表';