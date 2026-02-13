package database

const schema = `
-- repositories

CREATE TABLE IF NOT EXISTS repositories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    author TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
`
