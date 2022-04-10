create table if not exists image
(
    id           bigint primary key,
    name         text      not null,
    size         bigint    not null,
    description  text,
    s3_link      text      not null,
    label        text,
    category     text,
    dateUploaded timestamp not null,
    user_id      bigint,
    constraint image_user_id_fk foreign key (user_id) references "user" (id)
);

create sequence if not exists image_seq increment 1 start 1;
