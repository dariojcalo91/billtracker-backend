package usecase

import (
	"context"
	"errors"

	"github.com/dariojcalo91/billtracker/internal/domain"
	"github.com/dariojcalo91/billtracker/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyExists = errors.New("email already registered")

type RegisterService struct {
	repo ports.UserRepository
}

func NewRegisterService(repo ports.UserRepository) *RegisterService {
	return &RegisterService{repo: repo}
}

func (s *RegisterService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(email, string(hash))
	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, user)
}
