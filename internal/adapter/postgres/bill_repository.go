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

type BillRepository struct {
	q *db.Queries
}

func NewBillRepository(pool *pgxpool.Pool) *BillRepository {
	return &BillRepository{q: db.New(pool)}
}

func (r *BillRepository) Create(ctx context.Context, b *domain.Bill) (*domain.Bill, error) {
	userUUID, err := parseUUID(b.UserID)
	if err != nil {
		return nil, err
	}

	var amount pgtype.Numeric
	if err := amount.Scan(fmt.Sprintf("%.2f", b.ExpectedAmount)); err != nil {
		return nil, err
	}

	row, err := r.q.CreateBill(ctx, db.CreateBillParams{
		UserID:          userUUID,
		Name:            b.Name,
		Category:        b.Category,
		ServiceProvider: b.ServiceProvider,
		ExpectedAmount:  amount,
		DueDay:          int16(b.DueDay),
	})
	if err != nil {
		return nil, err
	}
	return toDomainBill(row), nil
}

func (r *BillRepository) GetByID(ctx context.Context, id, userID string) (*domain.Bill, error) {
	billUUID, err := parseUUID(id)
	if err != nil {
		return nil, err
	}
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	row, err := r.q.GetBillByID(ctx, db.GetBillByIDParams{
		ID:     billUUID,
		UserID: userUUID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomainBill(row), nil
}

func (r *BillRepository) ListByUser(ctx context.Context, userID string) ([]*domain.Bill, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	rows, err := r.q.ListBillsByUser(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	bills := make([]*domain.Bill, len(rows))
	for i, row := range rows {
		bills[i] = toDomainBill(row)
	}
	return bills, nil
}

func (r *BillRepository) Update(ctx context.Context, b *domain.Bill) (*domain.Bill, error) {
	billUUID, err := parseUUID(b.ID)
	if err != nil {
		return nil, err
	}
	userUUID, err := parseUUID(b.UserID)
	if err != nil {
		return nil, err
	}

	var amount pgtype.Numeric
	if err := amount.Scan(fmt.Sprintf("%.2f", b.ExpectedAmount)); err != nil {
		return nil, err
	}

	row, err := r.q.UpdateBill(ctx, db.UpdateBillParams{
		ID:              billUUID,
		UserID:          userUUID,
		Name:            b.Name,
		Category:        b.Category,
		ServiceProvider: b.ServiceProvider,
		ExpectedAmount:  amount,
		DueDay:          int16(b.DueDay),
		Status:          string(b.Status),
	})
	if err != nil {
		return nil, err
	}
	return toDomainBill(row), nil
}

func (r *BillRepository) Delete(ctx context.Context, id, userID string) error {
	billUUID, err := parseUUID(id)
	if err != nil {
		return err
	}
	userUUID, err := parseUUID(userID)
	if err != nil {
		return err
	}

	return r.q.DeleteBill(ctx, db.DeleteBillParams{
		ID:     billUUID,
		UserID: userUUID,
	})
}

// helpers

func parseUUID(s string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}

func toDomainBill(row db.Bill) *domain.Bill {
	amountFloat, _ := row.ExpectedAmount.Float64Value()
	return &domain.Bill{
		ID:              uuid.UUID(row.ID.Bytes).String(),
		UserID:          uuid.UUID(row.UserID.Bytes).String(),
		Name:            row.Name,
		Category:        row.Category,
		ServiceProvider: row.ServiceProvider,
		ExpectedAmount:  amountFloat.Float64,
		DueDay:          int(row.DueDay),
		Status:          domain.BillStatus(row.Status),
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
}
