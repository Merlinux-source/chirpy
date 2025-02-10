-- +goose Up
CREATE TABLE refresh_tokens
(
    token      TEXT PRIMARY KEY,
    created_at TIMESTAMP                                    NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP                                    NOT NULL DEFAULT current_timestamp,
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    expires_at TIMESTAMP                                    NOT NULL,
    revoked_at TIMESTAMP                                             DEFAULT NULL
);
-- +goose Down
DROP TABLE refresh_tokens;
