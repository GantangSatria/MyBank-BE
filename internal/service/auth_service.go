package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/GantangSatria/MyBank-BE/config"
	"github.com/GantangSatria/MyBank-BE/internal/domain"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
	"github.com/GantangSatria/MyBank-BE/internal/repository"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/request"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
	mapper "github.com/GantangSatria/MyBank-BE/pkg/mapper"
)

// Skema autentikasi MyBank:
//   - Identifier  : Email ATAU Nomor HP
//   - Kredensial  : PIN 6 digit numerik (bcrypt hash, cost=12)
//   - Access Token : JWT short-lived (default 24 jam)
//   - Refresh Token: JWT long-lived (default 7 hari), subject="refresh"
type AuthService interface {
	Register(ctx context.Context, req *request.RegisterRequest) (*response.TokenResponse, error)
	Login(ctx context.Context, req *request.LoginRequest) (*response.TokenResponse, error)
	RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.TokenResponse, error)
	ChangePassword(ctx context.Context, userID uint64, req *request.ChangePasswordRequest) error
	ChangePIN(ctx context.Context, userID uint64, req *request.ChangePINRequest) error
	SetupPIN(ctx context.Context, userID uint64, req *request.SetPINRequest) error
}

type authService struct {
	userRepo    repository.UserRepository
	accountRepo repository.AccountRepository
	cfg         *config.Config
}

func NewAuthService(userRepo repository.UserRepository, accountRepo repository.AccountRepository, cfg *config.Config) AuthService {
	return &authService{userRepo: userRepo, accountRepo: accountRepo, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, req *request.RegisterRequest) (*response.TokenResponse, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.Conflict("email sudah digunakan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, apperrors.InternalServerError("gagal memproses password")
	}

	var dob *time.Time
	if req.DateOfBirth != "" {
		parsed, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			dob = &parsed
		}
	}

	user := &domain.User{
		Email:                    req.Email,
		Password:                 string(hashedPassword),
		Name:                     &req.Name,
		Phone:                    &req.Phone,
		Occupation:               &req.Occupation,
		DateOfBirth:              dob,
		IsPersonalizationEnabled: true,
		IsActive:                 true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Create random account number "100" + 7 digits
	// Unix nano is usually enough for 7 random digits pseudo-randomness for MVP
	randomPart := fmt.Sprintf("%07d", time.Now().UnixNano()%10000000)
	accountNumber := "100" + randomPart

	account := &domain.Account{
		UserID:        user.ID,
		AccountNumber: accountNumber,
		AccountType:   "Saving",
		Balance:       0,
		Currency:      "IDR",
		IsActive:      true,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		// Logically we should rollback user creation, but for MVP we return err
		return nil, err
	}

	user.Password = ""

	return s.buildTokens(user)
}

func (s *authService) Login(ctx context.Context, req *request.LoginRequest) (*response.TokenResponse, error) {

	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.Unauthorized("email atau password salah")
	}

	if !user.IsActive {
		return nil, apperrors.Forbidden("user tidak aktif")
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.Unauthorized("email atau password salah")
	}

	user.Password = ""

	return s.buildTokens(user)
}


func (s *authService) RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.TokenResponse, error) {
	claims := &middleware.Claims{}

	token, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid || claims.Subject != "refresh" {
		return nil, apperrors.Unauthorized("refresh token tidak valid")
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperrors.Unauthorized("user tidak valid")
	}

	if !user.IsActive {
		return nil, apperrors.Forbidden("user tidak aktif")
	}

	return s.buildTokens(user)
}

// func (s *authService) ChangePIN(ctx context.Context, userID uint64, req *request.ChangePINRequest) error {
// 	if req.OldPIN == req.NewPIN {
// 		return apperrors.BadRequest("PIN baru tidak boleh sama dengan PIN lama")
// 	}

// 	user, err := s.userRepo.FindByIDWithPIN(ctx, userID)
// 	if err != nil {
// 		return err
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.PIN), []byte(req.OldPIN)); err != nil {
// 		return apperrors.ErrInvalidPIN
// 	}

// 	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPIN), 12)
// 	if err != nil {
// 		return apperrors.InternalServerError("gagal memproses PIN baru")
// 	}

// 	return s.userRepo.UpdatePIN(ctx, userID, string(newHash))
// }

func (s *authService) buildTokens(user *domain.User) (*response.TokenResponse, error) {
	now := time.Now()

	accessExp := now.Add(time.Duration(s.cfg.JWT.ExpiryHour) * time.Hour)
	refreshExp := now.Add(time.Duration(s.cfg.JWT.RefreshExpiryHour) * time.Hour)

	accessToken, err := s.sign(user, "access", accessExp)
	if err != nil {
		return nil, apperrors.InternalServerError("gagal membuat access token")
	}

	refreshToken, err := s.sign(user, "refresh", refreshExp)
	if err != nil {
		return nil, apperrors.InternalServerError("gagal membuat refresh token")
	}

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Until(accessExp).Seconds()),
		User:         mapper.MapUserToResponse(user),
	}, nil
}

func (s *authService) sign(user *domain.User, subject string, exp time.Time) (string, error) {
	claims := middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "mybank",
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(s.cfg.JWT.Secret))
}

func (s *authService) ChangePassword(ctx context.Context, userID uint64, req *request.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByIDWithAuth(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return apperrors.BadRequest("password lama tidak sesuai")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		return apperrors.InternalServerError("gagal memproses password baru")
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hashed))
}

func (s *authService) ChangePIN(ctx context.Context, userID uint64, req *request.ChangePINRequest) error {
	if req.OldPIN == req.NewPIN {
		return apperrors.BadRequest("PIN baru tidak boleh sama dengan PIN lama")
	}

	user, err := s.userRepo.FindByIDWithPIN(ctx, userID)
	if err != nil {
		return err
	}

	if user.PIN == nil {
		return apperrors.BadRequest("PIN belum diatur, gunakan endpoint setup PIN")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PIN), []byte(req.OldPIN)); err != nil {
		return apperrors.ErrInvalidPIN
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPIN), 12)
	if err != nil {
		return apperrors.InternalServerError("gagal memproses PIN baru")
	}

	return s.userRepo.UpdatePIN(ctx, userID, string(hashed))
}

func (s *authService) SetupPIN(ctx context.Context, userID uint64, req *request.SetPINRequest) error {
	user, err := s.userRepo.FindByIDWithPIN(ctx, userID)
	if err != nil {
		return err
	}

	if user.PIN != nil && *user.PIN != "" {
		return apperrors.BadRequest("PIN sudah diatur, gunakan endpoint change PIN")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.PIN), 12)
	if err != nil {
		return apperrors.InternalServerError("gagal memproses PIN")
	}

	return s.userRepo.UpdatePIN(ctx, userID, string(hashed))
}