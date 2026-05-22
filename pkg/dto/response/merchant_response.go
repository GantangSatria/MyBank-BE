package response

// MerchantResponse response detail merchant
type MerchantResponse struct {
	ID               uint64 `json:"id"`
	MerchantID       string `json:"merchant_id"`
	MerchantName     string `json:"merchant_name"`
	MerchantCategory string `json:"merchant_category"`
	MerchantCity     string `json:"merchant_city,omitempty"`
	MerchantType     string `json:"merchant_type"`
	MerchantStatus   string `json:"merchant_status"`
}
