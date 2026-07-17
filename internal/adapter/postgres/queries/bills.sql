-- name: CreateBill :one
INSERT INTO bills (user_id, name, category, service_provider, expected_amount, due_day)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetBillByID :one
SELECT * FROM bills
WHERE id = $1 AND user_id = $2;

-- name: ListBillsByUser :many
SELECT * FROM bills
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateBill :one
UPDATE bills
SET name = $3,
    category = $4,
    service_provider = $5,
    expected_amount = $6,
    due_day = $7,
    status = $8,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteBill :exec
DELETE FROM bills
WHERE id = $1 AND user_id = $2;
