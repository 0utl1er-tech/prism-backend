-- name: CreateBook :one
INSERT INTO "Book" (id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetBook :one
SELECT * FROM "Book"
WHERE id = $1 LIMIT 1;