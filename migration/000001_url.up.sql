DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
            CREATE TYPE role AS ENUM ('user', 'admin');
        END IF;
    END $$;

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
            CREATE TYPE status AS ENUM ('in progress','Approved','Rejected');
        END IF;
    END $$;

create table if not exists "user"(
    id           bigint unique,
    tg_username  text,
    created_at   timestamp,
    phone        varchar(20),
    channel_from varchar(150),
    user_role         role default 'user',
    primary key (id)
);

create table if not exists channel(
    id int generated always as identity,
    tg_id bigint,
    channel_name varchar(150),
    channel_url varchar(150),
    primary key (id)
);

create table if not exists notification(
    id int generated always as identity,
    channel_id bigint,
    notification_text text ,
    file_id varchar(150),
    file_type varchar(150),
    button_url varchar(150),
    primary key (id),
    foreign key (channel_id)
    references channel (id) on delete cascade
);


create table if not exists request(
    id int generated always as identity,
    user_id bigint,
    status_request status default 'in progress',
    foreign key (user_id)
    references "user" (id) on delete cascade
);
