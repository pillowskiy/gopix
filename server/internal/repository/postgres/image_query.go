package postgres

const createImageQuery = `
INSERT INTO images (author_id, path, title, description, access_level, expires_at, mime, ext)
VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, '')::access_level, 'public'::access_level), $6, $7, $8)
RETURNING *
`

const getByIdImageQuery = `SELECT * FROM images WHERE id = $1`

const deleteImageQuery = `DELETE FROM images WHERE id = $1`

const getDetailedImageQuery = `
SELECT
  i.*,
  u.id AS "author.id",
  u.username AS "author.username",
  u.avatar_url AS "author.avatar_url",
  COALESCE(a.likes_count, 0) AS likes,
  COALESCE(a.views_count, 0) AS views,
  TO_JSON(COALESCE(
    ARRAY_AGG(
      json_build_object('id', t.id, 'name', t.name)
    ) FILTER (WHERE t.id IS NOT NULL),
    '{}'
  )) AS tags
FROM
  images i
JOIN
  users u ON i.author_id = u.id
LEFT JOIN
  images_to_tags it ON i.id = it.image_id
LEFT JOIN
  tags t ON it.tag_id = t.id
LEFT JOIN
  images_analytics a ON a.image_id = i.id
WHERE
  i.id = $1
GROUP BY
  i.id, u.id, a.likes_count, a.views_count
`

const findManyImagesQuery = `
SELECT
  u.id AS "author.id",
  u.username AS "author.username",
  u.avatar_url AS "author.avatar_url",
  i.*
FROM images i
INNER JOIN users u ON i.author_id = u.id
WHERE i.id IN(?);
`

const updateImageQuery = `
UPDATE images SET
  title = COALESCE(NULLIF($1, ''), title),
  description = COALESCE(NULLIF($2, ''), description),
  access_level = COALESCE(NULLIF($3, '')::access_level, access_level)::access_level,
  expires_at = COALESCE($4, expires_at)
WHERE id = $5 RETURNING *`

const statesImageQuery = `
WITH params AS (SELECT $1::int AS image_id, $2::int AS user_id)
SELECT 
  $1::int as image_id,
  EXISTS (
    SELECT 1
    FROM images_to_views v
    JOIN params p ON v.image_id = p.image_id AND v.user_id = p.user_id
  ) AS viewed,
  EXISTS (
    SELECT 1
    FROM images_to_likes l
    JOIN params p ON l.image_id = p.image_id AND l.user_id = p.user_id
  ) AS liked;
`

const hasLikeImageQuery = `SELECT EXISTS (SELECT 1 FROM images_to_likes WHERE image_id = $1 AND user_id = $2)`
