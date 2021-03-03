create table t_messages
(
    id int auto_increment comment 'pk',
    type int not null comment '消息类型',
    relay varchar(512)  default ''  comment '类似cdn',
    tm int not null comment 'unix  时间戳',
    source varchar(20) not null comment '源手机号',
    source_device int not null comment '源手机号绑定的设备id',
    message blob null comment '消息',
    content blob not null,
    destination varchar(20) not null comment '目的手机号',
    create_time datetime default current_timestamp,
    PRIMARY KEY (`id`)
) comment '保存离线消息';

