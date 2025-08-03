-- name: CreateContact :one
INSERT INTO "Contact" (id, customer_id, staff_id, phone, mail, fax)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetContact :one
SELECT * FROM "Contact"
WHERE id = $1 LIMIT 1;

-- name: UpdateContact :one
UPDATE "Contact"
SET 
  phone = COALESCE(sqlc.narg(phone), phone),
  mail = COALESCE(sqlc.narg(mail), mail),
  fax = COALESCE(sqlc.narg(fax), fax)
WHERE 
  id = sqlc.arg(id)
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM "Contact"
WHERE id = sqlc.arg(id);