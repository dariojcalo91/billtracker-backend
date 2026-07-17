package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/dariojcalo91/billtracker/internal/domain"
	"github.com/dariojcalo91/billtracker/internal/usecase"
)

type fakeBillRepo struct {
	bills map[string]*domain.Bill
}

func newFakeBillRepo() *fakeBillRepo {
	return &fakeBillRepo{bills: map[string]*domain.Bill{}}
}

func (f *fakeBillRepo) Create(ctx context.Context, b *domain.Bill) (*domain.Bill, error) {
	b.ID = fmt.Sprintf("fake-id-%d", len(f.bills)+1)
	f.bills[b.ID] = b
	return b, nil
}

func (f *fakeBillRepo) GetByID(ctx context.Context, id, userID string) (*domain.Bill, error) {
	b, ok := f.bills[id]
	if !ok || b.UserID != userID {
		return nil, nil
	}
	return b, nil
}

func (f *fakeBillRepo) ListByUser(ctx context.Context, userID string) ([]*domain.Bill, error) {
	var result []*domain.Bill
	for _, b := range f.bills {
		if b.UserID == userID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (f *fakeBillRepo) Update(ctx context.Context, b *domain.Bill) (*domain.Bill, error) {
	f.bills[b.ID] = b
	return b, nil
}

func (f *fakeBillRepo) Delete(ctx context.Context, id, userID string) error {
	delete(f.bills, id)
	return nil
}

func TestBill_CreateSuccess(t *testing.T) {
	repo := newFakeBillRepo()
	svc := usecase.NewBillService(repo)

	bill, err := svc.Create(context.Background(), "user-1", "Electricity", "utilities", "ACME Power", 120.50, 15)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if bill.ID == "" {
		t.Error("expected a generated ID")
	}
	if bill.Status != domain.BillStatusActive {
		t.Errorf("expected status active, got %s", bill.Status)
	}
}

func TestBill_CreateInvalidDueDay(t *testing.T) {
	repo := newFakeBillRepo()
	svc := usecase.NewBillService(repo)

	_, err := svc.Create(context.Background(), "user-1", "Electricity", "utilities", "ACME Power", 120.50, 45)
	if err == nil {
		t.Error("expected error for invalid due day, got nil")
	}
}

func TestBill_CreateInvalidAmount(t *testing.T) {
	repo := newFakeBillRepo()
	svc := usecase.NewBillService(repo)

	_, err := svc.Create(context.Background(), "user-1", "Electricity", "utilities", "ACME Power", -10, 15)
	if err == nil {
		t.Error("expected error for negative amount, got nil")
	}
}

func TestBill_ListByUser(t *testing.T) {
	repo := newFakeBillRepo()
	svc := usecase.NewBillService(repo)

	_, _ = svc.Create(context.Background(), "user-1", "Electricity", "utilities", "ACME Power", 120.50, 15)
	_, _ = svc.Create(context.Background(), "user-1", "Water", "utilities", "City Water", 45.00, 10)

	bills, err := svc.ListByUser(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(bills) != 2 {
		t.Errorf("expected 2 bills, got %d", len(bills))
	}
}

func TestBill_DeleteSuccess(t *testing.T) {
	repo := newFakeBillRepo()
	svc := usecase.NewBillService(repo)

	bill, _ := svc.Create(context.Background(), "user-1", "Electricity", "utilities", "ACME Power", 120.50, 15)
	err := svc.Delete(context.Background(), bill.ID, "user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
