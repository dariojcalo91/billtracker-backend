package http_test

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/dariojcalo91/billtracker/internal/adapter/http"
	"github.com/dariojcalo91/billtracker/internal/domain"
	"github.com/dariojcalo91/billtracker/internal/usecase"
	"github.com/gin-gonic/gin"
)

type fakeUserRepo struct {
	users map[string]*domain.User
}

func (f *fakeUserRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	if f.users == nil {
		f.users = map[string]*domain.User{}
	}
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

func TestRegisterHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &fakeUserRepo{}
	svc := usecase.NewRegisterService(repo)
	handler := httpadapter.NewAuthHandler(svc, nil)

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	body := []byte(`{"email":"test@example.com","password":"supersecret"}`)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestRegisterHandler_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &fakeUserRepo{}
	svc := usecase.NewRegisterService(repo)
	handler := httpadapter.NewAuthHandler(svc, nil)

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	body := []byte(`{"email":"not-an-email","password":"supersecret"}`)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected status 400, got %d, body: %s", w.Code, w.Body.String())
	}
}
