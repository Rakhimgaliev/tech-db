CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS forum;
DROP TABLE IF EXISTS thread;
DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS vote;

CREATE TABLE "user" (
    nickname    citext  primary key COLLATE "POSIX",
    fullname    text    not null,
    about       text,
    email       citext  unique  not null
);

CREATE TABLE forum (
    title   text    not null,
    userNickname    citext  references "user"   not null,
    slug    citext    primary key,
    postCount   bigint  default 0 not null,
    threadCount bigint  default 0 not null
);

CREATE TABLE thread (
    id      bigserial   primary key,
    title   text    not null,
    userNickname    citext  references "user"   not null,
    forum   citext  references forum    not null,
    message text    not null,
    votes   integer default 0   not null,
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
    created timestamp with time zone default now(),
    children    integer[]
);

CREATE TABLE vote (
    nickname    citext  references "user",
    voice   boolean not null,
    threadId integer references thread,
    CONSTRAINT uniqueVote UNIQUE (nickname, threadId)
);


CREATE TABLE forum_user (
    nickname    citext  references "user",
    forum   citext  references forum,
    CONSTRAINT  uniqueForumUser UNIQUE (nickname, forum)
);

CREATE OR REPLACE FUNCTION create_children() RETURNS trigger as $create_children$
BEGIN
   IF NEW.parent IS NULL THEN
     NEW.children := (ARRAY [NEW.id]);
     return NEW;
   end if;

   NEW.children := (SELECT array_append(p.children, NEW.id::integer)
                from post p where p.id = NEW.parent);
  RETURN NEW;
END;
$create_children$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS create_children ON post;

CREATE TRIGGER create_children BEFORE INSERT ON post
  FOR EACH ROW EXECUTE PROCEDURE create_children();