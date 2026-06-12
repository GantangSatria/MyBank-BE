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
)

type AccountService interface {
	CreateAccount(ctx context.Context, userID uint64, req *request.CreateAccountRequest) (*response.AccountResponse, error)
	GetAccountsByUserID(ctx context.Context, userID uint64) ([]response.AccountResponse, error)
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, userID uint64, req *request.CreateAccountRequest) (*response.AccountResponse, error) {
	// Create random account number "100" + 7 digits
	randomPart := fmt.Sprintf("%07d", time.Now().UnixNano()%10000000)
	accountNumber := "100" + randomPart

	account := &domain.Account{
		UserID:        userID,
		AccountNumber: accountNumber,
		AccountType:   req.AccountType,
		Balance:       0, // Initial balance 0
		Currency:      "IDR",
		IsActive:      true,
	}

	err := s.accountRepo.Create(ctx, account)
	if err != nil {
		return nil, apperrors.InternalServerError("Failed to create account")
	}

	return s.mapToResponse(account), nil
}

func (s *accountService) GetAccountsByUserID(ctx context.Context, userID uint64) ([]response.AccountResponse, error) {
	accounts, err := s.accountRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var res []response.AccountResponse
	for _, acc := range accounts {
		res = append(res, *s.mapToResponse(&acc))
	}

	return res, nil
}

func (s *accountService) mapToResponse(account *domain.Account) *response.AccountResponse {
	return &response.AccountResponse{
		ID:            account.ID,
		AccountNumber: account.AccountNumber,
		AccountType:   account.AccountType,
		Balance:       account.Balance,
		Currency:      account.Currency,
		IsActive:      account.IsActive,
		CreatedAt:     account.CreatedAt,
	}
}
