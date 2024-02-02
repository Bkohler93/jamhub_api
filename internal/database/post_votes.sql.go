// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: post_votes.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPostVote = `-- name: CreatePostVote :one
INSERT INTO post_votes(id, post_id, user_id, created_at, updated_at, is_up)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING id, post_id, user_id, created_at, updated_at, is_up
`

type CreatePostVoteParams struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	IsUp      bool
}

func (q *Queries) CreatePostVote(ctx context.Context, arg CreatePostVoteParams) (PostVote, error) {
	row := q.db.QueryRowContext(ctx, createPostVote,
		arg.ID,
		arg.PostID,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.IsUp,
	)
	var i PostVote
	err := row.Scan(
		&i.ID,
		&i.PostID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsUp,
	)
	return i, err
}

const deletePostVote = `-- name: DeletePostVote :exec
DELETE FROM post_votes WHERE id=$1
`

func (q *Queries) DeletePostVote(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePostVote, id)
	return err
}
