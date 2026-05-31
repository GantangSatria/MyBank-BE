package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/GantangSatria/MyBank-BE/internal/domain"
	"github.com/GantangSatria/MyBank-BE/internal/repository"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
	"github.com/GantangSatria/MyBank-BE/pkg/mapper"
)

type RecommendationService interface {
	GetRecommendations(ctx context.Context, userID uint64) ([]response.RecommendationResponse, error)
	GenerateRecommendations(ctx context.Context, userID uint64) error
	TrackClick(ctx context.Context, userID uint64, recommendationID uint64) error
	GetReason(ctx context.Context, userID uint64, recommendationID uint64) (*response.RecommendationReasonResponse, error)
	TrackFeatureClick(ctx context.Context, userID uint64, featureName string) error
	GetFeatureClicks(ctx context.Context, userID uint64) ([]response.FeatureClickResponse, error)
}

type recommendationService struct {
	recRepo     repository.RecommendationRepository
	txRepo      repository.TransactionRepository
	fcRepo      repository.FeatureClickRepository
	userRepo    repository.UserRepository
	auditRepo   repository.AuditLogRepository
	mlServiceURL string
}

func NewRecommendationService(
	recRepo repository.RecommendationRepository,
	txRepo repository.TransactionRepository,
	fcRepo repository.FeatureClickRepository,
	userRepo repository.UserRepository,
	auditRepo repository.AuditLogRepository,
	mlServiceURL string,
) RecommendationService {
	return &recommendationService{
		recRepo:   recRepo,
		txRepo:    txRepo,
		fcRepo:       fcRepo,
		userRepo:     userRepo,
		auditRepo:    auditRepo,
		mlServiceURL: mlServiceURL,
	}
}

// GetRecommendations mengembalikan rekomendasi aktif untuk user.
// Jika belum ada, generate dulu secara otomatis.
func (s *recommendationService) GetRecommendations(ctx context.Context, userID uint64) ([]response.RecommendationResponse, error) {
	// Cek consent
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !user.IsPersonalizationEnabled {
		// Kembalikan rekomendasi generik (non-personalisasi)
		return s.getGenericRecommendations(), nil
	}

	// Cek apakah sudah ada rekomendasi aktif
	recs, err := s.recRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Jika belum ada, generate
	if len(recs) == 0 {
		if err := s.GenerateRecommendations(ctx, userID); err != nil {
			return nil, err
		}
		recs, err = s.recRepo.FindActiveByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
	}

	// Audit log
	s.auditRepo.Create(ctx, &domain.AuditLog{
		UserID: userID,
		Action: "VIEW_RECOMMENDATION",
		Detail: fmt.Sprintf("Mengakses %d rekomendasi personalisasi", len(recs)),
	})

	return mapper.MapRecommendationsToResponse(recs), nil
}

