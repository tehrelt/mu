CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
  IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    house_id UUID NOT NULL,
    balance BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now (),
    updated_at TIMESTAMP
  );