package domain

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type UserPayload struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func (u *User) PreCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.PasswordHash = strings.TrimSpace(u.PasswordHash)

	if err := u.hashPassword(); err != nil {
		return err
	}

	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

func (u *User) HidePassword() {
	u.PasswordHash = ""
}

func (u *User) hashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}
