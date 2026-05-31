package response

type MLRecommendationResponse struct {
	UserID      int                  `json:"user_id"`
	Nama        string               `json:"nama"`
	Cluster     MLClusterInfo        `json:"cluster"`
	Rekomendasi MLRecommendationData `json:"rekomendasi"`
}

type MLClusterInfo struct {
	ID         int    `json:"id"`
	Label      string `json:"label"`
	Penjelasan string `json:"penjelasan"`
}

type MLRecommendationData struct {
	Widget1CFMerchant []MLWidgetCFMerchant `json:"widget_1_cf_merchant"`
	Widget2CFChannel  []MLWidgetCFChannel  `json:"widget_2_cf_channel"`
	Widget3CBFFitur   []MLWidgetCBFFitur   `json:"widget_3_cbf_fitur"`
	Widget4CBFPromo   []MLWidgetCBFPromo   `json:"widget_4_cbf_promo"`
}

type MLWidgetCFMerchant struct {
	MerchantName  string `json:"merchant_name"`
	PenjelasanXAI string `json:"penjelasan_xai"`
}

type MLWidgetCFChannel struct {
	Channel       string `json:"channel"`
	PenjelasanXAI string `json:"penjelasan_xai"`
}

type MLWidgetCBFFitur struct {
	Feature       string `json:"feature"`
	PenjelasanXAI string `json:"penjelasan_xai"`
}

type MLWidgetCBFPromo struct {
	Merchant      string `json:"merchant"`
	PenjelasanXAI string `json:"penjelasan_xai"`
}
