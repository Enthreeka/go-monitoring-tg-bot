create table if not exists channel(
    id int generated always as identity,
    tg_id bigint unique not null,
    channel_name varchar(150) null,
    channel_url varchar(150) null,
    channel_status chan_status not null,
    primary key (id)
);