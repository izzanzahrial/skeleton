-- name: GetPostsFullText :many
SELECT * FROM posts
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', sqlc.arg(title)::text) OR sqlc.arg(title)::text = '')
AND (to_tsvector('simple', content) @@ plainto_tsquery('simple', sqlc.arg(content)::text) OR sqlc.arg(content)::text = '')
LIMIT COALESCE(sqlc.narg(limit_param)::int, 10) 
OFFSET $1;

-- name: CreatePost :one
INSERT INTO posts (
    user_id,
    title,
    content
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetPostByUserID :many
SELECT * FROM posts 
WHERE user_id = $1;