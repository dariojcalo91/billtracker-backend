package domain

import (
	"errors"
	"time"
)

type BillStatus string

const (
	BillStatusActive   BillStatus = "active"
	BillStatusInactive BillStatus = "inactive"
)

var (
	ErrInvalidDueDay         = errors.New("due day must be between 1 and 31")
	ErrInvalidExpectedAmount = errors.New("expected amount must be greater than zero")
	ErrBillNameRequired      = errors.New("bill name is required")
)

type Bill struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	Name            string     `json:"name"`
	Category        string     `json:"category"`
	ServiceProvider string     `json:"service_provider"`
	ExpectedAmount  float64    `json:"expected_amount"`
	DueDay          int        `json:"due_day"`
	Status          BillStatus `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func NewBill(userID, name, category, serviceProvider string, expectedAmount float64, dueDay int) (*Bill, error) {
	if name == "" {
		return nil, ErrBillNameRequired
	}
	if dueDay < 1 || dueDay > 31 {
		return nil, ErrInvalidDueDay
	}
	if expectedAmount <= 0 {
		return nil, ErrInvalidExpectedAmount
	}
	return &Bill{
		UserID:          userID,
		Name:            name,
		Category:        category,
		ServiceProvider: serviceProvider,
		ExpectedAmount:  expectedAmount,
		DueDay:          dueDay,
		Status:          BillStatusActive,
	}, nil
}
