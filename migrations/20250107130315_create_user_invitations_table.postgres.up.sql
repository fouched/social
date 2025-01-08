create table user_invitations (
    token bytea not null,
    user_id bigint not null,
    expiry timestamp with time zone not null,

    primary key (token, user_id)
);