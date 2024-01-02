-- name: CreateUser :one
INSERT INTO users ("firstName", "lastName", email, password)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;