// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: refresh_tokens.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at
`

type CreateRefreshTokenParams struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken,
		arg.Token,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.ExpiresAt,
	)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT u.id, u.created_at, u.updated_at, u.email, u.hashed_password, u.is_chirpy_red FROM users AS u
JOIN refresh_tokens AS rt ON u.id = rt.user_id
WHERE rt.token = $1
AND rt.revoked_at IS NULL
AND rt.expires_at > NOW()
`

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $1, updated_at = $2
WHERE token = $3
`

type RevokeRefreshTokenParams struct {
	RevokedAt sql.NullTime
	UpdatedAt time.Time
	Token     string
}

func (q *Queries) RevokeRefreshToken(ctx context.Context, arg RevokeRefreshTokenParams) error {
	_, err := q.db.ExecContext(ctx, revokeRefreshToken, arg.RevokedAt, arg.UpdatedAt, arg.Token)
	return err
}
