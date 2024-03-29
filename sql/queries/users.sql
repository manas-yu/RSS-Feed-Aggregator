-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, api_keys)
VALUES (
        $1,
        $2,
        $3,
        $4,
        encode(sha256(random()::text::bytea), 'hex')
    )
RETURNING *;
-- name: GetUserByApiKey :one
SELECT *
FROM users
WHERE api_keys = $1;
-- name: GetUserByName :one
SELECT *
FROM users
WHERE name = $1;