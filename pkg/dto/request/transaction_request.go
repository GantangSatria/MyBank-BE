package request

// CreateTransactionRequest untuk membuat transaksi baru
type CreateTransactionRequest struct {
	Type                     string  `json:"type" validate:"required,oneof=TRANSFER TOPUP" example:"TRANSFER"`
	Amount                   float64 `json:"amount" validate:"required,gt=0" example:"50000"`
	AccountNumber            string  `json:"account_number" validate:"required" example:"1001234567"`
	DestinationAccountNumber string  `json:"destination_account_number" validate:"omitempty" example:"1007654321"`
	PIN                      string  `json:"pin" validate:"required,len=6,numeric" example:"123456"`
	Description              string  `json:"description" validate:"omitempty,max=255" example:"Bayar hutang"`
}


// TransactionListRequest untuk query parameter list transaksi
type TransactionListRequest struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"`
}
