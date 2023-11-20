-- +goose Up
CREATE TABLE revoked_tokens (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL,
	revoked_at TIMESTAMP NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE revoked_tokens;