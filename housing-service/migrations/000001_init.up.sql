CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS house (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  address VARCHAR NOT NULL,
  rooms_qty INT NOT NULL,
  residents_qty INT NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS connected_services (
  house_id UUID PRIMARY KEY REFERENCES house(id) NOT NULL,
  service_id UUID NOT NULL,
  connected_at TIMESTAMP DEFAULT now()
);

create unique index house_connected_service_idx on connected_services (house_id, service_id);