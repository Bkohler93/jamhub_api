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
INSERT INTO post_votes(post_id, user_id, created_at, updated_at, is_up)
VALUES ($1,$2,$3,$4,$5)
RETURNING post_id, user_id, created_at, updated_at, is_up
`

type CreatePostVoteParams struct {
	PostID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	IsUp      bool
}

func (q *Queries) CreatePostVote(ctx context.Context, arg CreatePostVoteParams) (PostVote, error) {
	row := q.db.QueryRowContext(ctx, createPostVote,
		arg.PostID,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.IsUp,
	)
	var i PostVote
	err := row.Scan(
		&i.PostID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsUp,
	)
	return i, err
}

const deletePostVote = `-- name: DeletePostVote :exec
DELETE FROM post_votes WHERE post_id=$1 AND user_id=$2
`

type DeletePostVoteParams struct {
	PostID uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeletePostVote(ctx context.Context, arg DeletePostVoteParams) error {
	_, err := q.db.ExecContext(ctx, deletePostVote, arg.PostID, arg.UserID)
	return err
}

const getPostPostVotes = `-- name: GetPostPostVotes :many
SELECT post_id, user_id, created_at, updated_at, is_up FROM post_votes WHERE post_id=$1
`

func (q *Queries) GetPostPostVotes(ctx context.Context, postID uuid.UUID) ([]PostVote, error) {
	rows, err := q.db.QueryContext(ctx, getPostPostVotes, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PostVote
	for rows.Next() {
		var i PostVote
		if err := rows.Scan(
			&i.PostID,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.IsUp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPostVote = `-- name: GetPostVote :one
SELECT post_id, user_id, created_at, updated_at, is_up FROM post_votes WHERE post_id=$1 AND user_id=$2
`

type GetPostVoteParams struct {
	PostID uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) GetPostVote(ctx context.Context, arg GetPostVoteParams) (PostVote, error) {
	row := q.db.QueryRowContext(ctx, getPostVote, arg.PostID, arg.UserID)
	var i PostVote
	err := row.Scan(
		&i.PostID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsUp,
	)
	return i, err
}
