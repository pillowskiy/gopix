-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS image_properties (
    mime VARCHAR(255),
    ext VARCHAR(10),
    height INT,
    width INT,
    image_id BIGINT NOT NULL UNIQUE,

    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_image_properties_image_id ON image_properties(image_id);

INSERT INTO image_properties (mime, ext, height, width, image_id)
SELECT mime, ext, 0, 0, id FROM images;

ALTER TABLE images DROP COLUMN mime;
ALTER TABLE images DROP COLUMN ext;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE images ADD COLUMN mime VARCHAR(50) NOT NULL;
ALTER TABLE images ADD COLUMN ext VARCHAR(7) NOT NULL;

UPDATE images
SET mime = ip.mime,
    ext = ip.ext
FROM image_properties ip
WHERE images.id = ip.image_id;

DROP INDEX IF EXISTS idx_image_properties_image_id;

DROP TABLE IF EXISTS image_properties;
-- +goose StatementEnd
