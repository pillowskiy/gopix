-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_images_to_views_image_id
  ON images_to_views(
    image_id
  );
CREATE INDEX IF NOT EXISTS idx_images_to_likes_image_id
  ON images_to_likes(
    image_id
  );

CREATE INDEX IF NOT EXISTS idx_images_to_tags_image_id
  ON images_to_tags(
    image_id
  );

CREATE INDEX IF NOT EXISTS idx_images_to_tags_tag_id
  ON images_to_tags(
    tag_id
  );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_images_to_views_image_id;
DROP INDEX IF EXISTS idx_images_to_likes_image_id;

DROP INDEX IF EXISTS idx_images_to_tags_image_id
DROP INDEX IF EXISTS idx_images_to_tags_tag_id
-- +goose StatementEnd
