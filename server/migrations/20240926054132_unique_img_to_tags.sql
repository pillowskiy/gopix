-- +goose Up
-- +goose StatementBegin
ALTER TABLE images_to_tags ADD CONSTRAINT unique_image_tag UNIQUE(tag_id, image_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE images_to_tags DROP CONSTRAINT IF EXISTS unique_image_tag;
-- +goose StatementEnd
