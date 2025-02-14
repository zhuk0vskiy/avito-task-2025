create extension if not exists pgcrypto;

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

create table if not exists bought_merchs
(
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id),
    merch_id uuid not null references merchs(id),
    amount smallint
);

create table if not exists transactions
(
    id uuid primary key default gen_random_uuid(),
    from_user_id uuid not null references users(id),
    to_user_id uuid not null references users(id),
    coins_amount smallint not null
);

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

