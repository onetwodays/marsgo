create table t_pending_accounts
(
    id int auto_increment comment 'pk
'
        primary key,
    number varchar(64) not null comment '手机号',
    verification_code varchar(10) not null comment '验证码',
    push_code varchar(10) not null comment '推送码',
    timestamp bigint not null,
    create_time datetime default CURRENT_TIMESTAMP not null,
    update_time datetime default CURRENT_TIMESTAMP null,
    deleted_at datetime null,
    constraint t_pending_accounts_number_uindex
        unique (number),
    constraint t_pending_accounts_verification_code_uindex
        unique (verification_code)
)
    comment '待注册帐号';

