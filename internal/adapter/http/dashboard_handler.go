package http

import (
	"net/http"
	"time"

	"github.com/dariojcalo91/billtracker/internal/usecase"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	svc *usecase.DashboardService
}

func NewDashboardHandler(svc *usecase.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) Get(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	userID := c.GetString(userIDKey)
	summary, err := h.svc.GetDashboard(c.Request.Context(), userID, month, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
