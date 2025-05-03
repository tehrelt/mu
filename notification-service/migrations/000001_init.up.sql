CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS integrations (
    user_id UUID PRIMARY KEY NOT NULL,
    telegram_chat_id TEXT
);
