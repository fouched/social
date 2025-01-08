create extension if not exists citext;

create table users (
    id bigserial primary key,
    email citext unique not null,
    username varchar(255) unique not null,
    password bytea not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);