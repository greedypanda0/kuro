package database

const schema = `
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- repositories

CREATE TABLE IF NOT EXISTS repositories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    author TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
`
