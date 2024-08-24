package domain

import (
	"time"
)

type Image struct {
	ID          int        `json:"id" db:"id"`
	AuthorID    int        `json:"-" db:"author_id"`
	Path        string     `json:"path" db:"path"`
	Title       string     `json:"title,omitempty" db:"title"`
	Description string     `json:"description,omitempty" db:"description"`
	AccessLevel string     `json:"accessLevel" db:"access_level"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty" db:"expires_at"`
	CreatedAt   time.Time  `json:"createdAt" db:"uploaded_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}

type DetailedImage struct {
	Image
	Author ImageAuthor `json:"author" db:"author"`
	Likes  int         `json:"likes" db:"likes"`
	Views  int         `json:"views" db:"views"`
	Tags   []ImageTag  `json:"tags" db:"tags"`
}

type ImageAuthor struct {
	ID        int    `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	AvatarURL string `json:"avatarURL" db:"avatar_url"`
}

type ImageTag struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type ImageView struct {
	ImageID int  `db:"image_id"`
	UserID  *int `db:"user_id"`
}
