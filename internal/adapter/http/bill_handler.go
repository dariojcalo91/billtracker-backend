package http

import (
	"errors"
	"net/http"

	"github.com/dariojcalo91/billtracker/internal/usecase"
	"github.com/gin-gonic/gin"
)

type BillHandler struct {
	svc *usecase.BillService
}

func NewBillHandler(svc *usecase.BillService) *BillHandler {
	return &BillHandler{svc: svc}
}

type billRequest struct {
	Name            string  `json:"name" binding:"required"`
	Category        string  `json:"category" binding:"required"`
	ServiceProvider string  `json:"service_provider" binding:"required"`
	ExpectedAmount  float64 `json:"expected_amount" binding:"required,gt=0"`
	DueDay          int     `json:"due_day" binding:"required,min=1,max=31"`
}

func (h *BillHandler) Create(c *gin.Context) {
	var req billRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString(userIDKey)
	bill, err := h.svc.Create(c.Request.Context(), userID, req.Name, req.Category, req.ServiceProvider, req.ExpectedAmount, req.DueDay)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, bill)
}

func (h *BillHandler) List(c *gin.Context) {
	userID := c.GetString(userIDKey)
	bills, err := h.svc.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, bills)
}

func (h *BillHandler) Get(c *gin.Context) {
	userID := c.GetString(userIDKey)
	bill, err := h.svc.GetByID(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrBillNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "bill not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, bill)
}

func (h *BillHandler) Update(c *gin.Context) {
	var req billRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString(userIDKey)
	bill, err := h.svc.Update(c.Request.Context(), c.Param("id"), userID, req.Name, req.Category, req.ServiceProvider, req.ExpectedAmount, req.DueDay)
	if err != nil {
		if errors.Is(err, usecase.ErrBillNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "bill not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bill)
}

func (h *BillHandler) Delete(c *gin.Context) {
	userID := c.GetString(userIDKey)
	err := h.svc.Delete(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrBillNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "bill not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
