-- name: CreatePlannedExpense :one
INSERT INTO planned_expenses ("user_id", "account_id", "category_id", "amount", "due_day", "start_date", "end_date", "description", "frequency", "is_active")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
RETURNING *;

-- name: GetPlannedExpenseByID :one
SELECT * FROM planned_expenses WHERE "id" = $1;

-- name: GetPlannedExpenseByUserID :many
SELECT * FROM planned_expenses WHERE "user_id" = $1;

-- name: UpdatePlannedExpense :one
UPDATE planned_expenses SET "amount" = $2, "due_day" = $3, "start_date" = $4, "end_date" = $5, "description" = $6, "frequency" = $7, "is_active" = $8
WHERE "id" = $1
RETURNING *;

-- name: DeletePlannedExpense :exec
DELETE FROM planned_expenses WHERE "id" = $1;