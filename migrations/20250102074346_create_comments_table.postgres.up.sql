create table comments (
    id bigserial primary key,
    post_id bigint references posts(id) not null,
    user_id bigint references users(id) not null,
    content text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
)