create table t_messages
(
    id bigint auto_increment
        primary key,
    type smallint not null,
    relay varchar(256) not null,
    timestamp bigint not null,
    source varchar(256) null,
    source_device int null,
    destination varchar(256) not null,
    destination_device int not null,
    message blob null,
    content blob null,
    guid varchar(36) null,
    server_timestamp bigint null,
    source_uuid varchar(36) null,
    ctime datetime default CURRENT_TIMESTAMP null
);

create index destination_and_type_index
    on t_messages (destination, destination_device, type);

create index destination_index
    on t_messages (destination, destination_device);

