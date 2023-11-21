-- +goose Up
INSERT INTO users (id, email, phone, password_hash, display_name, created_at, updated_at) VALUES
    ('1f4ab37a-788b-4e9e-8a5e-3a7b6b8bb4b1', 'user1@example.com', '1234567890', '$2a$10$z7tSrmzYJDfqcyQpRmsvX.05thN8gpTRXOvNBP7RxGbCnHqcAx1iO', 'user1', '2023-01-01 12:00:00', '2023-01-01 12:00:00'),
    ('2f8cde2a-7d9f-4b3c-92b3-2a1e5d4c3f8d', 'user2@example.com', '9876543210', '$2a$10$py5xHHKn5QXOIiQpWE9DP.QvpjFCgMkKZvTMsj1pFnyI29MyOPjXW', 'user2', '2023-01-02 12:30:00', '2023-01-02 12:30:00'),
    ('3e6cf8a4-865d-4f9a-b4f2-1c4b5d6e7c3e', 'user3@example.com', '5556667777', '$2a$10$k.77jlGkUvanqL5WLbJb8uEFT0vawCeBKrhMb/oaTWyvyV4DOMbR6', 'user3', '2023-01-03 13:00:00', '2023-01-03 13:00:00');

INSERT INTO rooms (id, name, created_at, updated_at) VALUES
    ('4a4d5e1b-2cf4-4a47-b1f7-6d67e3b8a5f1', 'Chumab Wumba', '2023-01-01 12:00:00', '2023-01-01 12:00:00'),
    ('5b5f6c2d-8e9a-4f7c-a3d2-1e2d3f4b5c6d', 'The Beatles', '2023-01-02 12:30:00', '2023-01-02 12:30:00'),
    ('6c6f7a3e-9d8b-4c6a-b5d4-3e4f5a6b7a8e', '80s Pop', '2023-01-03 13:00:00', '2023-01-03 13:00:00');

INSERT INTO posts (id, user_id, room_id, link, created_at, updated_at) VALUES
    ('7d7e6b4a-2c1b-4a3f-8b7a-9e8e7f6d5b4a', '1f4ab37a-788b-4e9e-8a5e-3a7b6b8bb4b1', '4a4d5e1b-2cf4-4a47-b1f7-6d67e3b8a5f1', 'https://open.spotify.com/track/2BeInbvK9KLJVKGyNKmyne?si=cd362e47ce8347e6', '2023-01-01 12:15:00', '2023-01-01 12:15:00'),
    ('8d9f8c2a-4c3f-4b7c-9d2a-8d7b3f4c3f8d', '2f8cde2a-7d9f-4b3c-92b3-2a1e5d4c3f8d', '5b5f6c2d-8e9a-4f7c-a3d2-1e2d3f4b5c6d', 'https://open.spotify.com/track/44V3prBS5nmjRlQN0lcbeK?si=d7cef68232c147b0', '2023-01-02 13:00:00', '2023-01-02 13:00:00'),
    ('9a4d5e1b-7d9f-4a3f-8b7a-3e7b6d5e4f5a', '3e6cf8a4-865d-4f9a-b4f2-1c4b5d6e7c3e', '6c6f7a3e-9d8b-4c6a-b5d4-3e4f5a6b7a8e', 'https://open.spotify.com/track/1LVE84dyyEKAW1uptJguTL?si=5f1c3e69d9bc4c5b', '2023-01-03 14:30:00', '2023-01-03 14:30:00');

INSERT INTO room_subscriptions (id, room_id, user_id, created_at, updated_at) VALUES
    ('1a2b3c4d-5e6f-7a8b-9c1d-2e3f4a5b6c7d', '4a4d5e1b-2cf4-4a47-b1f7-6d67e3b8a5f1', '1f4ab37a-788b-4e9e-8a5e-3a7b6b8bb4b1', '2023-01-01 12:30:00', '2023-01-01 12:30:00'),
    ('2d3e4f5a-6b7c-8d9e-1f2a-3b4c5d6e7f8a', '5b5f6c2d-8e9a-4f7c-a3d2-1e2d3f4b5c6d', '2f8cde2a-7d9f-4b3c-92b3-2a1e5d4c3f8d', '2023-01-02 14:00:00', '2023-01-02 14:00:00'),
    ('3b4c5d6e-7f8a-9b1c-2d3e-4f5a6b7c8d9e', '6c6f7a3e-9d8b-4c6a-b5d4-3e4f5a6b7a8e', '3e6cf8a4-865d-4f9a-b4f2-1c4b5d6e7c3e', '2023-01-03 15:30:00', '2023-01-03 15:30:00');

INSERT INTO post_votes (id, post_id, user_id, created_at, updated_at, is_up) VALUES
    ('4a5b6c7d-8e9f-1a2b-3c4d-5e6f7a8b9c1d', '7d7e6b4a-2c1b-4a3f-8b7a-9e8e7f6d5b4a', '1f4ab37a-788b-4e9e-8a5e-3a7b6b8bb4b1', '2023-01-01 13:00:00', '2023-01-01 13:00:00', TRUE),
    ('5c6d7e8f-9a1b-2c3d-4e5f-6a7b8c9d1e2f', '8d9f8c2a-4c3f-4b7c-9d2a-8d7b3f4c3f8d', '2f8cde2a-7d9f-4b3c-92b3-2a1e5d4c3f8d', '2023-01-02 15:00:00', '2023-01-02 15:00:00', FALSE),
    ('6e7f8a9b-1c2d-3e4f-5a6b-7c8d9e1f2a3b', '9a4d5e1b-7d9f-4a3f-8b7a-3e7b6d5e4f5a', '3e6cf8a4-865d-4f9a-b4f2-1c4b5d6e7c3e', '2023-01-03 16:30:00', '2023-01-03 16:30:00', TRUE);


-- +goose Down
