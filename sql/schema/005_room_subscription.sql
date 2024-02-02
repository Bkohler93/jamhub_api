-- +goose Up
CREATE TABLE room_subscriptions(
	room_id UUID NOT NULL REFERENCES rooms(id),
	user_id UUID NOT NULL REFERENCES users(id),
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	PRIMARY KEY (room_id, user_id)
);

-- +goose Down
DROP TABLE room_subscriptions;