-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens WHERE token = $1 AND revoked_at IS NULL;