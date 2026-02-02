-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1,
    $2
)
RETURNING *;

-- name: GetChirpByID :one
SELECT * FROM chirps WHERE id = $1;

-- name: GetChirps :many
SELECT * FROM chirps WHERE true
ORDER BY created_at;

-- name: DeleteChirpByID :exec
DELETE FROM chirps WHERE id = $1;