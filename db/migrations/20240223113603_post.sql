-- +goose Up
-- +goose StatementBegin

-- Index
CREATE INDEX IF NOT EXISTS posts_title_idx ON posts USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS posts_content_idx ON posts USING GIN (to_tsvector('simple', contet));

CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    title text NOT NULL,
    content text NOT NULL,
    CONSTRAINT fk_user 
        FOREIGN KEY (user_id) 
            REFERENCES users (id)
            ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
DROP INDEX IF EXISTS posts_title_idx;
DROP INDEX IF EXISTS posts_content_idx;
-- +goose StatementEnd
