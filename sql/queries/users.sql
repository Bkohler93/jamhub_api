-- name: CreateUser :one
INSERT INTO users(id, email, phone, password_hash, display_name, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING *;

-- name: GetUserByDisplayName :one
SELECT * FROM users WHERE display_name = $1;

-- name: UpdateUser :one
UPDATE users
SET 
	email=COALESCE(sqlc.narg(email)::VARCHAR(64),email),
	phone=COALESCE(sqlc.narg(phone)::VARCHAR(64),phone),
	updated_at=sqlc.arg(updated_at)::TIMESTAMP,
	password_hash = CASE WHEN sqlc.narg('password_hash')::VARCHAR(100) IS NULL THEN password_hash ELSE sqlc.narg('password_hash')::VARCHAR(100) END,
	display_name = CASE WHEN sqlc.narg('display_name')::VARCHAR(64) IS NULL THEN display_name ELSE sqlc.narg('display_name')::VARCHAR(64) END
WHERE id=sqlc.arg(id)::UUID
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id=$1;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone=$1;


-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;