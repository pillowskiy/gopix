package domain

import (
	"time"
)

type CommentSortMethod string

const (
	CommentNewestSort CommentSortMethod = "newest"
	CommentOldestSort CommentSortMethod = "oldest"
)

type Comment struct {
	ID        ID        `json:"id" db:"id"`
	AuthorID  ID        `json:"-" db:"author_id"`
	ImageID   ID        `json:"-" db:"image_id"`
	ParentID  *ID       `json:"-" db:"parent_id"`
	Text      string    `json:"text" db:"comment"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type DetailedComment struct {
	Comment
	Stats  CommentStats  `json:"stats" db:"stats"`
	Author CommentAuthor `json:"author" db:"author"`
}

type CommentStats struct {
	RepliesCount int  `json:"repliesCount" db:"replies_count"`
	Likes        int  `json:"likes"        db:"likes"`
	Liked        bool `json:"liked"      db:"liked"`
}

type CommentAuthor struct {
	ID        ID     `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	AvatarURL string `json:"avatarURL" db:"avatar_url"`
}
