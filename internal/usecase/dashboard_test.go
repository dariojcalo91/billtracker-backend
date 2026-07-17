package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/dariojcalo91/billtracker/internal/domain"
	"github.com/dariojcalo91/billtracker/internal/usecase"
)

type fakePaymentRepo struct {
	payments map[string]*domain.Payment
}

func newFakePaymentRepo() *fakePaymentRepo {
	return &fakePaymentRepo{payments: map[string]*domain.Payment{}}
}

func (f *fakePaymentRepo) Create(ctx context.Context, p *domain.Payment) (*domain.Payment, error) {
	key := p.BillID + "|" + p.Month
	f.payments[key] = p
	return p, nil
}

func (f *fakePaymentRepo) GetByBillAndMonth(ctx context.Context, billID, month string) (*domain.Payment, error) {
	p, ok := f.payments[billID+"|"+month]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (f *fakePaymentRepo) ListByUserAndMonth(ctx context.Context, userID, month string) ([]*domain.Payment, error) {
	var result []*domain.Payment
	for _, p := range f.payments {
		result = append(result, p)
	}
	return result, nil
}

func TestDashboard_GroupsBillsCorrectly(t *testing.T) {
	billRepo := newFakeBillRepo()
	paymentRepo := newFakePaymentRepo()
	svc := usecase.NewDashboardService(billRepo, paymentRepo)
	billSvc := usecase.NewBillService(billRepo)

	// due day 5 — will be overdue (reference date is day 15)
	b1, _ := billSvc.Create(context.Background(), "user-1", "Electricity", "utilities", "ACME", 100, 5)
	// due day 20 — will be upcoming (reference date is day 15)
	b2, _ := billSvc.Create(context.Background(), "user-1", "Water", "utilities", "City", 50, 20)
	// due day 10 — will be done (has a payment)
	b3, _ := billSvc.Create(context.Background(), "user-1", "Gas", "utilities", "GasCo", 75, 10)

	// mark b3 as paid
	_, _ = paymentRepo.Create(context.Background(), &domain.Payment{
		BillID:     b3.ID,
		Month:      "2026-07",
		AmountPaid: 75,
	})

	// reference date: July 15 2026
	ref := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	summary, err := svc.GetDashboard(context.Background(), "user-1", "2026-07", ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(summary.Done) != 1 {
		t.Errorf("expected 1 done, got %d", len(summary.Done))
	}
	if len(summary.Upcoming) != 1 {
		t.Errorf("expected 1 upcoming, got %d", len(summary.Upcoming))
	}
	if len(summary.Overdue) != 1 {
		t.Errorf("expected 1 overdue, got %d", len(summary.Overdue))
	}
	if summary.Summary.Total != 3 {
		t.Errorf("expected total 3, got %d", summary.Summary.Total)
	}

	// verify the done bill is b3
	if summary.Done[0].Bill.ID != b3.ID {
		t.Errorf("expected done bill to be b3, got %s", summary.Done[0].Bill.ID)
	}
	// verify upcoming bill is b2 (due day 20, after reference day 15)
	if summary.Upcoming[0].Bill.ID != b2.ID {
		t.Errorf("expected upcoming bill to be b2, got %s", summary.Upcoming[0].Bill.ID)
	}
	// verify overdue bill is b1 (due day 5, before reference day 15)
	if summary.Overdue[0].Bill.ID != b1.ID {
		t.Errorf("expected overdue bill to be b1, got %s", summary.Overdue[0].Bill.ID)
	}
}
