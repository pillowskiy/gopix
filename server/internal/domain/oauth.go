package domain

import "time"

type OAuthService string

const (
	OAuthServiceGoogle OAuthService = "google"
)

type OAuth struct {
	UserID    ID           `json:"userID" db:"user_id"`
	OAuthID   string       `json:"-" db:"oauth_id"`
	Service   OAuthService `json:"-" db:"service"`
	CreatedAt time.Time    `json:"-" db:"created_at"`
	UpdatedAt time.Time    `json:"-" db:"updated_at"`
}

type OAuthUser struct {
	ID      string       `json:"id"`
	Email   string       `json:"email"`
	Name    string       `json:"name"`
	Picture string       `json:"picture"`
	Service OAuthService `json:"service"`
}
