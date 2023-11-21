-- name: CreatePost :one
INSERT INTO posts(id, user_id, room_id, link, updated_at, created_at)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id=$1;

-- name: GetRoomPosts :many
SELECT * FROM posts WHERE room_id=$1;

-- name: GetPost :one
SELECT * FROM posts WHERE id=$1;

--use to retrieve the "New" posts in a given room
-- name: GetNewRoomPosts :many
SELECT 
	posts.id,
	posts.user_id,
	posts.room_id,
	posts.link,
	posts.created_at,
	posts.updated_at
FROM
	posts
INNER JOIN
	rooms
ON
	posts.room_id=rooms.id
WHERE posts.room_id = $1
ORDER BY
	posts.created_at DESC;


--use to retrieve the "Top" posts in a given room
-- name: GetTopRoomPosts :many
SELECT
	posts.id as post_id,
	posts.room_id,
	posts.link,
	posts.created_at,
	posts.updated_at,
	SUM(CASE WHEN post_votes.is_up=TRUE THEN 1 ELSE 0 END) as num_upvotes
FROM
	posts
LEFT JOIN
	post_votes
ON
	posts.id = post_votes.post_id
WHERE posts.room_id=$1
GROUP BY
	posts.id, posts.room_id, posts.link, posts.created_at, posts.updated_at
ORDER BY
	num_upvotes DESC;
