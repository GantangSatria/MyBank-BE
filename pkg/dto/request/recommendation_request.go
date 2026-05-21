package request

// TrackRecommendationClickRequest untuk mencatat klik rekomendasi
type TrackRecommendationClickRequest struct {
	RecommendationID uint64 `json:"recommendation_id" validate:"required"`
}

// TrackFeatureClickRequest untuk mencatat klik fitur
type TrackFeatureClickRequest struct {
	FeatureName string `json:"feature_name" validate:"required,max=100"`
}
