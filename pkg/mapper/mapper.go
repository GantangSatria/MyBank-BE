package mapper

import (
	"github.com/GantangSatria/MyBank-BE/internal/domain"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
)

func MapUserToResponse(u *domain.User) response.UserResponse {
	resp := response.UserResponse{
		ID:                       u.ID,
		Name:                     u.Name,
		Email:                    u.Email,
		Phone:                    u.Phone,
		Gender:                   u.Gender,
		Occupation:               derefStr(u.Occupation),
		MaritalStatus:            u.MaritalStatus,
		Segment:                  u.Segment,
		MonthlyIncomeRange:       u.MonthlyIncomeRange,
		IsPersonalizationEnabled: u.IsPersonalizationEnabled,
		IsActive:                 u.IsActive,
		CreatedAt:                u.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if u.DateOfBirth != nil {
		dob := u.DateOfBirth.Format("2006-01-02")
		resp.DateOfBirth = &dob
	}
	if u.LastLoginAt != nil {
		lla := u.LastLoginAt.Format("2006-01-02 15:04:05")
		resp.LastLoginAt = &lla
	}

	return resp
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}


func MapTransactionToResponse(t *domain.Transaction) response.TransactionResponse {
	return response.TransactionResponse{
		ID:                       t.ID,
		ReferenceNumber:          t.ReferenceNumber,
		Type:                     string(t.Type),
		Status:                   string(t.Status),
		Amount:                   t.Amount,
		Fee:                      t.Fee,
		BalanceBefore:            t.BalanceBefore,
		BalanceAfter:             t.BalanceAfter,
		DestinationAccountNumber: t.DestinationAccountNumber,
		DestinationBankCode:      t.DestinationBankCode,
		DestinationName:          t.DestinationName,
		MerchantName:             t.MerchantName,
		MerchantCategory:         t.MerchantCategory,
		MerchantLocation:         t.MerchantLocation,
		Channel:                  t.Channel,
		Description:              t.Description,
		Note:                     t.Note,
		IsRecommended:            t.IsRecommended,
		RecommendationID:         t.RecommendationID,
		TransactedAt:             t.TransactedAt.Format("2006-01-02 15:04:05"),
		CreatedAt:                t.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func MapTransactionsToResponse(txs []domain.Transaction) []response.TransactionResponse {
	result := make([]response.TransactionResponse, len(txs))
	for i, t := range txs {
		result[i] = MapTransactionToResponse(&t)
	}
	return result
}

func MapRecommendationToResponse(r *domain.Recommendation) response.RecommendationResponse {
	return response.RecommendationResponse{
		ID:          r.ID,
		Type:        string(r.Type),
		Title:       r.Title,
		Description: r.Description,
		ImageURL:    r.ImageURL,
		Reason:      r.Reason,
		Priority:    r.Priority,
	}
}

func MapRecommendationsToResponse(recs []domain.Recommendation) []response.RecommendationResponse {
	result := make([]response.RecommendationResponse, len(recs))
	for i, r := range recs {
		result[i] = MapRecommendationToResponse(&r)
	}
	return result
}

func MapFeatureClickToResponse(fc *domain.FeatureClick) response.FeatureClickResponse {
	return response.FeatureClickResponse{
		FeatureName: fc.FeatureName,
		ClickCount:  fc.ClickCount,
		LastClicked: fc.LastClicked.Format("2006-01-02 15:04:05"),
	}
}

func MapFeatureClicksToResponse(fcs []domain.FeatureClick) []response.FeatureClickResponse {
	result := make([]response.FeatureClickResponse, len(fcs))
	for i, fc := range fcs {
		result[i] = MapFeatureClickToResponse(&fc)
	}
	return result
}

func MapAuditLogToResponse(al *domain.AuditLog) response.AuditLogResponse {
	return response.AuditLogResponse{
		ID:        al.ID,
		Action:    al.Action,
		Detail:    al.Detail,
		CreatedAt: al.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func MapMerchantToResponse(m *domain.Merchant) response.MerchantResponse {
	return response.MerchantResponse{
		ID:               m.ID,
		MerchantID:       m.MerchantID,
		MerchantName:     m.MerchantName,
		MerchantCategory: m.MerchantCategory,
		MerchantCity:     m.MerchantCity,
		MerchantType:     m.MerchantType,
		MerchantStatus:   m.MerchantStatus,
	}
}

func MapMerchantsToResponse(merchants []domain.Merchant) []response.MerchantResponse {
	result := make([]response.MerchantResponse, len(merchants))
	for i, m := range merchants {
		result[i] = MapMerchantToResponse(&m)
	}
	return result
}
