create table posts
(
    id bigserial primary key,
    user_id bigint references users(id) not null,
    title varchar(128) not null,
    content text not null,
    tags varchar(128) not null default '',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
