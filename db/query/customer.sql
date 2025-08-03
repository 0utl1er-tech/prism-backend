-- name: CreateCustomer :one
INSERT INTO "Customer" (id, book_id, category_id, name, corporation, phone, prefecture, address, leader, leader_sex,)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;