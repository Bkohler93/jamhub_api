-- +goose Up
CREATE TABLE post_votes(
	id UUID PRIMARY KEY,
	post_id UUID NOT NULL REFERENCES posts(id),
	user_id UUID NOT NULL REFERENCES users(id),
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	is_up BOOLEAN NOT NULL
);

-- +goose Down
DROP TABLE post_votes;