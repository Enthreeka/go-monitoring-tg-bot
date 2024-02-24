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