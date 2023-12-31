// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: revoked_tokens.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createRevokedToken = `-- name: CreateRevokedToken :one
INSERT INTO revoked_tokens(id, user_id, revoked_at)
VALUES ($1,$2,$3)
RETURNING id, user_id, revoked_at
` // #nosec G101

type CreateRevokedTokenParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	RevokedAt time.Time
}

func (q *Queries) CreateRevokedToken(ctx context.Context, arg CreateRevokedTokenParams) (RevokedToken, error) {
	row := q.db.QueryRowContext(ctx, createRevokedToken, arg.ID, arg.UserID, arg.RevokedAt)
	var i RevokedToken
	err := row.Scan(&i.ID, &i.UserID, &i.RevokedAt)
	return i, err
}

const getRevokedToken = `-- name: GetRevokedToken :one
SELECT id, user_id, revoked_at FROM revoked_tokens WHERE id=$1
` // #nosec G101 --no hardcoded credentials found

func (q *Queries) GetRevokedToken(ctx context.Context, id uuid.UUID) (RevokedToken, error) {
	row := q.db.QueryRowContext(ctx, getRevokedToken, id)
	var i RevokedToken
	err := row.Scan(&i.ID, &i.UserID, &i.RevokedAt)
	return i, err
}
