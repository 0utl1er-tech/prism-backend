-- name: CreateRedial :one
INSERT INTO "Redial" (id, user_id, date, time)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRedial :one
SELECT * FROM "Redial"
WHERE id = $1 LIMIT 1;