// GenerateRecommendations membuat rekomendasi berdasarkan rule-based engine.
// Rules:
// 1. Fav category → tampilkan promo terkait
// 2. Spending trend (naik/turun) → insight
// 3. Fitur yang belum dicoba → rekomendasi fitur
func (s *recommendationService) GenerateRecommendations(ctx context.Context, userID uint64) error {
	// Audit log: generate
	s.auditRepo.Create(ctx, &domain.AuditLog{
		UserID: userID,
		Action: "GENERATE_RECOMMENDATION",
		Detail: "Menganalisis data perilaku nasabah untuk personalisasi",
	})

	// Coba generate dari ML Service terlebih dahulu
	if s.mlServiceURL != "" {
		err := s.generateFromMLService(ctx, userID)
		if err == nil {
			log.Printf("[ML Service] Successfully generated recommendations for user %d", userID)
			return nil
		}
		log.Printf("[ML Service] Failed to generate recommendations: %v. Fallback to rule-based engine.", err)
	}

	priority := 1

	// === RULE 1: Rekomendasi berdasarkan kategori favorit ===
	favCat, _, err := s.txRepo.GetFavCategoryAndMethod(ctx, userID)
	if err == nil && favCat != "" {
		promoRec := s.buildCategoryPromo(userID, favCat, priority)
		s.recRepo.Create(ctx, promoRec)
		priority++
	}

	// === RULE 2: Rekomendasi berdasarkan tren spending mingguan ===
	weekly, err := s.txRepo.GetWeeklySpendingByCategory(ctx, userID)
	if err == nil && len(weekly) > 0 {
		for _, w := range weekly {
			if w.CurrentWeekSpend > w.PreviousWeekSpend*1.3 {
				// Spending naik > 30% → beri insight
				trendRec := &domain.Recommendation{
					UserID:      userID,
					Type:        domain.RecommendationTypeProduct,
					Title:       fmt.Sprintf("Spending %s Meningkat!", w.Category),
					Description: fmt.Sprintf("Pengeluaran %s Anda naik %.0f%% minggu ini. Cek promo hemat kami!", w.Category, calcPercentChange(w.PreviousWeekSpend, w.CurrentWeekSpend)),
					Reason:      fmt.Sprintf("Karena pengeluaran Anda di kategori %s meningkat signifikan minggu ini dibanding minggu lalu", w.Category),
					Priority:    priority,
					IsActive:    true,
				}
				s.recRepo.Create(ctx, trendRec)
				priority++
				break // hanya 1 insight tren
			}
		}
	}

	// === RULE 3: Rekomendasi fitur yang belum dicoba ===
	allFeatures := []string{
		"QRIS", "Savings", "Investment", "Insurance",
		"Apply Credit Card", "Apply Loan", "Cardless Withdrawal",
		"Convert to Instalment", "Promo Code", "Tukar Point Xtra", "Voucher",
		"My Schedule", "History Transaction",
	}

	clickedFeatures := make(map[string]bool)
	fcs, err := s.fcRepo.FindByUserID(ctx, userID)
	if err == nil {
		for _, fc := range fcs {
			clickedFeatures[fc.FeatureName] = true
		}
	}

	featureRecCount := 0
	for _, feature := range allFeatures {
		if featureRecCount >= 2 {
			break
		}
		if !clickedFeatures[feature] {
			featureRec := &domain.Recommendation{
				UserID:      userID,
				Type:        domain.RecommendationTypeFeature,
				Title:       fmt.Sprintf("Coba Fitur %s!", feature),
				Description: fmt.Sprintf("Anda belum pernah menggunakan fitur %s. Fitur ini bisa membantu mengelola keuangan Anda lebih baik.", feature),
				Reason:      fmt.Sprintf("Karena Anda belum pernah menggunakan fitur %s, kami rekomendasikan untuk dicoba", feature),
				Priority:    priority,
				IsActive:    true,
			}
			s.recRepo.Create(ctx, featureRec)
			priority++
			featureRecCount++
		}
	}

	// === RULE 4: Top merchants — loyalty promo ===
	merchants, err := s.txRepo.GetTopMerchants(ctx, userID, 1)
	if err == nil && len(merchants) > 0 {
		topMerchant := merchants[0]
		if topMerchant.TransactionCount >= 3 {
			loyaltyRec := &domain.Recommendation{
				UserID:      userID,
				Type:        domain.RecommendationTypePromo,
				Title:       fmt.Sprintf("Promo Spesial di %s!", topMerchant.MerchantName),
				Description: fmt.Sprintf("Anda sudah %d kali bertransaksi di %s. Dapatkan cashback spesial untuk transaksi berikutnya!", topMerchant.TransactionCount, topMerchant.MerchantName),
				Reason:      fmt.Sprintf("Karena Anda adalah pelanggan setia %s dengan %d transaksi", topMerchant.MerchantName, topMerchant.TransactionCount),
				Priority:    priority,
				IsActive:    true,
			}
			s.recRepo.Create(ctx, loyaltyRec)
		}
	}

	return nil
}

