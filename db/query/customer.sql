-- name: CreateCustomer :one
INSERT INTO "Customer" (id, book_id, category_id, name, corporation, address, leader, pic, memo)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;