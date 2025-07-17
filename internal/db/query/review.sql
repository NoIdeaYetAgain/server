-- name: CreateReview :one
INSERT INTO reviews (
    student_id, rating, comment
) VALUES (
             $1, $2, $3
         )
    RETURNING *;

-- name: GetReviewByID :one
SELECT * FROM reviews
WHERE id = $1;

-- name: ListReviews :many
SELECT * FROM reviews
ORDER BY created_at DESC
    LIMIT $1
OFFSET $2;

-- name: UpdateReview :one
UPDATE reviews
SET rating = $2,
    comment = $3,
    created_at = CURRENT_TIMESTAMP
WHERE id = $1
    RETURNING *;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1;
