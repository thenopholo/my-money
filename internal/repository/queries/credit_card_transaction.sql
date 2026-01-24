-- name: CreateCreditCardTransaction :one
INSERT INTO credit_card_transactions (
    "card_id", "category_id", "invoice_id", "amount", "description",
    "installments", "current_installment", "installment_value", "transaction_date"
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetCreditCardTransactionByID :one
SELECT * FROM credit_card_transactions WHERE "id" = $1;

-- name: GetCreditCardTransactionsByCardID :many
SELECT * FROM credit_card_transactions WHERE "card_id" = $1;

-- name: GetCreditCardTransactionsByCategoryID :many
SELECT * FROM credit_card_transactions WHERE "category_id" = $1;

-- name: GetCreditCardTransactionsByInvoiceID :many
SELECT * FROM credit_card_transactions WHERE "invoice_id" = $1;

-- name: GetPendingCreditCardTransactions :many
SELECT * FROM credit_card_transactions
WHERE "card_id" = $1 AND "invoice_id" IS NULL;

-- name: UpdateCreditCardTransaction :one
UPDATE credit_card_transactions
SET "amount" = $2, "description" = $3, "installments" = $4,
    "current_installment" = $5, "installment_value" = $6, "transaction_date" = $7
WHERE "id" = $1
RETURNING *;

-- name: AssignTransactionToInvoice :exec
UPDATE credit_card_transactions
SET "invoice_id" = $2
WHERE "id" = $1;

-- name: DeleteCreditCardTransaction :exec
DELETE FROM credit_card_transactions WHERE "id" = $1;