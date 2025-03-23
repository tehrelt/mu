CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  last_name VARCHAR,
  first_name VARCHAR,
  middle_name VARCHAR,
  email VARCHAR UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS personal_data (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    passport_number INTEGER NOT NULL,
    passport_series INTEGER NOT NULL,
    phone CHAR(11) NOT NULL,
    snils VARCHAR NOT NULL
);