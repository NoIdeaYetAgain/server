-- name: CreateInspectionRequest :one
INSERT INTO inspection_requests (
    property_id, student_id, requested_date
) VALUES (
             $1, $2, $3
         )
    RETURNING *;

-- name: GetInspectionRequest :one
SELECT * FROM inspection_requests
WHERE id = $1;

-- name: ListInspections :many
SELECT * FROM inspection_requests
ORDER BY created_at DESC
    LIMIT $1
OFFSET $2;

-- name: DeleteInspectionRequest :exec
DELETE FROM inspection_requests
WHERE id = $1;
