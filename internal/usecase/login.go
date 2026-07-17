package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dariojcalo91/billtracker/internal/ports"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type LoginService struct {
	repo      ports.UserRepository
	jwtSecret string
}

func NewLoginService(repo ports.UserRepository, jwtSecret string) *LoginService {
	return &LoginService{repo: repo, jwtSecret: jwtSecret}
}

func (s *LoginService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}
