-- name: CreateCategory :one
INSERT INTO "Category" (id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetCategory :one
SELECT * FROM "Category"
WHERE id = $1 LIMIT 1;