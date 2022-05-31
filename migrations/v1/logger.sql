CREATE TABLE public.signup_log (
	id bigserial NOT NULL,
	company_uuid varchar(64) NOT NULL,
	operator_uuid varchar(64) NOT NULL,
	operator_name varchar(64) NOT NULL,
	way varchar(16) NOT NULL,
	platform varchar(16) NOT NULL,
	ip inet NOT NULL,
	extra jsonb NOT NULL,
	created_time timestamp(0) NOT NULL,
	CONSTRAINT signup_log_pkey PRIMARY KEY (id)
);

COMMENT on table public.signup_log is '登录日志表';
COMMENT on column public.signup_log.company_uuid is '企业uuid';
COMMENT on column public.signup_log.operator_uuid is '操作人uuid';
COMMENT on column public.signup_log.operator_name is '操作人名称';
COMMENT on column public.signup_log.way is '登录方式';
COMMENT on column public.signup_log.ip is 'ip地址';
COMMENT on column public.signup_log.platform is '平台';
COMMENT on column public.signup_log.extra is '其他内容';
COMMENT on column public.signup_log.created_time is '创建时间';