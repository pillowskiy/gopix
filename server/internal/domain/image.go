package domain

import "time"

type Image struct {
	ID          int       `json:"id"`
	AuthorID    int       `json:"-"`
	Path        string    `json:"path"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	AccessLevel string    `json:"accessLevel"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type DetailedImage struct {
	Image
	Author ImageAuthor `json:"author"`
	Likes  int         `json:"likes"`
	Views  int         `json:"views"`
	Tags   []ImageTag  `json:"tags"`
}

type ImageAuthor struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarURL"`
}

type ImageTag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
