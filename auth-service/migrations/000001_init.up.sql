CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE ROLE AS ENUM ('admin', 'regular');

CREATE TABLE IF NOT EXISTS credentials (
  id UUID PRIMARY KEY,
  hashed_password VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS roles (
  user_id UUID REFERENCES credentials(id) ON DELETE CASCADE,
  role ROLE NOT NULL DEFAULT 'regular',
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

create unique index roles_user_id_role_unq on roles (user_id, role);