func (s *recommendationService) generateFromMLService(ctx context.Context, userID uint64) error {
	url := fmt.Sprintf("%s/recommend", s.mlServiceURL)
	payload := map[string]interface{}{
		"user_id": userID,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ML service returned status %d", resp.StatusCode)
	}

	var mlResp response.MLRecommendationResponse
	if err := json.NewDecoder(resp.Body).Decode(&mlResp); err != nil {
		return err
	}

	// Masukkan rekomendasi ke database
	priority := 1

	// Widget 1: CF Merchant -> Rekomendasi Promo / Merchant
	if len(mlResp.Rekomendasi.Widget1CFMerchant) > 0 {
		for i, w := range mlResp.Rekomendasi.Widget1CFMerchant {
			if i >= 1 { // Ambil top 1 saja
				break
			}
			rec := &domain.Recommendation{
				UserID:      userID,
				Type:        domain.RecommendationTypePromo,
				Title:       fmt.Sprintf("Rekomendasi Merchant: %s", w.MerchantName),
				Description: fmt.Sprintf("Kunjungi %s dan nikmati transaksi yang lebih mudah.", w.MerchantName),
				Reason:      w.PenjelasanXAI,
				Priority:    priority,
				IsActive:    true,
			}
			s.recRepo.Create(ctx, rec)
			priority++
		}
	}

	// Widget 2: CF Channel -> Rekomendasi Penggunaan Channel
	if len(mlResp.Rekomendasi.Widget2CFChannel) > 0 {
		for i, w := range mlResp.Rekomendasi.Widget2CFChannel {
			if i >= 1 { // Ambil top 1 saja
				break
			}
			rec := &domain.Recommendation{
				UserID:      userID,
				Type:        domain.RecommendationTypeProduct,
				Title:       fmt.Sprintf("Gunakan Channel %s", w.Channel),
				Description: fmt.Sprintf("Transaksi lebih lancar dan nyaman menggunakan %s.", w.Channel),
				Reason:      w.PenjelasanXAI,
				Priority:    priority,
				IsActive:    true,
			}
			s.recRepo.Create(ctx, rec)
			priority++
		}
	}

	// Widget 3: CBF Fitur -> Rekomendasi Fitur
	if len(mlResp.Rekomendasi.Widget3CBFFitur) > 0 {
		for i, w := range mlResp.Rekomendasi.Widget3CBFFitur {
			if i >= 1 { // Ambil top 1 saja
				break
			}
			rec := &domain.Recommendation{
				UserID:      userID,
				Type:        domain.RecommendationTypeFeature,
				Title:       fmt.Sprintf("Coba Fitur %s!", w.Feature),
				Description: fmt.Sprintf("Anda mungkin akan menyukai fitur %s untuk membantu keuangan Anda.", w.Feature),
				Reason:      w.PenjelasanXAI,
				Priority:    priority,
				IsActive:    true,
			}
			s.recRepo.Create(ctx, rec)
			priority++
		}
	}

	// Widget 4: CBF Promo -> Rekomendasi Promo Spesifik
	if len(mlResp.Rekomendasi.Widget4CBFPromo) > 0 {
		for i, w := range mlResp.Rekomendasi.Widget4CBFPromo {
			if i >= 1 { // Ambil top 1 saja
				break
			}
			rec := &domain.Recommendation{
				UserID:      userID,
				Type:        domain.RecommendationTypePromo,
				Title:       fmt.Sprintf("Promo Spesial di %s!", w.Merchant),
				Description: fmt.Sprintf("Jangan lewatkan penawaran spesial di %s.", w.Merchant),
				Reason:      w.PenjelasanXAI,
				Priority:    priority,
				IsActive:    true,
			}
			s.recRepo.Create(ctx, rec)
			priority++
		}
	}

	return nil
}

func (s *recommendationService) TrackClick(ctx context.Context, userID uint64, recommendationID uint64) error {
	// Verifikasi rekomendasi milik user
	rec, err := s.recRepo.FindByID(ctx, recommendationID)
	if err != nil {
		return err
	}
	if rec.UserID != userID {
		return apperrors.Forbidden("rekomendasi bukan milik anda")
	}

	click := &domain.RecommendationClick{
		UserID:           userID,
		RecommendationID: recommendationID,
	}

	if err := s.recRepo.LogClick(ctx, click); err != nil {
		return err
	}

	// Audit
	s.auditRepo.Create(ctx, &domain.AuditLog{
		UserID: userID,
		Action: "CLICK_RECOMMENDATION",
		Detail: fmt.Sprintf("Klik rekomendasi #%d: %s", rec.ID, rec.Title),
	})

	return nil
}

