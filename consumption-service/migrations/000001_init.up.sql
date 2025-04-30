CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists cabinet (
    id uuid primary key not null default uuid_generate_v4 (),
    account_id uuid not null,
    service_id uuid not null,
    consumed bigint not null default 0,
    created_at timestamp not null default now (),
    updated_at timestamp,
    unique (account_id, service_id)
);

create table if not exists consumption_log (
    id uuid primary key not null default uuid_generate_v4 (),
    amount bigint not null,
    payment_id uuid not null,
    cabinet_id uuid not null references cabinet (id) on delete cascade,
    created_at timestamp not null default now ()
);
