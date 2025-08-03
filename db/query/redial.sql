-- name: CreateRedial :one
INSERT INTO "Redial" (id, user_id, date, time)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRedial :one
SELECT * FROM "Redial"
WHERE id = $1 LIMIT 1;

-- name: UpdateRedial :one
UPDATE "Redial"
SET 
  user_id = COALESCE(sqlc.narg(user_id), user_id),
  date = COALESCE(sqlc.narg(date), date),
  time = COALESCE(sqlc.narg(time), time)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteRedial :exec
DELETE FROM "Redial"
WHERE id = sqlc.arg(id);