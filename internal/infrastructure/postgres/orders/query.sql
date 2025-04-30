-- name: CreateOrder :exec
INSERT INTO orders (
    number,
    user_id,
    status,
    accrual,
    created,
    updated
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: GetOrderByNumber :one
SELECT * FROM orders
WHERE number = $1;

-- name: GetOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1;

-- name: GetOrdersByLimitAndOffset :many
SELECT * FROM orders
WHERE status = ANY (sqlc.slice('statuses'))
ORDER BY id ASC
LIMIT $1
OFFSET $2;

-- name: UpdateOrder :exec
UPDATE orders
SET status = $1, accrual = $2, updated = $3
WHERE number = $4;

-- name: GetTotalAccrualByUserID :one
SELECT SUM(accrual)::real as total_accrual FROM orders
WHERE user_id = $1 AND status = $2;