package domain

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Permission int

const (
	PermissionsUploadImage Permission = 1 << 0
	PermissionsAdmin       Permission = 1 << 10
)

type User struct {
	ID           ID        `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	Permissions  int       `json:"permissions" db:"permissions"`
	AvatarURL    string    `json:"avatarURL" db:"avatar_url"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	UpdatedAt    time.Time `json:"-" db:"updated_at"`
}

type DetailedUser struct {
	User
	Subscription FollowingStats `json:"subscription"`
}

type FollowingStats struct {
	Followers   int  `json:"followers" db:"followers"`
	Following   int  `json:"following" db:"following"`
	IsFollowing bool `json:"isFollowing" db:"is_following"`
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type UserPayload struct {
	ID       ID     `json:"sub"`
	Username string `json:"username"`
}

func (u *User) HasPermission(permission Permission) bool {
	return u.Permissions&int(permission) == int(permission)
}

func (u *User) PrepareMutation() error {
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
