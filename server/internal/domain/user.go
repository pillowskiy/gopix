package domain

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Permissions  int       `json:"permissions"`
	AvatarURL    string    `json:"avatarURL"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserWithToken struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
