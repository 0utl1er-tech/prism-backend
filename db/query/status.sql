-- name: CreateStatus :one
INSERT INTO "Status" (id, name, effective, ng)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetStatus :one
SELECT * FROM "Status"
WHERE id = $1 LIMIT 1;

-- name: UpdateStatus :one
UPDATE "Status"
SET 
  name = COALESCE(sqlc.narg(name), name),
  effective = COALESCE(sqlc.narg(effective), effective),
  ng = COALESCE(sqlc.narg(ng), ng)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteStatus :exec
DELETE FROM "Status"
WHERE id = sqlc.arg(id);