CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(100) UNIQUE                            NOT NULL,
    password_hash VARCHAR(255)                                   NOT NULL
);

CREATE TABLE posts (
    id            SERIAL PRIMARY KEY,
    user_id       SERIAL REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    title         VARCHAR(255) NOT NULL,
    body          TEXT 
);

CREATE TABLE comments (
    id            SERIAL PRIMARY KEY,
    user_id       SERIAL REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    post_id       SERIAL REFERENCES posts (id) ON DELETE CASCADE NOT NULL,
    email         VARCHAR(100)                                   NOT NULL,
    body          TEXT
);