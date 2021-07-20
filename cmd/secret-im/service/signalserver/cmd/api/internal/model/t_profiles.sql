create table t_profiles
(
    id bigint auto_increment
        primary key,
    uuid varchar(64) not null comment 'uuid',
    version varchar(1024) not null comment 'prekey version',
    name text not null comment 'name',
    avatar varchar(128) null comment '头像文件相对路径',
    commitment text null,
    create_time datetime default CURRENT_TIMESTAMP null,
    update_time datetime default CURRENT_TIMESTAMP null
);

create index t_profiles_uuid_index
    on t_profiles (uuid);

