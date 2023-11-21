-- name: CreatePostVote :one
INSERT INTO post_votes(id, post_id, user_id, created_at, updated_at, is_up)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: DeletePostVote :exec
DELETE FROM post_votes WHERE id=$1;