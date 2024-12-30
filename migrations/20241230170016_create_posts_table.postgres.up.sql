create table posts
(
    id bigserial primary key,
    user_id bigint references users(id) not null,
    title text not null,
    content text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
