-- +goose Up
CREATE TABLE room_subscriptions(
	id UUID PRIMARY KEY,
	room_id UUID NOT NULL REFERENCES rooms(id),
	user_id UUID NOT NULL REFERENCES users(id),
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE room_subscriptions;