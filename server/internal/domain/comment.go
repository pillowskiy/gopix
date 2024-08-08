package domain

import "time"

type Comment struct {
	ID        int `json:"id"`
	AuthorID  int
	ImageID   int
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type DetailedComment struct {
	Comment
	Author CommentAuthor `json:"author"`
}

type CommentAuthor struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarURL"`
}
