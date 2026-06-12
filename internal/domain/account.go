package domain

import "time"

type Account struct {
	ID            uint64
	UserID        uint64
	AccountNumber string
	AccountType   string
	Balance       float64
	Currency      string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
 
type CategorySpending struct {
	MerchantCategory string
	TotalAmount      float64
	TransactionCount int64
}