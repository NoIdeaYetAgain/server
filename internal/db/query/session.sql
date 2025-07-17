-- name: CreateRefreshToken :one
INSERT INTO sessions (
    id,
    refresh_token,
    user_id,
    user_agent,
    client_ip,
    is_blocked,
    expires_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: BlockRefreshToken :exec
UPDATE sessions
SET is_blocked = TRUE
WHERE id = $1;
