-- name: CreateUser :exec
INSERT INTO users (
    login,
    password_hash,
    created,
    updated
) VALUES (
    $1,
    $2,
    $3,
    $4
);

-- name: GetUserByLogin :one
SELECT * FROM users
WHERE login = $1;