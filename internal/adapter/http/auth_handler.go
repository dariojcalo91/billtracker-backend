package http

import (
	"errors"
	"net/http"

	"github.com/dariojcalo91/billtracker/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	registerSvc *usecase.RegisterService
	loginSvc    *usecase.LoginService
}

func NewAuthHandler(registerSvc *usecase.RegisterService, loginSvc *usecase.LoginService) *AuthHandler {
	return &AuthHandler{registerSvc: registerSvc, loginSvc: loginSvc}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.registerSvc.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "email": user.Email})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.loginSvc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
