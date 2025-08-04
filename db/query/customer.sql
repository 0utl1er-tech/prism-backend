-- name: CreateCustomer :one
INSERT INTO "Customer" (id, book_id, category_id, name, corporation, address, leader, pic, memo)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetCustomer :one
SELECT 
    c.id as customer_id,
    c.book_id as customer_book_id,
    c.category_id as customer_category_id,
    c.job as customer_job,
    c.name as customer_name,
    c.corporation as customer_corporation,
    c.address as customer_address,
    c.leader as customer_leader,
    c.pic as customer_pic,
    c.memo as customer_memo,
    c.created_at as customer_created_at,
    ct.id as contact_id,
    ct.customer_id as contact_customer_id,
    ct.staff_id as contact_staff_id,
    ct.phone as contact_phone,
    ct.mail as contact_mail,
    ct.fax as contact_fax,
    ct.created_at as contact_created_at
FROM "Customer" c
LEFT JOIN "Contact" ct ON c.id = ct.id
WHERE c.id = $1;

-- name: GetCustomerByBookId :many
SELECT * FROM "Customer"
WHERE book_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchCustomer :many
SELECT * FROM "Customer"
WHERE book_id = COALESCE(sqlc.narg(book_id), book_id)
AND name ILIKE '%' || COALESCE(sqlc.narg(name), name) || '%'
AND corporation ILIKE '%' || COALESCE(sqlc.narg(corporation), corporation) || '%'
AND address ILIKE '%' || COALESCE(sqlc.narg(address), address) || '%'
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