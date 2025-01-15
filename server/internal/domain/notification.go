package domain

import (
	"time"
)

type Notification struct {
	ID      ID        `json:"id" db:"id"`
	UserID  ID        `json:"-" db:"user_id"`
	Title   string    `json:"title" db:"title"`
	Message string    `json:"message" db:"message"`
	Hidden  bool      `json:"hidden" db:"hidden"`
	Read    bool      `json:"read" db:"read"`
	SentAt  time.Time `json:"sentAt" db:"sent_at"`
}

type NotificationStats struct {
	Unread int `json:"unread" db:"unread"`
}
