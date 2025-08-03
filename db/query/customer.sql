-- name: CreateCustomer :one
INSERT INTO "Customer" (id, book_id, category_id, name, corporation, address, leader, pic, memo)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: SearchCustomer :many
SELECT * FROM "Customer"
WHERE book_id = COALESCE(sqlc.narg(book_id), book_id)
AND name ILIKE '%' || COALESCE(sqlc.narg(name), name) || '%'
AND corporation ILIKE '%' || COALESCE(sqlc.narg(corporation), corporation) || '%'
AND address ILIKE '%' || COALESCE(sqlc.narg(address), address) || '%'
AND leader = COALESCE(sqlc.narg(leader), leader)
AND pic = COALESCE(sqlc.narg(pic), pic) 
AND memo ILIKE '%' || COALESCE(sqlc.narg(memo), memo) || '%';

-- name: UpdateCustomer :one
UPDATE "Customer"
SET 
  name = COALESCE(sqlc.narg(name), name),
  book_id = COALESCE(sqlc.narg(book_id), book_id),
  corporation = COALESCE(sqlc.narg(corporation), corporation),
  address = COALESCE(sqlc.narg(address), address),
  memo = COALESCE(sqlc.narg(memo), memo)
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: DeleteCustomer :exec
DELETE FROM "Customer"
WHERE id = sqlc.arg(id);