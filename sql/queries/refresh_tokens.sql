-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT u.* FROM users AS u
JOIN refresh_tokens AS rt ON u.id = rt.user_id
WHERE rt.token = $1
AND rt.revoked_at IS NULL
AND rt.expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $1, updated_at = $2
WHERE token = $3;