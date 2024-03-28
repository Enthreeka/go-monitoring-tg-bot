create table if not exists sender(
    id int generated always as identity,
    channel_tg_id bigint not null,
    message text not null,
    primary key (id),
    foreign key (channel_tg_id)
    references channel (tg_id) on delete cascade
);