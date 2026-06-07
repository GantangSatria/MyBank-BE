package response

type MLRecommendationResponse struct {
	UserID      int                  `json:"user_id"`
	Nama        string               `json:"nama"`
	Cluster     MLClusterInfo        `json:"cluster"`
	Rekomendasi MLRecommendationData `json:"rekomendasi"`
}

type MLClusterInfo struct {
	ID         int    `json:"id" example:"1"`
	Label      string `json:"label" example:"High Spender"`
	Penjelasan string `json:"penjelasan" example:"Cluster pengguna dengan pengeluaran tinggi"`
}

type MLRecommendationData struct {
	Widget1CFMerchant []MLWidgetCFMerchant `json:"widget_1_cf_merchant"`
	Widget2CFChannel  []MLWidgetCFChannel  `json:"widget_2_cf_channel"`
	Widget3CBFFitur   []MLWidgetCBFFitur   `json:"widget_3_cbf_fitur"`
	Widget4CBFPromo   []MLWidgetCBFPromo   `json:"widget_4_cbf_promo"`
}

type MLWidgetCFMerchant struct {
	MerchantName  string `json:"merchant_name" example:"Kopi Kenangan"`
	PenjelasanXAI string `json:"penjelasan_xai" example:"Berdasarkan kemiripan dengan pengguna lain"`
}

type MLWidgetCFChannel struct {
	Channel       string `json:"channel" example:"Mobile Banking"`
	PenjelasanXAI string `json:"penjelasan_xai" example:"Berdasarkan preferensi transaksi Anda"`
}

type MLWidgetCBFFitur struct {
	Feature       string `json:"feature" example:"Transfer"`
	PenjelasanXAI string `json:"penjelasan_xai" example:"Karena Anda sering melakukan transfer"`
}

type MLWidgetCBFPromo struct {
	Merchant      string `json:"merchant" example:"Kopi Kenangan"`
	PenjelasanXAI string `json:"penjelasan_xai" example:"Diskon 50% untuk Anda"`
}
