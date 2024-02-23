-- +goose Up
-- +goose StatementBegin

-- Extension
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Index
-- using gin index and gin_trgm_ops for better searching performance
-- for ILIKE queries, example in 'GetUsersLikeUsername'
-- reference: https://niallburkley.com/blog/index-columns-for-like-in-postgres/
CREATE INDEX trgm_idx_users_username ON users USING gin (username gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS pg_trgm;
DROP INDEX IF EXISTS trgm_idx_users_username;
-- +goose StatementEnd
