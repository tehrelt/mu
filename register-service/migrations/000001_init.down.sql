DROP TABLE IF EXISTS payments;

DROP TYPE PAYMENT_STATUS AS ENUM ('pending', 'paid', 'canceled');

DROP EXTENSION IF EXISTS "uuid-ossp";