package request

// CreateTransactionRequest untuk membuat transaksi baru
type CreateTransactionRequest struct {
	Type                     string  `json:"type" validate:"required,oneof=TRANSFER PAYMENT TOPUP WITHDRAW DEPOSIT QRIS"`
	Amount                   float64 `json:"amount" validate:"required,gt=0"`
	AccountID                uint64  `json:"account_id" validate:"required"`
	PIN                      string  `json:"pin" validate:"required,len=6,numeric"`
	DestinationAccountNumber string  `json:"destination_account_number" validate:"omitempty"`
	DestinationBankCode      string  `json:"destination_bank_code" validate:"omitempty"`
	DestinationName          string  `json:"destination_name" validate:"omitempty"`
	MerchantName             string  `json:"merchant_name" validate:"omitempty,max=200"`
	MerchantCategory         string  `json:"merchant_category" validate:"omitempty,max=100"`
	MerchantLocation         string  `json:"merchant_location" validate:"omitempty,max=200"`
	Channel                  string  `json:"channel" validate:"omitempty,oneof=mobile QRIS virtual_account transfer"`
	Description              string  `json:"description" validate:"omitempty,max=255"`
	Note                     string  `json:"note" validate:"omitempty,max=255"`
}


// TransactionListRequest untuk query parameter list transaksi
type TransactionListRequest struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"`
}
