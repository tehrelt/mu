drop if exists index house_connected_service_idx on connected_services (house_id, service_id);
DROP TABLE IF EXISTS connected_services;
DROP TABLE IF EXISTS house;
DROP EXTENSION IF EXISTS "uuid-ossp";