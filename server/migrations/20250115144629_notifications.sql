-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "notifications" (
    "id" BIGINT DEFAULT generate_snowflake_id() PRIMARY KEY,
    "user_id" BIGINT NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "message" TEXT DEFAULT '',
    "hidden" BOOLEAN DEFAULT FALSE,
    "read" BOOLEAN DEFAULT FALSE,
    "sent_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_id FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS "idx_notifications_user_id" ON "notifications" ("user_id");
CREATE INDEX IF NOT EXISTS "idx_notifications_read" ON "notifications" ("read");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "notifications";
DROP INDEX IF EXISTS "idx_notifications_user_id";
DROP INDEX IF EXISTS "idx_notifications_read";
-- +goose StatementEnd
