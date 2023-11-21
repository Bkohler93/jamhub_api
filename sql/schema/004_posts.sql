-- +goose Up
CREATE TABLE posts(
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL REFERENCES users(id),
	room_id UUID NOT NULL REFERENCES rooms(id),
	link VARCHAR(1024) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE posts;