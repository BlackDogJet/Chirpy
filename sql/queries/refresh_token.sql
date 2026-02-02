-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $2,
    $3
)
RETURNING *;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE token = $1
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
  AND refresh_tokens.revoked_at IS NULL
  AND refresh_tokens.expires_at > CURRENT_TIMESTAMP;
