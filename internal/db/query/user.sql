-- name: CreateUser :one
INSERT INTO users (
    full_name, email, phone_number, password_hash, role
) VALUES (
             $1, $2, $3, $4, $5
         )
    RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET full_name = $2,
    email = $3,
    phone_number = $4,
    password_hash = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
