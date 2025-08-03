-- name: CreateUser :one
INSERT INTO "User" (id, email, name, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT * FROM "User"
WHERE id = $1 LIMIT 1;