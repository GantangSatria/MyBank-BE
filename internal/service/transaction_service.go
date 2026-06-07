package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

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
	txRepo      repository.TransactionRepository
	auditRepo   repository.AuditLogRepository
	accountRepo repository.AccountRepository
	userRepo    repository.UserRepository
}

func NewTransactionService(
	txRepo repository.TransactionRepository,
	auditRepo repository.AuditLogRepository,
	accountRepo repository.AccountRepository,
	userRepo repository.UserRepository,
) TransactionService {
	return &transactionService{
		txRepo:      txRepo,
		auditRepo:   auditRepo,
		accountRepo: accountRepo,
		userRepo:    userRepo,
	}
}
func (s *transactionService) CreateTransaction(ctx context.Context, userID uint64, req *request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	// Verify user PIN
	user, err := s.userRepo.FindByIDWithPIN(ctx, userID)
	if err != nil {
		return nil, apperrors.BadRequest("User tidak ditemukan")
	}
	if user.PIN == nil {
		return nil, apperrors.BadRequest("PIN belum diatur")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PIN), []byte(req.PIN)); err != nil {
		return nil, apperrors.BadRequest("PIN salah")
	}

	// Get source account
	account, err := s.accountRepo.FindByAccountNumber(ctx, req.AccountNumber)
	if err != nil {
		return nil, apperrors.BadRequest("Rekening sumber tidak ditemukan")
	}

	if account.UserID != userID {
		return nil, apperrors.Forbidden("Rekening ini bukan milik anda")
	}

	balanceBefore := account.Balance
	balanceAfter := balanceBefore

	// Update balance based on transaction type
	if req.Type == "TOPUP" {
		err = s.accountRepo.AddBalance(ctx, account.ID, req.Amount)
		if err != nil {
			return nil, err
		}
		balanceAfter += req.Amount
	} else if req.Type == "TRANSFER" {
		// Check if destination exists
		destAccount, err := s.accountRepo.FindByAccountNumber(ctx, req.DestinationAccountNumber)
		if err != nil {
			return nil, apperrors.BadRequest("Rekening tujuan tidak ditemukan")
		}

		err = s.accountRepo.SubtractBalance(ctx, account.ID, req.Amount)
		if err != nil {
			if err.Error() == "insufficient balance" {
				return nil, apperrors.ErrInsufficientBalance
			}
			return nil, err
		}
		balanceAfter -= req.Amount

		// Add to destination
		err = s.accountRepo.AddBalance(ctx, destAccount.ID, req.Amount)
		if err != nil {
			// Rollback logic could be implemented here for real scenario
			return nil, apperrors.InternalServerError("Gagal menambah saldo tujuan")
		}
	} else {
		return nil, apperrors.BadRequest("Tipe transaksi tidak didukung")
	}

	tx := &domain.Transaction{
		UserID:                   userID,
		AccountID:                account.ID,
		ReferenceNumber:          repository.GenerateRefNumber(req.Type),
		Type:                     domain.TransactionType(req.Type),
		Status:                   domain.TransactionStatusSuccess,
		Amount:                   req.Amount,
		BalanceBefore:            balanceBefore,
		BalanceAfter:             balanceAfter,
		DestinationAccountNumber: req.DestinationAccountNumber,
		Description:              req.Description,
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
