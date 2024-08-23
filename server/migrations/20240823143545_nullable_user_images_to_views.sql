-- +goose Up
-- +goose StatementBegin
ALTER TABLE "images_to_views" ALTER COLUMN "user_id" DROP NOT NULL;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE "images_to_views" ALTER COLUMN "user_id" SET NOT NULL;
-- +goose StatementEnd
