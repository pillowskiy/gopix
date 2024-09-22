-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "comments_to_likes" (
    "user_id" BIGINT NOT NULL,
    "comment_id" BIGINT NOT NULL,
    "created_at" timestamp DEFAULT current_timestamp,

    CONSTRAINT unique_comment_like UNIQUE ("user_id", "comment_id"),

    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY ("comment_id") REFERENCES "comments" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX idx_comments_to_likes_user_id ON comments_to_likes(user_id);
CREATE INDEX idx_comments_to_likes_comment_id ON comments_to_likes(comment_id);
ALTER TABLE "comments"
    ADD COLUMN "parent_id" BIGINT DEFAULT NULL,
    ADD CONSTRAINT fk_comments_parent FOREIGN KEY ("parent_id") REFERENCES "comments" ("id") ON DELETE CASCADE;

CREATE INDEX idx_comments_parent_id ON comments(parent_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "comments_to_likes";
ALTER TABLE "comments" DROP COLUMN IF EXISTS "parent_id";
-- +goose StatementEnd
