-- name: CreateInvoice :one
INSERT INTO invoices ("card_id", "reference_date", "due_date", "total_amount", "status", "paid_at")
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetInvoiceByID :one
SELECT * FROM invoices WHERE "id" = $1;

-- name: GetInvoiceByCardID :many
SELECT * FROM invoices WHERE "card_id" = $1;

-- name: UpdateInvoice :one
UPDATE invoices SET "total_amount" = $2, "status" = $3, "paid_at" = $4
WHERE "id" = $1
RETURNING *;

-- name: DeleteInvoice :exec
DELETE FROM invoices WHERE "id" = $1;