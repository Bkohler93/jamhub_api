-- +goose Up
CREATE TABLE rooms(
	id UUID PRIMARY KEY,
	name VARCHAR(128) UNIQUE NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);


-- +goose Down
DROP TABLE rooms;