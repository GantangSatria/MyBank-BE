package request

// CreateMerchantRequest untuk membuat merchant baru
type CreateMerchantRequest struct {
	MerchantID       string `json:"merchant_id" validate:"required,max=20"`
	MerchantName     string `json:"merchant_name" validate:"required,max=200"`
	MerchantCategory string `json:"merchant_category" validate:"required,max=100"`
	MerchantCity     string `json:"merchant_city" validate:"omitempty,max=100"`
	MerchantType     string `json:"merchant_type" validate:"required,oneof=offline online hybrid"`
	MerchantStatus   string `json:"merchant_status" validate:"omitempty,oneof=aktif nonaktif"`
}

// UpdateMerchantRequest untuk update data merchant
type UpdateMerchantRequest struct {
	MerchantName     string `json:"merchant_name" validate:"omitempty,max=200"`
	MerchantCategory string `json:"merchant_category" validate:"omitempty,max=100"`
	MerchantCity     string `json:"merchant_city" validate:"omitempty,max=100"`
	MerchantType     string `json:"merchant_type" validate:"omitempty,oneof=offline online hybrid"`
	MerchantStatus   string `json:"merchant_status" validate:"omitempty,oneof=aktif nonaktif"`
}
