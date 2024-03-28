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