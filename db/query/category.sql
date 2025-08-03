-- name: CreateCategory :one
INSERT INTO "Category" (id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetCategory :one
SELECT * FROM "Category"
WHERE id = $1 LIMIT 1;

-- name: UpdateCategory :one
UPDATE "Category"
SET 
  name = COALESCE(sqlc.narg(name), name)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM "Category"
WHERE id = sqlc.arg(id);