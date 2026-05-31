package service

import (
	"context"
	"fmt"
	"time"

	"github.com/GantangSatria/MyBank-BE/internal/domain"
	"github.com/GantangSatria/MyBank-BE/internal/repository"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/request"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
	"github.com/GantangSatria/MyBank-BE/pkg/mapper"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, userID uint64, req *request.CreateTransactionRequest) (*response.TransactionResponse, error)
	GetTransactions(ctx context.Context, userID uint64, page, limit int) ([]response.TransactionResponse, int64, error)
	GetTransactionDetail(ctx context.Context, userID uint64, txID uint64) (*response.TransactionResponse, error)
	GetSpendingSummary(ctx context.Context, userID uint64) (*response.SpendingSummaryResponse, error)
}

type transactionService struct {
	txRepo       repository.TransactionRepository
	auditRepo    repository.AuditLogRepository
}

func NewTransactionService(
	txRepo repository.TransactionRepository,
	auditRepo repository.AuditLogRepository,
) TransactionService {
	return &transactionService{
		txRepo:    txRepo,
		auditRepo: auditRepo,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, userID uint64, req *request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	tx := &domain.Transaction{
		UserID:                   userID,
		AccountID:                req.AccountID,
		ReferenceNumber:          repository.GenerateRefNumber(req.Type),
		Type:                     domain.TransactionType(req.Type),
		Status:                   domain.TransactionStatusSuccess,
		Amount:                   req.Amount,
		DestinationAccountNumber: req.DestinationAccountNumber,
		DestinationBankCode:      req.DestinationBankCode,
		DestinationName:          req.DestinationName,
		MerchantName:             req.MerchantName,
		MerchantCategory:         req.MerchantCategory,
		MerchantLocation:         req.MerchantLocation,
		Channel:                  req.Channel,
		Description:              req.Description,
		Note:                     req.Note,
	}


	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	// Audit log
	s.auditRepo.Create(ctx, &domain.AuditLog{
		UserID: userID,
		Action: "CREATE_TRANSACTION",
		Detail: "Transaksi " + req.Type + " sebesar " + formatAmount(req.Amount),
	})

	result := mapper.MapTransactionToResponse(tx)
	return &result, nil
}

func (s *transactionService) GetTransactions(ctx context.Context, userID uint64, page, limit int) ([]response.TransactionResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	txs, total, err := s.txRepo.FindByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return mapper.MapTransactionsToResponse(txs), total, nil
}

func (s *transactionService) GetTransactionDetail(ctx context.Context, userID uint64, txID uint64) (*response.TransactionResponse, error) {
	tx, err := s.txRepo.FindByID(ctx, txID)
	if err != nil {
		return nil, err
	}

	if tx.UserID != userID {
		return nil, apperrors.Forbidden("transaksi bukan milik anda")
	}

	result := mapper.MapTransactionToResponse(tx)
	return &result, nil
}

func (s *transactionService) GetSpendingSummary(ctx context.Context, userID uint64) (*response.SpendingSummaryResponse, error) {
	// Total spend & count
	totalSpend, totalTrx, err := s.txRepo.GetTotalSpendAndCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Fav category & method
	favCat, favMethod, err := s.txRepo.GetFavCategoryAndMethod(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Category breakdown (all time)
	now := time.Now()
	yearAgo := now.AddDate(-1, 0, 0)
	categories, err := s.txRepo.GetCategorySpending(ctx, userID, yearAgo, now)
	if err != nil {
		return nil, err
	}

	catResp := make([]response.CategorySpendingResponse, len(categories))
	for i, c := range categories {
		catResp[i] = response.CategorySpendingResponse{
			MerchantCategory: c.MerchantCategory,
			TotalAmount:      c.TotalAmount,
			TransactionCount: c.TransactionCount,
		}
	}

	// Top merchants
	merchants, err := s.txRepo.GetTopMerchants(ctx, userID, 3)
	if err != nil {
		return nil, err
	}

	merchResp := make([]response.MerchantSummaryResponse, len(merchants))
	for i, m := range merchants {
		merchResp[i] = response.MerchantSummaryResponse{
			MerchantName:     m.MerchantName,
			MerchantCategory: m.MerchantCategory,
			TotalAmount:      m.TotalAmount,
			TransactionCount: m.TransactionCount,
		}
	}

	// Weekly comparison
	weekly, err := s.txRepo.GetWeeklySpendingByCategory(ctx, userID)
	if err != nil {
		return nil, err
	}

	weeklyResp := make([]response.WeeklySpendingResponse, len(weekly))
	for i, w := range weekly {
		weeklyResp[i] = response.WeeklySpendingResponse{
			Category:          w.Category,
			CurrentWeekSpend:  w.CurrentWeekSpend,
			PreviousWeekSpend: w.PreviousWeekSpend,
		}
	}

	// Audit log
	s.auditRepo.Create(ctx, &domain.AuditLog{
		UserID: userID,
		Action: "VIEW_SPENDING_SUMMARY",
		Detail: "Akses ringkasan spending untuk personalisasi",
	})

	return &response.SpendingSummaryResponse{
		TotalSpend:        totalSpend,
		TotalTransactions: totalTrx,
		FavCategory:       favCat,
		FavMethod:         favMethod,
		CategoryBreakdown: catResp,
		TopMerchants:      merchResp,
		WeeklyComparison:  weeklyResp,
	}, nil
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("Rp%.0f", amount)
}
