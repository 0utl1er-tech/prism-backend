-- name: CreateStatus :one
INSERT INTO "Status" (id, name, effective, ng)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetStatus :one
SELECT * FROM "Status"
WHERE id = $1 LIMIT 1;