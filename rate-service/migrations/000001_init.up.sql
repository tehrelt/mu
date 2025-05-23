CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE SERVICE_TYPE AS ENUM (
  'water_supply',
  'heating',
  'power_supply',
  'gas_supply'
);

CREATE TABLE
  IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    s_name VARCHAR(255) NOT NULL,
    measure_unit VARCHAR(255) NOT NULL,
    s_type SERVICE_TYPE NOT NULL,
    rate BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now (),
    updated_at TIMESTAMP
  );