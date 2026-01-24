-- name: CreateBankAccount :one
INSERT INTO bank_accounts ("user_id", "name", "bank_name", "account_type", "balance", "is_active")
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetBankAccountByID :one
SELECT * FROM bank_accounts WHERE "id" = $1;

-- name: GetBankAccountByUserID :many
SELECT * FROM bank_accounts WHERE "user_id" = $1;

-- name: UpdateBankAccount :one
UPDATE bank_accounts SET "name" = $2, "bank_name" = $3, "account_type" = $4, "balance" = $5, "is_active" = $6
WHERE "id" = $1
RETURNING *;

-- name: DeleteBankAccount :exec
DELETE FROM bank_accounts WHERE "id" = $1;