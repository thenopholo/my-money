-- name: CreateCategory :one
INSERT INTO categories ("user_id", "name", "category_type", "color", "icon")
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE "id" = $1;

-- name: GetCategoryByUserID :many
SELECT * FROM categories WHERE "user_id" = $1;

-- name: UpdateCategory :one
UPDATE categories SET "name" = $2, "category_type" = $3, "color" = $4, "icon" = $5
WHERE "id" = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE "id" = $1;