create table t_accounts
(
    id int auto_increment comment 'pk'
        primary key,
    number varchar(40) not null,
    uuid varchar(60) not null,
    data text null,
    create_time datetime default CURRENT_TIMESTAMP null,
    update_time datetime default CURRENT_TIMESTAMP null,
    delete_time datetime null,
    constraint t_acconts_number_uindex
        unique (number),
    constraint t_acconts_uuid_uindex
        unique (uuid)
)
    comment '账户表';

