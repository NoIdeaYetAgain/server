-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
     user_id,
     email,
     token
) VALUES (
          $1, $2, $3
         ) RETURNING *;