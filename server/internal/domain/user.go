package domain

import "time"

type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	Permissions  int       `json:"permissions" db:"permissions"`
	AvatarURL    string    `json:"avatarURL" db:"avatar_url"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type UserWithToken struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
