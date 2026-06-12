package response

// RecommendationResponse adalah response untuk satu rekomendasi personalisasi
type RecommendationResponse struct {
	ID          uint64 `json:"id" example:"1"`
	Type        string `json:"type" example:"PROMO"`
	Title       string `json:"title" example:"Diskon 50%"`
	Description string `json:"description" example:"Diskon khusus untuk Anda"`
	ImageURL    string `json:"image_url,omitempty" example:"https://example.com/image.jpg"`
	Reason      string `json:"reason" example:"Karena Anda sering bertransaksi di E-Wallet"`
	Priority    int    `json:"priority" example:"1"`
}

// RecommendationReasonResponse detail "Kenapa saya melihat ini?"
type RecommendationReasonResponse struct {
	RecommendationID uint64 `json:"recommendation_id" example:"1"`
	Title            string `json:"title" example:"Diskon 50%"`
	Reason           string `json:"reason" example:"Karena Anda sering bertransaksi di E-Wallet"`
	DataUsed         string `json:"data_used" example:"Riwayat transaksi, Preferensi kategori"`
	ConsentStatus    bool   `json:"consent_status" example:"true"`
}

// FeatureClickResponse adalah response klik fitur
type FeatureClickResponse struct {
	FeatureName string `json:"feature_name" example:"Transfer"`
	ClickCount  int64  `json:"click_count" example:"5"`
	LastClicked string `json:"last_clicked,omitempty" example:"2023-10-01 12:00:00"`
}

// AuditLogResponse untuk audit trail
type AuditLogResponse struct {
	ID        uint64 `json:"id" example:"1"`
	Action    string `json:"action" example:"LOGIN"`
	Detail    string `json:"detail" example:"User berhasil login"`
	CreatedAt string `json:"created_at" example:"2023-10-01 12:00:00"`
}
