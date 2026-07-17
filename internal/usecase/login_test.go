package usecase_test

import (
	"context"
	"testing"

	"github.com/dariojcalo91/billtracker/internal/domain"
	"github.com/dariojcalo91/billtracker/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_Success(t *testing.T) {
	repo := newFakeUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	user, _ := domain.NewUser("test@example.com", string(hash))
	_, _ = repo.Create(context.Background(), user)

	svc := usecase.NewLoginService(repo, "test-secret")

	token, err := svc.Login(context.Background(), "test@example.com", "correct-password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected a non-empty token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := newFakeUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	user, _ := domain.NewUser("test@example.com", string(hash))
	_, _ = repo.Create(context.Background(), user)

	svc := usecase.NewLoginService(repo, "test-secret")

	_, err := svc.Login(context.Background(), "test@example.com", "wrong-password")
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	repo := newFakeUserRepo()
	svc := usecase.NewLoginService(repo, "test-secret")

	_, err := svc.Login(context.Background(), "ghost@example.com", "whatever")
	if err == nil {
		t.Error("expected error for unknown email, got nil")
	}
}
