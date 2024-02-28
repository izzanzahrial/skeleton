-- name: CreateUser :one
INSERT INTO users (
    email,
    username,
    password_hash,
    role,
    origin
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE id = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users 
WHERE id = $1 LIMIT 1 
FOR UPDATE;

-- name: GetuserByEmailOrUsername :one
SELECT * FROM users 
WHERE (email = $1 OR $1 = '')
AND (username = $2 OR $2 = '')
AND deleted_at IS NULL
LIMIT 1;

-- name: GetUsersByRole :many
SELECT * FROM users
WHERE role = $1 AND deleted_at IS NULL
ORDER BY id DESC
LIMIT COALESCE(sqlc.narg(limit_param)::int, 10) 
OFFSET $2;

-- name: GetUsersLikeUsername :many
SELECT * FROM users
WHERE username ILIKE $1
ORDER BY id DESC
LIMIT COALESCE(sqlc.narg(limit_param)::int, 10) 
OFFSET $2;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = $1, updated_at = NOW()
WHERE id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserRole :one
UPDATE users
SET role = $1, updated_at = NOW()
WHERE id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET updated_at = NOW(), email = $1, username = $2, password_hash = $3
WHERE id = $4 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateUserGoogle :one
INSERT INTO users (
    email,
    first_name,
    last_name,
    picture_url,
    refresh_token,
    role,
    origin
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetuserByEmail :one
SELECT * FROM users 
WHERE (email = $1 OR $1 = '')
AND deleted_at IS NULL
LIMIT 1;
