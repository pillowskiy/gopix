-- +goose Up
-- +goose StatementBegin

CREATE TABLE images_analytics (
  image_id BIGINT NOT NULL,
  likes_count INT DEFAULT 0,
  views_count INT DEFAULT 0,

  FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION create_images_analytics()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO images_analytics (image_id) VALUES (NEW.id);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_images_analytics
AFTER INSERT ON images
FOR EACH ROW
EXECUTE FUNCTION create_images_analytics();

INSERT INTO images_analytics (image_id, likes_count, views_count)
SELECT
    i.id AS image_id,
    COALESCE(l.likes_count, 0) AS likes_count,
    COALESCE(v.views_count, 0) AS views_count
FROM images i
LEFT JOIN (
    SELECT image_id, COUNT(*) AS likes_count
    FROM images_to_likes
    GROUP BY image_id
) l ON i.id = l.image_id
LEFT JOIN (
    SELECT image_id, COUNT(*) AS views_count
    FROM images_to_views
    GROUP BY image_id
) v ON i.id = v.image_id
WHERE i.id NOT IN (SELECT image_id FROM images_analytics);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_create_images_analytics ON images;
DROP FUNCTION IF EXISTS create_images_analytics();
DROP TABLE IF EXISTS images_analytics;
-- +goose StatementEnd