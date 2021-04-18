-- citext делает строки нечувствительными к регистру
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS Users
(
    id          BIGSERIAL,
    fullname    VARCHAR(100) NOT NULL,
    email       citext       NOT NULL UNIQUE
    nickname    citext       NOT NULL PRIMARY KEY,
    about       TEXT,
    
);

CREATE TABLE IF NOT EXISTS Forum
(
    title           TEXT   NOT NULL,
    user_nickname   citext NOT NULL REFERENCES users (nickname),
    slug            citext NOT NULL PRIMARY KEY,
    posts           BIGINT,
    threads         INT
);


CREATE TABLE IF NOT EXISTS Threads
(
    id          BIGSERIAL NOT NULL PRIMARY KEY,
    title       TEXT      NOT NULL,
    author      citext    NOT NULL REFERENCES users (nickname),
    forum       citext    NOT NULL REFERENCES forum (slug),
    msg         TEXT      NOT NULL,
    votes       INTEGER   DEFAULT 0,
    slug        citext,
    created     TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS Posts
(
    id          BIGSERIAL NOT NULL CONSTRAINT posts_pkey PRIMARY KEY,
    author      citext    NOT NULL CONSTRAINT posts_author_fkey REFERENCES users,
    forum       citext    NOT NULL CONSTRAINT posts_forum_fkey REFERENCES forum,
    thread      INTEGER   NOT NULL CONSTRAINT posts_thread_fkey REFERENCES threads,
    msg         TEXT      NOT NULL,
    parent      INTEGER   DEFAULT 0,
    is_edited   BOOLEAN   DEFAULT false,
    path        BIGINT[]  DEFAULT ARRAY []::INTEGER[]
    created     TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
);


CREATE TABLE IF NOT EXISTS Votes
(
    nickname  citext NOT NULL REFERENCES users (nickname),
    thread_id INT    NOT NULL REFERENCES threads (id),
    voice     INT    NOT NULL,
    UNIQUE (nickname, thread_id)
);

CREATE TABLE IF NOT EXISTS Users_to_forums
(
    nickname citext NOT NULL REFERENCES users (nickname),
    forum    citext NOT NULL REFERENCES forum (slug),
    UNIQUE (nickname, forum)
);