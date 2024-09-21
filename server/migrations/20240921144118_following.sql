-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "following" (
    "follower_id" BIGINT NOT NULL,
    "followed_id" BIGINT NOT NULL,
    "created_at" TIMESTAMP DEFAULT (current_timestamp),

    CONSTRAINT unique_follow UNIQUE ("follower_id", "followed_id"),

    FOREIGN KEY ("follower_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY ("followed_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_following_follower_id ON "following" ("follower_id");
CREATE INDEX idx_following_followed_id ON "following" ("followed_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_following_follower_id;
DROP INDEX IF EXISTS idx_following_followed_id;
DROP TABLE IF EXISTS "following";
-- +goose StatementEnd
