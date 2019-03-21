CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS forum;
DROP TABLE IF EXISTS thread;
DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS vote;

-- primary key = not null + unique

CREATE TABLE "user" (
    nickname    text,   -- Имя пользователя (уникальное поле).
                        -- Данное поле допускает только латиницу, цифры и знак подчеркивания.
                        -- Сравнение имени регистронезависимо.
    fullname    text    not null,   -- Полное имя пользователя.
    about       text,   -- Описание пользователя.
    email       citext  primary key,    -- Почтовый адрес пользователя (уникальное поле).
)

CREATE TABLE forum (
    title   text    not null,   -- Название форума.
    userNickname    citext  references  "user"  not null,   -- Nickname пользователя, который отвечает за форум.
    slug    text    not null,   -- Человекопонятный URL (https://ru.wikipedia.org/wiki/%D0%A1%D0%B5%D0%BC%D0%B0%D0%BD%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9_URL), уникальное поле.
    posts   bigint, -- Общее кол-во сообщений в данном форуме.
    threads bigint  -- Общее кол-во ветвей обсуждения в данном форуме.
)

CREATE TABLE thread (
    id      integer primary key,    -- Идентификатор ветки обсуждения.
    title   text    not null,   -- Заголовок ветки обсуждения.
    author  citext  references "user"   not null,   -- Пользователь, создавший данную тему.
    forum   citext  references forum not null,   -- Форум, в котором расположена данная ветка обсуждения.
    message text    not null,   -- Описание ветки обсуждения.
    votes   integer,    -- Кол-во голосов непосредственно за данное сообщение форума.
    slug    citext  unique, -- Человекопонятный URL (https://ru.wikipedia.org/wiki/%D0%A1%D0%B5%D0%BC%D0%B0%D0%BD%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9_URL).
                                -- В данной структуре slug опционален и не может быть число
    created timestamp with time zone default now()  not null -- Дата создания ветки на форуме.
)

CREATE TABLE POST (

)

CREATE TABLE VOTE (
    
)
