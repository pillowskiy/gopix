package domain

import "time"

type Album struct {
	ID          ID        `json:"id" db:"id"`
	AuthorID    ID        `json:"-" db:"author_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type DetailedAlbum struct {
	Album
	Cover  []Image     `json:"cover" db:"cover"`
	Author AlbumAuthor `json:"author" db:"author"`
}

type AlbumAuthor struct {
	ID        ID     `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	AvatarURL string `json:"avatarURL" db:"avatar_url"`
}
