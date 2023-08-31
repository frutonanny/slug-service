-- +goose Up
-- Таблица сегментов.
create table if not exists slugs
(
    id         bigserial primary key,
    "name"     varchar(255) not null unique,
    options    jsonb,
    created_at timestamptz  not null default now(),
    deleted_at timestamptz
);

-- Таблица пользователей.
create table if not exists users
(
    id         bigserial primary key,
    user_id    uuid unique,
    created_at timestamptz not null default now()
);

-- Таблица пользователей и их сегментов.
create table if not exists users_slugs
(
    id        bigserial primary key,
    user_id   uuid         not null references users (user_id),
    slug_id   bigint       not null references slugs (id),
    -- Сознательно дублируем данные, чтобы избежать join-а на "slugs" при получении сегментов пользователя.
    slug_name varchar(255) not null,
    -- Время, до которого пользователь находится в сегменте.
    slug_ttl  timestamptz
);

-- Таблица истории добавления и удаления пользователей в сегменты.
create table if not exists events
(
    id         bigserial primary key,
    user_id    uuid        not null,
    slug_id    bigint      not null,
    event      text        not null,
    created_at timestamptz not null default now()
);

create index events_user_id_idx on events (user_id);
create index events_created_at_idx on events (created_at);

-- Таблица хранения задач для дальнейшего выполнение.
create table if not exists outbox
(
    id             bigserial primary key,
    name           text        not null,
    data           text,
    reserved_until timestamp   not null default now(),
    created_at     timestamptz not null default now()
);

-- +goose Down
drop table outbox;
drop table events;
drop table users_slugs;
drop table users;
drop table slugs;
