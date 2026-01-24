-- name: CreateCreditCard :one
INSERT INTO credit_cards ("user_id", "name", "credit_limit", "close_day", "due_day", "is_active")
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetCreditCardByID :one
SELECT * FROM credit_cards WHERE "id" = $1;

-- name: GetCreditCardByUserID :many
SELECT * FROM credit_cards WHERE "user_id" = $1;

-- name: UpdateCreditCard :one
UPDATE credit_cards SET "name" = $2, "credit_limit" = $3, "close_day" = $4, "due_day" = $5, "is_active" = $6
WHERE "id" = $1
RETURNING *;

-- name: DeleteCreditCard :exec
DELETE FROM credit_cards WHERE "id" = $1;