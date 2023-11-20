-- name: CreateRevokedToken :one
INSERT INTO revoked_tokens(id, user_id, revoked_at)
VALUES ($1,$2,$3)
RETURNING *;

-- name: GetRevokedToken :one
SELECT * FROM revoked_tokens WHERE id=$1;
