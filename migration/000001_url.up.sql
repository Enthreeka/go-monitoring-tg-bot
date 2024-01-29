sql create table if not exists "user"(
                                         id bigint unique,
                                         tg_username text ,
                                         created_at timestamp,
                                         phone int,
                                         channel_from varchar(150),
    role role default 'user',
    primary key (id)
    );

create table if not exists notification(
                                           id int generated always as identity,
                                           notification_text text ,
                                           file_id varchar(150),
    file_type varchar(150),
    button_url varchar(150),
    channel_id bigint,
    primary key (id),
    foreign key (channel_id)
    references channel (id) on delete cascade
    );


create table if not exists channel(
                                      id int generated always as identity,
                                      tg_id bigint,
                                      channel_name varchar(150),
    channel_url varchar(150),
    primary key (id)
    );

create table if not exists request(
                                      id int generated always as identity,
                                      user_id bigint,
                                      status_request bigint, --to-do
                                      foreign key (user_id)
    references "user" (id) on delete cascade
    );
