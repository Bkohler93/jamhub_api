-- name: CreatePostVote :one
INSERT INTO post_votes(post_id, user_id, created_at, updated_at, is_up)
VALUES ($1,$2,$3,$4,$5)
RETURNING *;

-- name: DeletePostVote :exec
DELETE FROM post_votes WHERE post_id=$1 AND user_id=$2;

-- name: GetPostVote :one
SELECT * FROM post_votes WHERE post_id=$1 AND user_id=$2;

-- name: GetPostPostVotes :many
SELECT * FROM post_votes WHERE post_id=$1;