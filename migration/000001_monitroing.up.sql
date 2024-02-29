set timezone = 'Europe/Moscow';

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
            CREATE TYPE role AS ENUM ('user', 'admin','superAdmin');
        END IF;
    END $$;

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'req_status') THEN
            CREATE TYPE req_status AS ENUM ('in progress','approved','rejected');
        END IF;
    END $$;

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'chan_status') THEN
            CREATE TYPE chan_status AS ENUM ('kicked','administrator','left','member');
        END IF;
    END $$;

create table if not exists "user"(
    id           bigint unique,
    tg_username  text not null ,
    created_at   timestamp not null,
    phone        varchar(20) null,
    channel_from varchar(150) null,
    user_role         role default 'user' not null,
    primary key (id)
);

create table if not exists channel(
    id int generated always as identity,
    tg_id bigint unique not null,
    channel_name varchar(150) null,
    channel_url varchar(150) null,
    channel_status chan_status not null,
    primary key (id)
);

create table if not exists user_channel(
    user_id bigint,
    channel_tg_id bigint,
    primary key (user_id,channel_tg_id),
    foreign key (user_id)
        references "user" (id) on delete cascade,
    foreign key (channel_tg_id)
        references channel (tg_id) on delete cascade
);

create table if not exists notification(
    id int generated always as identity,
    channel_id bigint,
    notification_text text null ,
    file_id varchar(150) null,
    file_type varchar(150) null,
    button_url varchar(150) null,
    button_text varchar(150) null,
    primary key (id),
    foreign key (channel_id)
    references channel (tg_id) on delete cascade
);

create table if not exists request(
    id int generated always as identity,
    channel_tg_id bigint not null,
    user_id bigint not null,
    status_request req_status default 'in progress' not null,
    date_request timestamp not null,
    foreign key (user_id)
    references "user" (id) on delete cascade,
    foreign key (channel_tg_id)
        references channel (tg_id) on delete cascade
);

create table if not exists sender(
    id int generated always as identity,
    channel_tg_id bigint not null,
    message text not null,
    primary key (id),
    foreign key (channel_tg_id)
        references channel (tg_id) on delete cascade
);

create table if not exists spam_bot(
    id int generated always as identity,
    token varchar(255),
    bot_name varchar(150) null,
    primary key (id)
);