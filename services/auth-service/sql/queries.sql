-- name: GetUserByEmail :one
SELECT id, email, password, created_at, updated_at
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES ($1, $2)
RETURNING id, email, created_at, updated_at;