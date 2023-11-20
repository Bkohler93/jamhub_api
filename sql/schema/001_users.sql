-- +goose Up
CREATE TABLE users(
	id UUID PRIMARY KEY,
	email VARCHAR(64),
	phone VARCHAR(64),
	password_hash VARCHAR(100) NOT NULL,
	display_name VARCHAR(64) UNIQUE NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);


-- +goose Down
DROP TABLE users;