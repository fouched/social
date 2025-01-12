create table roles (
    id serial primary key,
    name varchar(255) not null unique,
    level int not null default 0,
    description text
);

insert into roles (name, description, level)
values ('user', 'create posts and comments', 1);

insert into roles (name, description, level)
values ('moderator', 'update other users posts', 2);

insert into roles (name, description, level)
values ('admin', 'update and delete other users posts', 3);