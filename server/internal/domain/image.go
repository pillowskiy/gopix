package domain

import (
	"time"
)

type ImageSortMethod string

const (
	ImageNewestSort     ImageSortMethod = "newest"
	ImageOldestSort     ImageSortMethod = "oldest"
	ImagePopularSort    ImageSortMethod = "popular"
	ImageMostViewedSort ImageSortMethod = "mostViewed"
)

type ImageAccessLevel string

const (
	ImageAccessPublic  ImageAccessLevel = "public"
	ImageAccessPrivate ImageAccessLevel = "private"
	ImageAccessLink    ImageAccessLevel = "link"
)

type Image struct {
	ID          ID               `json:"id" db:"id"`
	AuthorID    ID               `json:"-" db:"author_id"`
	Path        string           `json:"path" db:"path"`
	Title       string           `json:"title,omitempty" db:"title"`
	Description string           `json:"description,omitempty" db:"description"`
	AccessLevel ImageAccessLevel `json:"accessLevel" db:"access_level"`
	Mime        string           `json:"mime" db:"mime"`
	Ext         string           `json:"ext" db:"ext"`
	Url         string           `json:"url"`
	ExpiresAt   *time.Time       `json:"expiresAt,omitempty" db:"expires_at"`
	CreatedAt   time.Time        `json:"createdAt" db:"uploaded_at"`
	UpdatedAt   time.Time        `json:"updatedAt" db:"updated_at"`
}

type ImageStates struct {
	ImageID ID   `json:"imageID" db:"image_id"`
	Viewed  bool `json:"viewed" db:"viewed"`
	Liked   bool `json:"liked" db:"liked"`
}

type ImageWithAuthor struct {
	Image
	Author ImageAuthor `json:"author" db:"author"`
}

type DetailedImage struct {
	ImageWithAuthor
	Likes int        `json:"likes" db:"likes"`
	Views int        `json:"views" db:"views"`
	Tags  []ImageTag `json:"tags" db:"tags"`
}

type ImageAuthor struct {
	ID        ID     `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	AvatarURL string `json:"avatarURL" db:"avatar_url"`
}

type ImageTag struct {
	ID   ID     `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
