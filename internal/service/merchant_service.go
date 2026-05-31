package service

import (
	"context"

	"github.com/GantangSatria/MyBank-BE/internal/domain"
	"github.com/GantangSatria/MyBank-BE/internal/repository"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/request"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
	"github.com/GantangSatria/MyBank-BE/pkg/mapper"
)

type MerchantService interface {
	Create(ctx context.Context, req *request.CreateMerchantRequest) (*response.MerchantResponse, error)
	GetAll(ctx context.Context) ([]response.MerchantResponse, error)
	GetByMerchantID(ctx context.Context, merchantID string) (*response.MerchantResponse, error)
	GetByCategory(ctx context.Context, category string) ([]response.MerchantResponse, error)
	Update(ctx context.Context, id uint64, req *request.UpdateMerchantRequest) (*response.MerchantResponse, error)
	Delete(ctx context.Context, id uint64) error
}

type merchantService struct {
	repo repository.MerchantRepository
}

func NewMerchantService(repo repository.MerchantRepository) MerchantService {
	return &merchantService{repo: repo}
}

func (s *merchantService) Create(ctx context.Context, req *request.CreateMerchantRequest) (*response.MerchantResponse, error) {
	status := req.MerchantStatus
	if status == "" {
		status = "aktif"
	}

	m := &domain.Merchant{
		MerchantID:       req.MerchantID,
		MerchantName:     req.MerchantName,
		MerchantCategory: req.MerchantCategory,
		MerchantCity:     req.MerchantCity,
		MerchantType:     req.MerchantType,
		MerchantStatus:   status,
	}

	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}

	result := mapper.MapMerchantToResponse(m)
	return &result, nil
}

func (s *merchantService) GetAll(ctx context.Context) ([]response.MerchantResponse, error) {
	merchants, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.MapMerchantsToResponse(merchants), nil
}

func (s *merchantService) GetByMerchantID(ctx context.Context, merchantID string) (*response.MerchantResponse, error) {
	m, err := s.repo.FindByMerchantID(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	result := mapper.MapMerchantToResponse(m)
	return &result, nil
}

func (s *merchantService) GetByCategory(ctx context.Context, category string) ([]response.MerchantResponse, error) {
	merchants, err := s.repo.FindByCategory(ctx, category)
	if err != nil {
		return nil, err
	}
	return mapper.MapMerchantsToResponse(merchants), nil
}

func (s *merchantService) Update(ctx context.Context, id uint64, req *request.UpdateMerchantRequest) (*response.MerchantResponse, error) {
	m := &domain.Merchant{
		ID:               id,
		MerchantName:     req.MerchantName,
		MerchantCategory: req.MerchantCategory,
		MerchantCity:     req.MerchantCity,
		MerchantType:     req.MerchantType,
		MerchantStatus:   req.MerchantStatus,
	}

	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}

	result := mapper.MapMerchantToResponse(m)
	return &result, nil
}

func (s *merchantService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
