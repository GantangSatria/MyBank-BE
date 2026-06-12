package service

import (
	"context"

	"github.com/GantangSatria/MyBank-BE/internal/domain"
	"github.com/GantangSatria/MyBank-BE/internal/repository"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/request"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
	"github.com/GantangSatria/MyBank-BE/pkg/mapper"
)

type UserService interface {
	GetMe(ctx context.Context, userID uint64) (*response.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint64, req *request.UpdateProfileRequest) (*response.UserResponse, error)
	UpdatePhone(ctx context.Context, userID uint64, req *request.UpdatePhoneRequest) (*response.UserResponse, error)
	UpdatePersonalizationConsent(ctx context.Context, userID uint64, req *request.UpdatePersonalizationRequest) (*response.UserResponse, error)
	DeleteAccount(ctx context.Context, userID uint64) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetMe(ctx context.Context, userID uint64) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := mapper.MapUserToResponse(user)
	return &res, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint64, req *request.UpdateProfileRequest) (*response.UserResponse, error) {
	u := &domain.User{
		Name:        req.Name,
		Occupation:  req.Occupation,
		DateOfBirth: req.ParsedDateOfBirth(),
	}

	if err := s.userRepo.UpdateProfile(ctx, userID, u); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := mapper.MapUserToResponse(user)
	return &res, nil
}

func (s *userService) UpdatePhone(ctx context.Context, userID uint64, req *request.UpdatePhoneRequest) (*response.UserResponse, error) {
	exists, err := s.userRepo.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.Conflict("nomor HP sudah digunakan")
	}

	u := &domain.User{Phone: &req.Phone}
	if err := s.userRepo.UpdatePhone(ctx, userID, u); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := mapper.MapUserToResponse(user)
	return &res, nil
}

func (s *userService) UpdatePersonalizationConsent(ctx context.Context, userID uint64, req *request.UpdatePersonalizationRequest) (*response.UserResponse, error) {
	if err := s.userRepo.UpdatePersonalizationConsent(ctx, userID, req.Enabled); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := mapper.MapUserToResponse(user)
	return &res, nil
}

func (s *userService) DeleteAccount(ctx context.Context, userID uint64) error {
	return s.userRepo.SoftDelete(ctx, userID)
}