func (s *recommendationService) GetReason(ctx context.Context, userID uint64, recommendationID uint64) (*response.RecommendationReasonResponse, error) {
	rec, err := s.recRepo.FindByID(ctx, recommendationID)
	if err != nil {
		return nil, err
	}
	if rec.UserID != userID {
		return nil, apperrors.Forbidden("rekomendasi bukan milik anda")
	}

	// Cek consent status
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dataUsed := "Riwayat transaksi"
	switch rec.Type {
	case domain.RecommendationTypeFeature:
		dataUsed = "Riwayat penggunaan fitur"
	case domain.RecommendationTypePromo:
		dataUsed = "Riwayat transaksi, Preferensi kategori"
	case domain.RecommendationTypeProduct:
		dataUsed = "Riwayat transaksi, Tren pengeluaran mingguan"
	}

	// Audit
	s.auditRepo.Create(ctx, &domain.AuditLog{
		UserID: userID,
		Action: "VIEW_RECOMMENDATION_REASON",
		Detail: fmt.Sprintf("Melihat alasan rekomendasi #%d", rec.ID),
	})

	return &response.RecommendationReasonResponse{
		RecommendationID: rec.ID,
		Title:            rec.Title,
		Reason:           rec.Reason,
		DataUsed:         dataUsed,
		ConsentStatus:    user.IsPersonalizationEnabled,
	}, nil
}

func (s *recommendationService) TrackFeatureClick(ctx context.Context, userID uint64, featureName string) error {
	if err := s.fcRepo.Upsert(ctx, userID, featureName); err != nil {
		return err
	}

	// Audit
	s.auditRepo.Create(ctx, &domain.AuditLog{
		UserID: userID,
		Action: "CLICK_FEATURE",
		Detail: fmt.Sprintf("Klik fitur: %s", featureName),
	})

	return nil
}

func (s *recommendationService) GetFeatureClicks(ctx context.Context, userID uint64) ([]response.FeatureClickResponse, error) {
	fcs, err := s.fcRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapper.MapFeatureClicksToResponse(fcs), nil
}

// === Helper Functions ===

func (s *recommendationService) buildCategoryPromo(userID uint64, category string, priority int) *domain.Recommendation {
	promoMap := map[string]struct {
		title string
		desc  string
	}{
		"E-Wallet":                  {"Promo E-Wallet Spesial!", "Dapatkan cashback hingga 30% untuk top-up e-wallet favorit Anda."},
		"Food & Beverage":           {"Promo Kuliner Untukmu!", "Nikmati diskon hingga 25% di restoran dan cafe favorit Anda."},
		"Internet":                  {"Promo Paket Internet!", "Hemat hingga 20% untuk pembayaran internet bulanan Anda."},
		"Telco":                     {"Promo Pulsa & Paket Data!", "Bonus pulsa 10% untuk setiap pembelian pulsa di MyBank."},
		"Transport & Mobility":      {"Promo Transportasi!", "Cashback hingga 15% untuk pembayaran transportasi online."},
		"Retail & Convenience":      {"Promo Belanja Hemat!", "Diskon spesial hingga 20% di merchant retail favorit Anda."},
		"Lifestyle & Entertainment": {"Promo Lifestyle!", "Nikmati penawaran spesial untuk hiburan dan gaya hidup Anda."},
		"Utilities":                 {"Hemat Bayar Tagihan!", "Gratis biaya admin untuk pembayaran tagihan utilitas Anda."},
	}

	promo, ok := promoMap[category]
	if !ok {
		promo.title = fmt.Sprintf("Promo untuk %s!", category)
		promo.desc = fmt.Sprintf("Penawaran spesial untuk kategori %s favorit Anda.", category)
	}

	return &domain.Recommendation{
		UserID:      userID,
		Type:        domain.RecommendationTypePromo,
		Title:       promo.title,
		Description: promo.desc,
		Reason:      fmt.Sprintf("Karena Anda sering bertransaksi di kategori %s", category),
		Priority:    priority,
		IsActive:    true,
	}
}

func (s *recommendationService) getGenericRecommendations() []response.RecommendationResponse {
	return []response.RecommendationResponse{
		{
			ID:          0,
			Type:        "PROMO",
			Title:       "Promo Cashback untuk Semua!",
			Description: "Dapatkan cashback 10% untuk semua transaksi pertama Anda di MyBank.",
			Reason:      "Rekomendasi umum untuk semua nasabah MyBank",
			Priority:    1,
		},
		{
			ID:          0,
			Type:        "FEATURE",
			Title:       "Aktifkan Personalisasi!",
			Description: "Aktifkan fitur personalisasi untuk mendapatkan rekomendasi yang lebih relevan sesuai kebutuhan Anda.",
			Reason:      "Personalisasi belum diaktifkan. Aktifkan untuk pengalaman yang lebih baik.",
			Priority:    2,
		},
	}
}

func calcPercentChange(old, new float64) float64 {
	if old == 0 {
		return 100
	}
	return math.Round(((new - old) / old) * 100)
}
