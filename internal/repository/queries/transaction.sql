-- name: CreateTransaction :one
INSERT INTO transactions ("account_id", "category_id", "invoice_id", "planned_income_id", "planned_expense_id", "amount", "transaction_type", "description", "transaction_date")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transactions WHERE "id" = $1;

-- name: GetTransactionByAccountID :many
SELECT * FROM transactions WHERE "account_id" = $1;

-- name: GetTransactionByCategoryID :many
SELECT * FROM transactions WHERE "category_id" = $1;

-- name: GetTransactionByInvoiceID :many
SELECT * FROM transactions WHERE "invoice_id" = $1;

-- name: GetTransactionByPlannedIncomeID :many
SELECT * FROM transactions WHERE "planned_income_id" = $1;

-- name: GetTransactionByPlannedExpenseID :many
SELECT * FROM transactions WHERE "planned_expense_id" = $1;

-- name: UpdateTransaction :one
UPDATE transactions SET "amount" = $2, "transaction_type" = $3, "description" = $4, "transaction_date" = $5
WHERE "id" = $1
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions WHERE "id" = $1;