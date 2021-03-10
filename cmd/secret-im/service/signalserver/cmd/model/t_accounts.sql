create table t_accounts
(
    id int auto_increment,
    number varchar(20) not null comment 'phonenumber',
    data text not null comment 'account json byte',
    PRIMARY KEY (`id`),
    UNIQUE KEY `number_index` (`number`)
)
    comment '通过验证的真实客户';

