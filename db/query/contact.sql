-- name: CreateContact :one
INSERT INTO "Contact" (id, customer_id, staff_id, phone, mail, fax)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetContact :one
SELECT * FROM "Contact"
WHERE id = $1 LIMIT 1;