-- name: AddFavorite :one
INSERT INTO favorites (student_id, property_id)
VALUES ($1, $2)
    RETURNING *;

-- name: RemoveFavorite :exec
DELETE FROM favorites
WHERE student_id = $1 AND property_id = $2;

-- name: GetFavoritesApartments :many
SELECT p.*
FROM properties p
         JOIN favorites f ON f.property_id = p.id
WHERE f.student_id = $1
ORDER BY p.created_at DESC
    LIMIT $2 OFFSET $3;
