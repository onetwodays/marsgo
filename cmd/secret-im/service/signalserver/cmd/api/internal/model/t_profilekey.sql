create table t_profilekey
(
    id int auto_increment
        primary key,
    account_name varchar(64) not null,
    profile_key varchar(1024) default '' not null,
    constraint t_profilekey_account_name_uindex
        unique (account_name)
);

