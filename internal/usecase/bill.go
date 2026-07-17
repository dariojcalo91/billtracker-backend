package usecase

import (
	"context"
	"errors"

	"github.com/dariojcalo91/billtracker/internal/domain"
	"github.com/dariojcalo91/billtracker/internal/ports"
)

var ErrBillNotFound = errors.New("bill not found")

type BillService struct {
	repo ports.BillRepository
}

func NewBillService(repo ports.BillRepository) *BillService {
	return &BillService{repo: repo}
}

func (s *BillService) Create(ctx context.Context, userID, name, category, serviceProvider string, expectedAmount float64, dueDay int) (*domain.Bill, error) {
	bill, err := domain.NewBill(userID, name, category, serviceProvider, expectedAmount, dueDay)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, bill)
}

func (s *BillService) GetByID(ctx context.Context, id, userID string) (*domain.Bill, error) {
	bill, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, ErrBillNotFound
	}
	return bill, nil
}

func (s *BillService) ListByUser(ctx context.Context, userID string) ([]*domain.Bill, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *BillService) Update(ctx context.Context, id, userID, name, category, serviceProvider string, expectedAmount float64, dueDay int) (*domain.Bill, error) {
	bill, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, ErrBillNotFound
	}

	// validate the new values through the domain before persisting
	updated, err := domain.NewBill(userID, name, category, serviceProvider, expectedAmount, dueDay)
	if err != nil {
		return nil, err
	}
	updated.ID = bill.ID
	updated.CreatedAt = bill.CreatedAt
	updated.Status = bill.Status

	return s.repo.Update(ctx, updated)
}

func (s *BillService) Delete(ctx context.Context, id, userID string) error {
	bill, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}
	if bill == nil {
		return ErrBillNotFound
	}
	return s.repo.Delete(ctx, id, userID)
}
