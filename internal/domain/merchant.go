package domain

import "time"

// Merchant menyimpan data merchant/toko yang berelasi dengan transaksi.
type Merchant struct {
	ID               uint64
	MerchantID       string // e.g. "M001"
	MerchantName     string
	MerchantCategory string // F&B, E-commerce, Utilities, Travel, Entertainment, Retail, Healthcare, etc.
	MerchantCity     string
	MerchantType     string // offline, online, hybrid
	MerchantStatus   string // aktif, nonaktif
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
