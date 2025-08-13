CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hashed TEXT NOT NULL,
    subscription TEXT NOT NULL DEFAULT 'nosubs'
        CHECK (subscription IN ('nosubs','month','year')),
    register_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    is_paid  BOOLEAN NOT NULL DEFAULT FALSE
);
