package response

// MerchantResponse response detail merchant
type MerchantResponse struct {
	ID               uint64 `json:"id" example:"1"`
	MerchantID       string `json:"merchant_id" example:"M001"`
	MerchantName     string `json:"merchant_name" example:"Kopi Kenangan"`
	MerchantCategory string `json:"merchant_category" example:"F&B"`
	MerchantCity     string `json:"merchant_city,omitempty" example:"Jakarta"`
	MerchantType     string `json:"merchant_type" example:"offline"`
	MerchantStatus   string `json:"merchant_status" example:"aktif"`
}
