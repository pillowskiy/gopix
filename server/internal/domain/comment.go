package domain

import (
	"time"
)

type Comment struct {
	ID        int       `json:"id" db:"id"`
	AuthorID  int       `json:"-" db:"author_id"`
	ImageID   int       `json:"-" db:"image_id"`
	Text      string    `json:"text" db:"comment"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type DetailedComment struct {
	Comment
	Author CommentAuthor `json:"author" db:"author"`
}

type CommentAuthor struct {
	ID        int    `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	AvatarURL string `json:"avatarURL" db:"avatar_url"`
}
