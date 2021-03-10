
create table t_pend_accounts
(
	id bigint auto_increment  comment '自增id',
	number varchar(30) not null comment 'phone number',
	verification_code varchar(16) not null,
	create_time timestamp default CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`),
	UNIQUE KEY `number_index` (`number`)
)
comment ' 验证码表';

