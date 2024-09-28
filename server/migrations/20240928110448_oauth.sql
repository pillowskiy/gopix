-- +goose Up
-- +goose StatementBegin
-- Додаємо колонку "external" до таблиці "users"
ALTER TABLE "users"
ADD COLUMN external BOOLEAN DEFAULT FALSE NOT NULL;

CREATE TABLE IF NOT EXISTS "oauth" (
    "user_id" BIGINT NOT NULL,
    "oauth_id" VARCHAR(255) NOT NULL,
    "service" VARCHAR(50) NOT NULL,
    "created_at" timestamp DEFAULT (current_timestamp),
    "updated_at" timestamp DEFAULT (current_timestamp),
    CONSTRAINT unique_user_oauth UNIQUE(user_id, oauth_id)
);

ALTER TABLE "oauth"
ADD CONSTRAINT fk_oauth_user_id FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "users"
DROP COLUMN IF EXISTS external;

DELETE TABLE IF EXISTS "oauth";
-- +goose StatementEnd
