-- name: CreateBook :one
INSERT INTO "Book" (id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetBook :one
SELECT * FROM "Book"
WHERE id = $1 LIMIT 1;

-- name: UpdateBook :one
UPDATE "Book"
SET 
  name = COALESCE(sqlc.narg(name), name)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM "Book"
WHERE id = sqlc.arg(id);