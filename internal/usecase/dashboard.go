package usecase

import (
	"context"
	"time"

	"github.com/dariojcalo91/billtracker/internal/domain"
	"github.com/dariojcalo91/billtracker/internal/ports"
)

type DashboardService struct {
	billRepo    ports.BillRepository
	paymentRepo ports.PaymentRepository
}

func NewDashboardService(billRepo ports.BillRepository, paymentRepo ports.PaymentRepository) *DashboardService {
	return &DashboardService{billRepo: billRepo, paymentRepo: paymentRepo}
}

func (s *DashboardService) GetDashboard(ctx context.Context, userID, month string, ref time.Time) (*domain.DashboardSummary, error) {
	bills, err := s.billRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	payments, err := s.paymentRepo.ListByUserAndMonth(ctx, userID, month)
	if err != nil {
		return nil, err
	}

	// index payments by bill ID for O(1) lookup
	paidBills := map[string]*domain.Payment{}
	for _, p := range payments {
		paidBills[p.BillID] = p
	}

	summary := &domain.DashboardSummary{
		Month:    month,
		Done:     []*domain.BillDashboardStatus{},
		Upcoming: []*domain.BillDashboardStatus{},
		Overdue:  []*domain.BillDashboardStatus{},
	}

	today := ref.Day()

	for _, bill := range bills {
		if bill.Status != domain.BillStatusActive {
			continue
		}

		entry := &domain.BillDashboardStatus{Bill: bill}

		if payment, paid := paidBills[bill.ID]; paid {
			entry.Payment = payment
			summary.Done = append(summary.Done, entry)
		} else if bill.DueDay >= today {
			summary.Upcoming = append(summary.Upcoming, entry)
		} else {
			summary.Overdue = append(summary.Overdue, entry)
		}
	}

	summary.Summary = domain.DashboardCounts{
		Total:    len(summary.Done) + len(summary.Upcoming) + len(summary.Overdue),
		Done:     len(summary.Done),
		Upcoming: len(summary.Upcoming),
		Overdue:  len(summary.Overdue),
	}

	return summary, nil
}
