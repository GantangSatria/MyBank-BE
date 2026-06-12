package response

// TransactionResponse adalah response untuk satu transaksi
type TransactionResponse struct {
	ID              uint64  `json:"id" example:"1"`
	ReferenceNumber string  `json:"reference_number" example:"TRX-TOPUP-1634567890"`
	Type            string  `json:"type" example:"TOPUP"`
	Status          string  `json:"status" example:"SUCCESS"`
	Amount          float64 `json:"amount" example:"50000"`
	Fee             float64 `json:"fee" example:"0"`
	BalanceBefore   float64 `json:"balance_before" example:"100000"`
	BalanceAfter    float64 `json:"balance_after" example:"150000"`

	DestinationAccountNumber string `json:"destination_account_number,omitempty" example:"1007654321"`
	DestinationBankCode      string `json:"destination_bank_code,omitempty" example:"BCA"`
	DestinationName          string `json:"destination_name,omitempty" example:"John Doe"`

	MerchantName     string `json:"merchant_name,omitempty" example:"Kopi Kenangan"`
	MerchantCategory string `json:"merchant_category,omitempty" example:"F&B"`
	MerchantLocation string `json:"merchant_location,omitempty" example:"Jakarta"`

	Channel     string `json:"channel,omitempty" example:"Mobile Banking"`
	Description string `json:"description,omitempty" example:"Topup saldo"`
	Note        string `json:"note,omitempty" example:"Bayar hutang"`

	IsRecommended    bool    `json:"is_recommended" example:"false"`
	RecommendationID *uint64 `json:"recommendation_id,omitempty" example:"1"`

	TransactedAt string `json:"transacted_at" example:"2023-10-01 12:00:00"`
	CreatedAt    string `json:"created_at" example:"2023-10-01 12:00:00"`
}


// CategorySpendingResponse menunjukkan total spending per kategori merchant
type CategorySpendingResponse struct {
	MerchantCategory string  `json:"merchant_category" example:"F&B"`
	TotalAmount      float64 `json:"total_amount" example:"150000"`
	TransactionCount int64   `json:"transaction_count" example:"3"`
}

// MerchantSummaryResponse menunjukkan top merchant
type MerchantSummaryResponse struct {
	MerchantName     string  `json:"merchant_name" example:"Kopi Kenangan"`
	MerchantCategory string  `json:"merchant_category" example:"F&B"`
	TotalAmount      float64 `json:"total_amount" example:"50000"`
	TransactionCount int64   `json:"transaction_count" example:"1"`
}

// WeeklySpendingResponse perbandingan spending minggu ini vs minggu lalu
type WeeklySpendingResponse struct {
	Category          string  `json:"category" example:"F&B"`
	CurrentWeekSpend  float64 `json:"current_week_spend" example:"150000"`
	PreviousWeekSpend float64 `json:"previous_week_spend" example:"100000"`
}

// SpendingSummaryResponse ringkasan spending keseluruhan user
type SpendingSummaryResponse struct {
	TotalSpend        float64                    `json:"total_spend" example:"500000"`
	TotalTransactions int64                      `json:"total_transactions" example:"10"`
	FavCategory       string                     `json:"fav_category" example:"F&B"`
	FavMethod         string                     `json:"fav_method" example:"QRIS"`
	CategoryBreakdown []CategorySpendingResponse `json:"category_breakdown"`
	TopMerchants      []MerchantSummaryResponse  `json:"top_merchants"`
	WeeklyComparison  []WeeklySpendingResponse   `json:"weekly_comparison"`
}
