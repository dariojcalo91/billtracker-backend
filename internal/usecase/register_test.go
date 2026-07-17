package usecase_test

import (
	"context"
	"testing"

	"github.com/dariojcalo91/billtracker/internal/domain"
	"github.com/dariojcalo91/billtracker/internal/usecase"
)

type fakeUserRepo struct {
	users map[string]*domain.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: map[string]*domain.User{}}
}

func (f *fakeUserRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	f.users[u.Email] = u
	return u, nil
}

func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, ok := f.users[email]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func TestRegister_CreatesUserWithHashedPassword(t *testing.T) {
	repo := newFakeUserRepo()
	svc := usecase.NewRegisterService(repo)

	user, err := svc.Register(context.Background(), "test@example.com", "plaintext-password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.PasswordHash == "plaintext-password" {
		t.Error("expected password to be hashed, got plaintext")
	}
}

func TestRegister_RejectsDuplicateEmail(t *testing.T) {
	repo := newFakeUserRepo()
	svc := usecase.NewRegisterService(repo)

	_, _ = svc.Register(context.Background(), "test@example.com", "password1")
	_, err := svc.Register(context.Background(), "test@example.com", "password2")

	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
}
