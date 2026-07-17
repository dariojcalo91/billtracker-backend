package domain

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

var ErrInvalidEmail = errors.New("invalid email")

func NewUser(email, passwordHash string) (*User, error) {
	if !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}
	return &User{Email: email, PasswordHash: passwordHash}, nil
}
