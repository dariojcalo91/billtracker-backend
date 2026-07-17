package ports

import (
	"context"

	"github.com/dariojcalo91/billtracker/internal/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	GetByBillAndMonth(ctx context.Context, billID, month string) (*domain.Payment, error)
	ListByUserAndMonth(ctx context.Context, userID, month string) ([]*domain.Payment, error)
}
