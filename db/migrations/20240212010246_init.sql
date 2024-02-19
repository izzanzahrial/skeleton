-- +goose Up
-- +goose StatementBegin

-- Extention 
-- For storing text data with case insensitive
CREATE EXTENSION IF NOT EXISTS citext;

-- Enum
CREATE TYPE roles AS ENUM (
    'admin',
    'user'
);

-- DB tables
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    -- email by standard is case insensitive
    email citext UNIQUE NOT NULL,
    username text UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    role roles NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
