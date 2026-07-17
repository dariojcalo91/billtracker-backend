package domain

import "time"

type Payment struct {
	ID           string
	BillID       string    `json:"bill_id"`
	Month        string    `json:"month"`
	AmountPaid   float64   `json:"amount_paid"`
	ProofFileURL *string   `json:"proof_file_url"`
	PaidAt       time.Time `json:"paid_at"`
}
