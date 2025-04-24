CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE PAYMENT_STATUS AS ENUM ('pending', 'paid', 'canceled');

CREATE TABLE
  IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    account_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    status PAYMENT_STATUS DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT now (),
    updated_at TIMESTAMP
  );