package response

// RecommendationResponse adalah response untuk satu rekomendasi personalisasi
type RecommendationResponse struct {
	ID          uint64 `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url,omitempty"`
	Reason      string `json:"reason"` // explainable: "Karena Anda sering bertransaksi di E-Wallet"
	Priority    int    `json:"priority"`
}

// RecommendationReasonResponse detail "Kenapa saya melihat ini?"
type RecommendationReasonResponse struct {
	RecommendationID uint64 `json:"recommendation_id"`
	Title            string `json:"title"`
	Reason           string `json:"reason"`
	DataUsed         string `json:"data_used"` // e.g. "Riwayat transaksi, Preferensi kategori"
	ConsentStatus    bool   `json:"consent_status"`
}

// FeatureClickResponse adalah response klik fitur
type FeatureClickResponse struct {
	FeatureName string `json:"feature_name"`
	ClickCount  int64  `json:"click_count"`
	LastClicked string `json:"last_clicked,omitempty"`
}

// AuditLogResponse untuk audit trail
type AuditLogResponse struct {
	ID        uint64 `json:"id"`
	Action    string `json:"action"`
	Detail    string `json:"detail"`
	CreatedAt string `json:"created_at"`
}
