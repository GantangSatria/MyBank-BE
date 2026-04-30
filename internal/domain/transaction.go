package domain
 
import "time"
 
type TransactionType string
 
const (
	TransactionTypeTransfer TransactionType = "TRANSFER"
	TransactionTypePayment  TransactionType = "PAYMENT"
	TransactionTypeTopUp    TransactionType = "TOPUP"
	TransactionTypeWithdraw TransactionType = "WITHDRAW"
	TransactionTypeDeposit  TransactionType = "DEPOSIT"
	TransactionTypeQRIS     TransactionType = "QRIS"
)
 
type TransactionStatus string
 
const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusSuccess   TransactionStatus = "SUCCESS"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
)
 
type Transaction struct {
	ID              uint64
	UserID          uint64
	AccountID       uint64
	ReferenceNumber string
	Type            TransactionType
	Status          TransactionStatus
	Amount          float64
	Fee             float64
	BalanceBefore   float64
	BalanceAfter    float64
 
	DestinationAccountNumber string
	DestinationBankCode      string
	DestinationName          string
 
	MerchantName     string
	MerchantCategory string
	MerchantLocation string
 
	Description string
	Note        string
	FailReason  string
 
	IsRecommended    bool
	RecommendationID *uint64
 
	TransactedAt time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}