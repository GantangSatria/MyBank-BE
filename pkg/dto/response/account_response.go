package response

import "time"

type AccountResponse struct {
	ID            uint64    `json:"id" example:"1"`
	AccountNumber string    `json:"account_number" example:"1001234567"`
	AccountType   string    `json:"account_type" example:"SAVING"`
	Balance       float64   `json:"balance" example:"150000"`
	Currency      string    `json:"currency" example:"IDR"`
	IsActive      bool      `json:"is_active" example:"true"`
	CreatedAt     time.Time `json:"created_at" example:"2023-10-01T12:00:00Z"`
}
