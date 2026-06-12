package domain

import "time"

// FeatureClick menyimpan jumlah klik fitur mobile banking per user.
// Kolom ini memetakan 13 fitur dari SCV: Apply Credit Card, QRIS, Savings, dll.
type FeatureClick struct {
	ID          uint64
	UserID      uint64
	FeatureName string
	ClickCount  int64
	LastClicked time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
