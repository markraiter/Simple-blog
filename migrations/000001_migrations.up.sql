CREATE TABLE users
(
    id            serial       not null unique,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE posts
(
    id          serial       not null unique,
    user_id     int          not null references users (id) on delete cascade,
    title       varchar(255) not null,
    body        text
);

CREATE TABLE comments
(
    id      serial                                           not null unique,
    user_id int references users (id) on delete cascade      not null,
    post_id int references posts (id) on delete cascade      not null,
    email   varchar(255)                                     not null unique,
    body    text
);