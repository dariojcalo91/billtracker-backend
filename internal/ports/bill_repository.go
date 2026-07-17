package ports

import (
	"context"

	"github.com/dariojcalo91/billtracker/internal/domain"
)

type BillRepository interface {
	Create(ctx context.Context, bill *domain.Bill) (*domain.Bill, error)
	GetByID(ctx context.Context, id, userID string) (*domain.Bill, error)
	ListByUser(ctx context.Context, userID string) ([]*domain.Bill, error)
	Update(ctx context.Context, bill *domain.Bill) (*domain.Bill, error)
	Delete(ctx context.Context, id, userID string) error
}
