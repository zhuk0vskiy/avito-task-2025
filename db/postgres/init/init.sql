create extension if not exists pgcrypto;
alter system set max_connections = 1000;

create user admin with password 'avito';
grant all privileges on database shop to admin;
alter default privileges in schema public grant all privileges on tables to admin;

-- create user test_admin with password 'test_avito';
-- create database test_shop;
-- grant all privileges on database test_shop to test_admin;

-- create role readaccess;
-- grant connect on database shop to readaccess;
-- grant select on all tables in schema public to readaccess;
-- create user admin_ro with password 'avito_ro';
-- grant readaccess to admin_ro;
-- alter default privileges in schema public grant select on tables to readaccess;

create table if not exists merchs
(
    id uuid primary key default gen_random_uuid(),
    type text not null,
    cost integer not null
);

create table if not exists users
(
    id uuid primary key default gen_random_uuid(),
    username text not null,
    password bytea,
    coins_amount integer not null
);

create index idx_username on users using hash (username);

create table if not exists bought_merchs
(
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id),
    merch_id uuid not null references merchs(id),
    amount smallint
);

create index idx_bought_merch_amount on bought_merchs using hash (amount);

create table if not exists transactions
(
    id uuid primary key default gen_random_uuid(),
    from_user_id uuid not null references users(id),
    to_user_id uuid not null references users(id),
    coins_amount smallint not null
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_bought_merchs_user_merch 
ON bought_merchs(user_id, merch_id);

insert into merchs(type, cost) 
values 
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 500);

