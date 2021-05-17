


create table t_keys
(
    id bigint auto_increment
        primary key,
    number varchar(255) not null,
    key_id bigint not null,
    public_key text not null,
    last_resort smallint default 0 not null,
    device_id bigint default 1 not null,
    create_time timestamp default CURRENT_TIMESTAMP not null,
    signed_prekey varchar(1024) default '' not null,
    identity_key varchar(1024) default '' not null
);

create index keys_number_index
	on t_keys (number);

