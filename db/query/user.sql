-- name: CreateUser :one
INSERT INTO "User" (id, email, name, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT * FROM "User"
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE "User"
SET 
  email = COALESCE(sqlc.narg(email), email),
  name = COALESCE(sqlc.narg(name), name),
  role = COALESCE(sqlc.narg(role), role)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "User"
WHERE id = sqlc.arg(id);