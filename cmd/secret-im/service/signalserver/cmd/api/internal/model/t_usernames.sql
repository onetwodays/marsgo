create table t_usernames
(
    id bigint auto_increment primary key,
    uuid varchar(64) not null,
    username varchar(256) not null,
    create_time datetime default CURRENT_TIMESTAMP null,
    update_time datetime default CURRENT_TIMESTAMP null,
    constraint t_usernames_username_uindex
        unique (username),
    constraint t_usernames_uuid_uindex
        unique (uuid)
)
    comment 'uuid与用户昵称的映射表';

