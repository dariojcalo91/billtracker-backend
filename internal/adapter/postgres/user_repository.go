package postgres

import (
	"context"
	"errors"

	"github.com/dariojcalo91/billtracker/internal/adapter/postgres/db"
	"github.com/dariojcalo91/billtracker/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{q: db.New(pool)}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	row, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	})
	if err != nil {
		return nil, err
	}
	return toDomainUser(row), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomainUser(row), nil
}

func toDomainUser(row db.User) *domain.User {
	return &domain.User{
		ID:           uuid.UUID(row.ID.Bytes).String(),
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt.Time,
	}
}
