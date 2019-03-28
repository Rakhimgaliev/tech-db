CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS forum;
DROP TABLE IF EXISTS thread;
DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS vote;

CREATE TABLE "user" (
    nickname    citext,
    fullname    text    not null,
    about       text,
    email       citext  primary key
);

CREATE TABLE forum (
    title   text    not null,
    userNickname    citext  references "user",
    slug    text    primary key,
    postCount   bigint  default 0 not null,
    threadCount bigint  default 0 not null
);

CREATE TABLE thread (
    id      integer primary key,
    title   text    not null,
    userNickname    citext  references "user",
    forum   citext  references forum,
    message text    not null,
    votes   integer,
    slug    citext  unique,
    created timestamp with time zone default now()  not null
);

CREATE TABLE post (
    id      bigserial   primary key,
    parent  bigint  references post,
    userNickname    citext  references "user",
    message text    not null,
    edited  boolean default false,
    forum   citext  references forum,
    thread  integer references thread,
    created timestamp with time zone default now()
);

CREATE TABLE vote (
    nickname    citext  references "user",
    voice   integer not null
);


CREATE TABLE forum_user (
    nickname    citext  references "user",
    forum   citext  references forum,
    CONSTRAINT  uniqueForumUser UNIQUE (nickname, forum)
);

