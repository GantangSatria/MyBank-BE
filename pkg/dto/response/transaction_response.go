package response

// TransactionResponse adalah response untuk satu transaksi
type TransactionResponse struct {
	ID              uint64  `json:"id"`
	ReferenceNumber string  `json:"reference_number"`
	Type            string  `json:"type"`
	Status          string  `json:"status"`
	Amount          float64 `json:"amount"`
	Fee             float64 `json:"fee"`
	BalanceBefore   float64 `json:"balance_before"`
	BalanceAfter    float64 `json:"balance_after"`

	DestinationAccountNumber string `json:"destination_account_number,omitempty"`
	DestinationBankCode      string `json:"destination_bank_code,omitempty"`
	DestinationName          string `json:"destination_name,omitempty"`

	MerchantName     string `json:"merchant_name,omitempty"`
	MerchantCategory string `json:"merchant_category,omitempty"`
	MerchantLocation string `json:"merchant_location,omitempty"`

	Description string `json:"description,omitempty"`
	Note        string `json:"note,omitempty"`

	IsRecommended    bool    `json:"is_recommended"`
	RecommendationID *uint64 `json:"recommendation_id,omitempty"`

	TransactedAt string `json:"transacted_at"`
	CreatedAt    string `json:"created_at"`
}

// CategorySpendingResponse menunjukkan total spending per kategori merchant
type CategorySpendingResponse struct {
	MerchantCategory string  `json:"merchant_category"`
	TotalAmount      float64 `json:"total_amount"`
	TransactionCount int64   `json:"transaction_count"`
}

// MerchantSummaryResponse menunjukkan top merchant
type MerchantSummaryResponse struct {
	MerchantName     string  `json:"merchant_name"`
	MerchantCategory string  `json:"merchant_category"`
	TotalAmount      float64 `json:"total_amount"`
	TransactionCount int64   `json:"transaction_count"`
}

// WeeklySpendingResponse perbandingan spending minggu ini vs minggu lalu
type WeeklySpendingResponse struct {
	Category          string  `json:"category"`
	CurrentWeekSpend  float64 `json:"current_week_spend"`
	PreviousWeekSpend float64 `json:"previous_week_spend"`
}

// SpendingSummaryResponse ringkasan spending keseluruhan user
type SpendingSummaryResponse struct {
	TotalSpend        float64                    `json:"total_spend"`
	TotalTransactions int64                      `json:"total_transactions"`
	FavCategory       string                     `json:"fav_category"`
	FavMethod         string                     `json:"fav_method"`
	CategoryBreakdown []CategorySpendingResponse `json:"category_breakdown"`
	TopMerchants      []MerchantSummaryResponse  `json:"top_merchants"`
	WeeklyComparison  []WeeklySpendingResponse   `json:"weekly_comparison"`
}
