-- name: CreateStaff :one
INSERT INTO "Staff" (id, name, sex)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetStaff :one
SELECT * FROM "Staff"
WHERE id = $1 LIMIT 1;