-- name: CreateWithdrawal :exec
INSERT INTO withdrawals (
    order_id,
    sum,
    user_id,
    created,
    updated
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: GetWithdrawalsByUserID :many
SELECT * FROM withdrawals
WHERE user_id = $1;

-- name: GetTotalWithdrawnByUserID :one
SELECT COALESCE(SUM(sum), 0)::real as total_withdrawn FROM withdrawals
WHERE user_id = $1;