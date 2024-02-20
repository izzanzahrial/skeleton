-- +goose Up
-- +goose StatementBegin
CREATE TYPE origins AS ENUM (
    'native',
    'google'
);

ALTER TABLE users ADD COLUMN first_name VARCHAR(255);
ALTER TABLE users ADD COLUMN last_name VARCHAR(255);
ALTER TABLE users ADD COLUMN picture_url text;
ALTER TABLE users ADD COLUMN refresh_token text;
ALTER TABLE users ADD COLUMN origin origins NOT NULL;
ALTER TABLE users ALTER COLUMN username DROP NOT NULL;
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
