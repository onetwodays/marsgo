create table t_keys
(
    id int auto_increment,
    number varchar(20) null,
    keyid int null,
    publickey varchar(120) null,
    last_resort int null,
    deviceid int null,
    PRIMARY KEY (`id`)
);
