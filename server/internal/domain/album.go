package domain

import "time"

type Album struct {
	ID          int       `json:"id"`
	AuthorID    int       `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type DetailedAlbum struct {
	Album
	Author AlbumAuthor `json:"author"`
}

type AlbumAuthor struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarURL"`
}
