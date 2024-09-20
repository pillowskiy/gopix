-- +goose Up
-- +goose StatementBegin

CREATE SEQUENCE public.global_id_seq;
ALTER SEQUENCE public.global_id_seq OWNER TO postgres;

CREATE OR REPLACE FUNCTION generate_snowflake_id()
    RETURNS bigint
    LANGUAGE 'plpgsql'
AS $BODY$
DECLARE
    our_epoch bigint := 1314220021721;
    seq_id bigint;
    now_millis bigint;
    shard_id int := 1;
    result bigint := 0;
BEGIN
    SELECT nextval('public.global_id_seq') % 1024 INTO seq_id;

    SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis;
    result := (now_millis - our_epoch) << 23;
    result := result | (shard_id << 10);
    result := result | (seq_id);
	return result;
END;
$BODY$;

ALTER FUNCTION generate_snowflake_id() OWNER TO postgres;

CREATE TABLE
    IF NOT EXISTS "users" (
        "id" BIGINT DEFAULT generate_snowflake_id() PRIMARY KEY,
        "username" varchar(255) UNIQUE NOT NULL,
        "email" varchar(255) UNIQUE NOT NULL,
        "permissions" bigint NOT NULL DEFAULT 0,
        "password_hash" varchar(255) NOT NULL,
        "avatar_url" varchar(255) DEFAULT '',
        "created_at" timestamp DEFAULT (current_timestamp),
        "updated_at" timestamp DEFAULT (current_timestamp)
    );

CREATE UNIQUE INDEX idx_users_username ON users(username);

CREATE TYPE access_level AS ENUM ('link', 'private', 'public');

CREATE TABLE
    IF NOT EXISTS "images" (
        "id" BIGINT DEFAULT generate_snowflake_id() PRIMARY KEY,
        "author_id" BIGINT NOT NULL,
        "path" varchar(255) NOT NULL,
        "title" varchar(255) DEFAULT '',
        "description" text DEFAULT '',
        "ext" varchar(7) NOT NULL,
        "mime" varchar(50) NOT NULL,
        "access_level" access_level DEFAULT 'link',
        "expires_at" timestamp,
        "uploaded_at" timestamp DEFAULT (current_timestamp),
        "updated_at" timestamp DEFAULT (current_timestamp)
    );

CREATE INDEX idx_images_author_id ON images(author_id);
CREATE INDEX idx_images_access_level ON images(access_level);

CREATE TABLE
    IF NOT EXISTS "albums" (
        "id" BIGINT DEFAULT generate_snowflake_id() PRIMARY KEY,
        "author_id" BIGINT NOT NULL,
        "name" varchar(255) NOT NULL,
        "description" text DEFAULT '',
        "created_at" timestamp DEFAULT (current_timestamp),
        "updated_at" timestamp DEFAULT (current_timestamp)
    );

CREATE INDEX idx_albums_author_id ON albums(author_id);

CREATE TABLE
    IF NOT EXISTS "images_to_albums" (
        "album_id" BIGINT NOT NULL,
        "image_id" BIGINT NOT NULL
    );

CREATE INDEX idx_images_to_albums_album_id ON images_to_albums(album_id);
CREATE INDEX idx_images_to_albums_image_id ON images_to_albums(image_id);

CREATE TABLE
    IF NOT EXISTS "tags" (
        "id" BIGINT DEFAULT generate_snowflake_id() PRIMARY KEY,
        "name" varchar(50) UNIQUE NOT NULL,
        "created_at" timestamp DEFAULT (current_timestamp),
        "updated_at" timestamp DEFAULT (current_timestamp)
    );

CREATE UNIQUE INDEX idx_tags_name ON tags(name);

CREATE TABLE
    IF NOT EXISTS "comments" (
        "id" BIGINT DEFAULT generate_snowflake_id() PRIMARY KEY,
        "author_id" BIGINT NOT NULL,
        "image_id" BIGINT NOT NULL,
        "comment" text NOT NULL,
        "created_at" timestamp DEFAULT (current_timestamp),
        "updated_at" timestamp DEFAULT (current_timestamp)
    );

CREATE INDEX idx_comments_author_id ON comments(author_id);
CREATE INDEX idx_comments_image_id ON comments(image_id);

CREATE TABLE
    IF NOT EXISTS "images_to_tags" (
        "image_id" BIGINT NOT NULL,
        "tag_id" BIGINT NOT NULL
    );

CREATE INDEX idx_images_to_tags_image_id ON images_to_tags(image_id);
CREATE INDEX idx_images_to_tags_tag_id ON images_to_tags(tag_id);

CREATE TABLE
    IF NOT EXISTS "images_to_views" (
        "user_id" BIGINT,
        "image_id" BIGINT NOT NULL
    );

