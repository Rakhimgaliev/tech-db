CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS forum;
DROP TABLE IF EXISTS thread;
DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS vote;

-- primary key = not null + unique
-- reference -> post(id), if id=primaryKey for post then
-- you can write reference -> post


CREATE TABLE "user" (
    nickname    citext, -- Имя пользователя (уникальное поле).
                        -- Данное поле допускает только латиницу, цифры и знак подчеркивания.
                        -- Сравнение имени регистронезависимо.
    fullname    text    not null,   -- Полное имя пользователя.
    about       text,   -- Описание пользователя.
    email       citext  primary key,    -- Почтовый адрес пользователя (уникальное поле).
);

CREATE TABLE forum (
    title   text    not null,   -- Название форума.
    userNickname    citext  references  "user", -- Nickname пользователя, который отвечает за форум.
    slug    text    not null,   -- Человекопонятный URL (https://ru.wikipedia.org/wiki/%D0%A1%D0%B5%D0%BC%D0%B0%D0%BD%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9_URL), уникальное поле.
    posts   bigint, -- Общее кол-во сообщений в данном форуме.
    threads bigint  -- Общее кол-во ветвей обсуждения в данном форуме.
);

CREATE TABLE thread (
    id      integer primary key,    -- Идентификатор ветки обсуждения.
    title   text    not null,   -- Заголовок ветки обсуждения.
    author  citext  references "user",   -- Пользователь, создавший данную тему.
    forum   citext  references forum,   -- Форум, в котором расположена данная ветка обсуждения.
    message text    not null,   -- Описание ветки обсуждения.
    votes   integer,    -- Кол-во голосов непосредственно за данное сообщение форума.
    slug    citext  unique, -- Человекопонятный URL (https://ru.wikipedia.org/wiki/%D0%A1%D0%B5%D0%BC%D0%B0%D0%BD%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9_URL).
                                -- В данной структуре slug опционален и не может быть число
    created timestamp with time zone default now()  not null -- Дата создания ветки на форуме.
);

CREATE TABLE POST (
    id      bigserial   primary key,
    parent  bigint  references post,
    author  citext  references "user",
    message text    not null,
    edited  boolean default false,
    forum   citext  references forum,
    thread  integer references thread,
    created timestamp with time default now()
);

CREATE TABLE VOTE (
    nickname    citext  references "user",
    voice   integer not null
);
