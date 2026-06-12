package domain

import "time"

// RecommendationType mendefinisikan tipe rekomendasi
type RecommendationType string

const (
	RecommendationTypePromo   RecommendationType = "PROMO"
	RecommendationTypeFeature RecommendationType = "FEATURE"
	RecommendationTypeProduct RecommendationType = "PRODUCT"
)

// Recommendation menyimpan rekomendasi yang dihasilkan engine untuk setiap user.
type Recommendation struct {
	ID          uint64
	UserID      uint64
	Type        RecommendationType
	Title       string
	Description string
	ImageURL    string
	Reason      string // explainable: "Karena Anda sering bertransaksi di E-Wallet"
	Priority    int
	IsActive    bool
	ExpiresAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RecommendationClick mencatat setiap kali user meng-klik rekomendasi — untuk menghitung CTR.
type RecommendationClick struct {
	ID               uint64
	UserID           uint64
	RecommendationID uint64
	ClickedAt        time.Time
}

// AuditLog mencatat akses/penggunaan data nasabah untuk transparansi & compliance.
type AuditLog struct {
	ID        uint64
	UserID    uint64
	Action    string // e.g. "VIEW_RECOMMENDATION", "GENERATE_RECOMMENDATION", "ACCESS_SCV"
	Detail    string // detail tambahan dalam format bebas
	IPAddress string
	CreatedAt time.Time
}

// MerchantSummary digunakan untuk aggregate top merchant dari transaksi.
type MerchantSummary struct {
	MerchantName     string
	MerchantCategory string
	TotalAmount      float64
	TransactionCount int64
}

// WeeklySpending menyimpan perbandingan spending minggu ini vs minggu lalu per kategori.
type WeeklySpending struct {
	Category          string
	CurrentWeekSpend  float64
	PreviousWeekSpend float64
}
