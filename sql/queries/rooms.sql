-- name: CreateRoom :one
INSERT INTO rooms(id, name, created_at, updated_at)
VALUES ($1,$2,$3,$4)
RETURNING *;

-- name: GetRooms :many
SELECT * FROM rooms
ORDER BY (
	SELECT COUNT(*) FROM room_subscriptions WHERE id=room_id
) 
LIMIT $1;

-- name: GetRoomByID :one
SELECT * FROM rooms WHERE id=$1;

-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id=$1;

--use this to retrieve a user's subbed rooms ordered by date subbed, seen in Home Page of app underneath "Your subscriptions" heading
-- name: GetUserRoomsOrderedBySubs :many
SELECT
	rooms.id AS room_id,
	rooms.name AS room_name,
	rooms.created_at,
	rooms.updated_at,
	COUNT(room_subscriptions.room_id) AS subscription_count
FROM	
	rooms
LEFT JOIN
	room_subscriptions on rooms.id = room_subscriptions.room_id
WHERE
	room_subscriptions.user_id = $1
GROUP BY
	rooms.id, rooms.name, rooms.created_at, rooms.updated_at, room_subscriptions.created_at
ORDER BY	room_subscriptions.created_at DESC
LIMIT $2 OFFSET $3;

--use this to retrieve "Top" rooms ordered by number of subs a room has, seen in Home Page of app underneath "TOP" heading
-- name: GetRoomsOrderedBySubs :many
SELECT
	rooms.id AS room_id,
	rooms.name AS room_name,
	rooms.created_at,
	rooms.updated_at,
	COUNT(room_subscriptions.room_id) AS subscription_count
FROM	
	rooms
LEFT JOIN
	room_subscriptions on rooms.id = room_subscriptions.room_id
GROUP BY
	rooms.id, rooms.name, rooms.created_at, rooms.updated_at
ORDER BY	subscription_count DESC
LIMIT $1 OFFSET $2;