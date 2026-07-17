package ports

import (
	"context"

	"github.com/dariojcalo91/billtracker/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
