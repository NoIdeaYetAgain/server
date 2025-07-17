-- name: CreateProperty :one
INSERT INTO properties (
    agent_id, title, description, price, location, property_type, image_url, video_url
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8
         )
    RETURNING *;

-- name: GetProperty :one
SELECT * FROM properties
WHERE id = $1 LIMIT 1;

-- name: ListProperties :many
SELECT * FROM properties
ORDER BY created_at DESC
    LIMIT $1
OFFSET $2;

-- name: UpdateProperty :one
UPDATE properties
SET title = $2,
    description = $3,
    price = $4,
    location = $5,
    property_type = $6,
    image_url = $7,
    video_url = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
    RETURNING *;

-- name: DeleteProperty :exec
DELETE FROM properties
WHERE id = $1;