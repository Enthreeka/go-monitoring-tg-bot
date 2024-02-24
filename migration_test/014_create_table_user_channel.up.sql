create table if not exists user_channel(
    user_id bigint,
    channel_tg_id bigint,
    primary key (user_id,channel_tg_id),
    foreign key (user_id)
    references "user" (id) on delete cascade,
    foreign key (channel_tg_id)
    references channel (tg_id) on delete cascade
);