CREATE INDEX idx_images_to_views_user_id ON images_to_views(user_id);
CREATE INDEX idx_images_to_views_image_id ON images_to_views(image_id);

CREATE TABLE
    IF NOT EXISTS "images_to_likes" (
        "user_id" BIGINT NOT NULL,
        "image_id" BIGINT NOT NULL,

        CONSTRAINT unique_like UNIQUE("user_id", "image_id")
    );

CREATE INDEX idx_images_to_likes_user_id ON images_to_likes(user_id);
CREATE INDEX idx_images_to_likes_image_id ON images_to_likes(image_id);

ALTER TABLE "images"
ADD CONSTRAINT fk_images_author_id FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE "albums"
ADD CONSTRAINT fk_albums_author_id FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE "images_to_albums"
ADD CONSTRAINT fk_images_to_albums_album_id FOREIGN KEY ("album_id") REFERENCES "albums" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "images_to_albums"
ADD CONSTRAINT fk_images_to_albums_image_id FOREIGN KEY ("image_id") REFERENCES "images" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "comments"
ADD CONSTRAINT fk_comments_author_id FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE "comments"
ADD CONSTRAINT fk_comments_image_id FOREIGN KEY ("image_id") REFERENCES "images" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "images_to_tags"
ADD CONSTRAINT fk_images_to_tags_image_id FOREIGN KEY ("image_id") REFERENCES "images" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "images_to_tags"
ADD CONSTRAINT fk_images_to_tags_tag_id FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "images_to_views"
ADD CONSTRAINT fk_images_to_views_user_id FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "images_to_views"
ADD CONSTRAINT fk_images_to_views_image_id FOREIGN KEY ("image_id") REFERENCES "images" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "images_to_likes"
ADD CONSTRAINT fk_images_to_likes_user_id FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "images_to_likes"
ADD CONSTRAINT fk_images_to_likes_image_id FOREIGN KEY ("image_id") REFERENCES "images" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Foreign Keys
ALTER TABLE "images_to_likes" DROP CONSTRAINT IF EXISTS fk_images_to_likes_image_id;
ALTER TABLE "images_to_likes" DROP CONSTRAINT IF EXISTS fk_images_to_likes_user_id;

ALTER TABLE "images_to_views" DROP CONSTRAINT IF EXISTS fk_images_to_views_image_id;
ALTER TABLE "images_to_views" DROP CONSTRAINT IF EXISTS fk_images_to_views_user_id;

ALTER TABLE "images_to_tags" DROP CONSTRAINT IF EXISTS fk_images_to_tags_tag_id;
ALTER TABLE "images_to_tags" DROP CONSTRAINT IF EXISTS fk_images_to_tags_image_id;

ALTER TABLE "comments" DROP CONSTRAINT IF EXISTS fk_comments_image_id;
ALTER TABLE "comments" DROP CONSTRAINT IF EXISTS fk_comments_author_id;

ALTER TABLE "images_to_albums" DROP CONSTRAINT IF EXISTS fk_images_to_albums_image_id;
ALTER TABLE "images_to_albums" DROP CONSTRAINT IF EXISTS fk_images_to_albums_album_id;

ALTER TABLE "albums" DROP CONSTRAINT IF EXISTS fk_albums_author_id;

ALTER TABLE "images" DROP CONSTRAINT IF EXISTS fk_images_author_id;

-- Tables
DROP TABLE IF EXISTS "images_to_likes";
DROP TABLE IF EXISTS "images_to_views";
DROP TABLE IF EXISTS "images_to_tags";
DROP TABLE IF EXISTS "comments";
DROP TABLE IF EXISTS "tags";
DROP TABLE IF EXISTS "images_to_albums";
DROP TABLE IF EXISTS "albums";
DROP TABLE IF EXISTS "images";
DROP TABLE IF EXISTS "users";

-- Indexes
DROP INDEX IF EXISTS idx_images_author_id;
DROP INDEX IF EXISTS idx_images_access_level;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_images_to_albums_album_id;
DROP INDEX IF EXISTS idx_images_to_albums_image_id;
DROP INDEX IF EXISTS idx_comments_author_id;
DROP INDEX IF EXISTS idx_comments_image_id;
DROP INDEX IF EXISTS idx_tags_name;
DROP INDEX IF EXISTS idx_images_to_tags_image_id;
DROP INDEX IF EXISTS idx_images_to_tags_tag_id;
DROP INDEX IF EXISTS idx_images_to_views_user_id;
DROP INDEX IF EXISTS idx_images_to_views_image_id;
DROP INDEX IF EXISTS idx_images_to_likes_user_id;
DROP INDEX IF EXISTS idx_images_to_likes_image_id;
DROP INDEX IF EXISTS idx_albums_author_id;

DROP TYPE IF EXISTS access_level;
DROP FUNCTION IF EXISTS generate_snowflake_id;
-- +goose StatementEnd
