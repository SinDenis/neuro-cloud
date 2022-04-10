create table if not exists "user"
(
    id       bigint primary key,
    username text not null,
    password text not null,
    email    text not null
);

create sequence if not exists user_seq increment 1 start 1;