-- name: CreateRoomSubscription :one
INSERT INTO room_subscriptions(id, room_id, user_id, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5) RETURNING *;

-- name: DeleteRoomSubscription :exec
DELETE FROM room_subscriptions WHERE room_id=$1 AND user_id=$2;

-- name: GetUserRoomSubscriptions :many
SELECT * FROM room_subscriptions WHERE user_id=$1;

-- name: GetRoomRoomSubscriptions :many
SELECT * FROM room_subscriptions WHERE room_id=$1;

-- name: GetRoomSubscription :one
SELECT * FROM room_subscriptions WHERE user_id=$1 AND room_id=$1;

-- name: GetAllRoomSubs :many
SELECT * FROM room_subscriptions;