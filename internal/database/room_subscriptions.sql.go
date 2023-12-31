// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: room_subscriptions.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createRoomSubscription = `-- name: CreateRoomSubscription :one
INSERT INTO room_subscriptions(id, room_id, user_id, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5) RETURNING id, room_id, user_id, created_at, updated_at
`

type CreateRoomSubscriptionParams struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateRoomSubscription(ctx context.Context, arg CreateRoomSubscriptionParams) (RoomSubscription, error) {
	row := q.db.QueryRowContext(ctx, createRoomSubscription,
		arg.ID,
		arg.RoomID,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i RoomSubscription
	err := row.Scan(
		&i.ID,
		&i.RoomID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteRoomSubscription = `-- name: DeleteRoomSubscription :exec
DELETE FROM room_subscriptions WHERE room_id=$1 AND user_id=$2
`

type DeleteRoomSubscriptionParams struct {
	RoomID uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteRoomSubscription(ctx context.Context, arg DeleteRoomSubscriptionParams) error {
	_, err := q.db.ExecContext(ctx, deleteRoomSubscription, arg.RoomID, arg.UserID)
	return err
}

const getAllRoomSubs = `-- name: GetAllRoomSubs :many
SELECT id, room_id, user_id, created_at, updated_at FROM room_subscriptions
`

func (q *Queries) GetAllRoomSubs(ctx context.Context) ([]RoomSubscription, error) {
	rows, err := q.db.QueryContext(ctx, getAllRoomSubs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RoomSubscription
	for rows.Next() {
		var i RoomSubscription
		if err := rows.Scan(
			&i.ID,
			&i.RoomID,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getRoomRoomSubscriptions = `-- name: GetRoomRoomSubscriptions :many
SELECT id, room_id, user_id, created_at, updated_at FROM room_subscriptions WHERE room_id=$1
`

func (q *Queries) GetRoomRoomSubscriptions(ctx context.Context, roomID uuid.UUID) ([]RoomSubscription, error) {
	rows, err := q.db.QueryContext(ctx, getRoomRoomSubscriptions, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RoomSubscription
	for rows.Next() {
		var i RoomSubscription
		if err := rows.Scan(
			&i.ID,
			&i.RoomID,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getRoomSubscription = `-- name: GetRoomSubscription :one
SELECT id, room_id, user_id, created_at, updated_at FROM room_subscriptions WHERE user_id=$1 AND room_id=$1
`

func (q *Queries) GetRoomSubscription(ctx context.Context, userID uuid.UUID) (RoomSubscription, error) {
	row := q.db.QueryRowContext(ctx, getRoomSubscription, userID)
	var i RoomSubscription
	err := row.Scan(
		&i.ID,
		&i.RoomID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserRoomSubscriptions = `-- name: GetUserRoomSubscriptions :many
SELECT id, room_id, user_id, created_at, updated_at FROM room_subscriptions WHERE user_id=$1
`

func (q *Queries) GetUserRoomSubscriptions(ctx context.Context, userID uuid.UUID) ([]RoomSubscription, error) {
	rows, err := q.db.QueryContext(ctx, getUserRoomSubscriptions, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RoomSubscription
	for rows.Next() {
		var i RoomSubscription
		if err := rows.Scan(
			&i.ID,
			&i.RoomID,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
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
