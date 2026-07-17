package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dariojcalo91/billtracker/internal/adapter/postgres/db"
	"github.com/dariojcalo91/billtracker/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	q *db.Queries
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{q: db.New(pool)}
}

func (r *PaymentRepository) Create(ctx context.Context, p *domain.Payment) (*domain.Payment, error) {
	billUUID, err := parseUUID(p.BillID)
	if err != nil {
		return nil, err
	}

	var amount pgtype.Numeric
	if err := amount.Scan(fmt.Sprintf("%.2f", p.AmountPaid)); err != nil {
		return nil, err
	}

	var proofURL pgtype.Text
	if p.ProofFileURL != nil {
		proofURL = pgtype.Text{String: *p.ProofFileURL, Valid: true}
	}

	row, err := r.q.CreatePayment(ctx, db.CreatePaymentParams{
		BillID:       billUUID,
		Month:        p.Month,
		AmountPaid:   amount,
		ProofFileUrl: proofURL,
	})
	if err != nil {
		return nil, err
	}
	return toDomainPayment(row), nil
}

func (r *PaymentRepository) GetByBillAndMonth(ctx context.Context, billID, month string) (*domain.Payment, error) {
	billUUID, err := parseUUID(billID)
	if err != nil {
		return nil, err
	}

	row, err := r.q.GetPaymentByBillAndMonth(ctx, db.GetPaymentByBillAndMonthParams{
		BillID: billUUID,
		Month:  month,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomainPayment(row), nil
}

func (r *PaymentRepository) ListByUserAndMonth(ctx context.Context, userID, month string) ([]*domain.Payment, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	rows, err := r.q.ListPaymentsByUserAndMonth(ctx, db.ListPaymentsByUserAndMonthParams{
		UserID: userUUID,
		Month:  month,
	})
	if err != nil {
		return nil, err
	}

	payments := make([]*domain.Payment, len(rows))
	for i, row := range rows {
		payments[i] = toDomainPayment(row)
	}
	return payments, nil
}

func toDomainPayment(row db.Payment) *domain.Payment {
	amountFloat, _ := row.AmountPaid.Float64Value()
	p := &domain.Payment{
		ID:         uuid.UUID(row.ID.Bytes).String(),
		BillID:     uuid.UUID(row.BillID.Bytes).String(),
		Month:      row.Month,
		AmountPaid: amountFloat.Float64,
		PaidAt:     row.PaidAt.Time,
	}
	if row.ProofFileUrl.Valid {
		p.ProofFileURL = &row.ProofFileUrl.String
	}
	return p
}
