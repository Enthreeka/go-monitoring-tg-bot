create table if not exists "user"(
    id           bigint unique,
    tg_username  text not null ,
    created_at   timestamp not null,
    phone        varchar(20) null,
    channel_from varchar(150) null,
    user_role         role default 'user' not null,
    primary key (id)
);