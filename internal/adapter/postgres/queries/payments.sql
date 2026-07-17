-- name: CreatePayment :one
INSERT INTO payments (bill_id, month, amount_paid, proof_file_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPaymentByBillAndMonth :one
SELECT * FROM payments
WHERE bill_id = $1 AND month = $2;

-- name: ListPaymentsByUserAndMonth :many
SELECT p.* FROM payments p
INNER JOIN bills b ON b.id = p.bill_id
WHERE b.user_id = $1 AND p.month = $2;
