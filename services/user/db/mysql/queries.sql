-- name: CreateUser :execresult
INSERT INTO users (email, password_hash)
VALUES (?, ?);

-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at
FROM users
WHERE email = ?;

-- name: ListUsers :many
SELECT id, email, password_hash, created_at
FROM users
ORDER BY id;