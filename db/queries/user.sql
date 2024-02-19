-- name: CreateUser :one
INSERT INTO users (
    email,
    username,
    password_hash,
    role
) VALUES (
    $1, $2, $3, $4
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
LIMIT COALESCE(sqlc.narg(limit_arg)::int, 10) 
OFFSET $2;

-- name: GetUsersLikeUsername :many
SELECT * FROM users
WHERE username ILIKE $1
ORDER BY id DESC
LIMIT COALESCE(sqlc.narg(limit_arg)::int, 10) 
OFFSET $2;

-- name: UpdateUserEmail :one
UPDATE users
SET email = $1
WHERE id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = $1
WHERE id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserRole :one
UPDATE users
SET role = $1
WHERE id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserUsername :one
UPDATE users
SET username = $1
WHERE id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;
