package handler

import (
	"github.com/gofiber/fiber/v3"

	// "github.com/GantangSatria/MyBank-BE/internal/middleware"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
	"github.com/GantangSatria/MyBank-BE/internal/service"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/request"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
	"github.com/GantangSatria/MyBank-BE/pkg/utils"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Register Request"
// @Success 201 {object} response.Base{data=response.TokenResponse}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Router /auth/register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	result, err := h.svc.Register(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success("registrasi berhasil", result))
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access & refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login Request"
// @Success 200 {object} response.Base{data=response.TokenResponse}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	result, err := h.svc.Login(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("login berhasil", result))
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} response.Base{data=response.TokenResponse}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	var req request.RefreshTokenRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	result, err := h.svc.RefreshToken(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("token diperbarui", result))
}

// ChangePassword godoc
// @Summary Change password
// @Description Change user password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.ChangePasswordRequest true "Change Password Request"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /auth/password [put]
func (h *AuthHandler) ChangePassword(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req request.ChangePasswordRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	if err := h.svc.ChangePassword(c.Context(), userID, &req); err != nil {
		return err
	}
	return c.JSON(response.Success("password berhasil diubah", nil))
}

// ChangePIN godoc
// @Summary Change PIN
// @Description Change user PIN
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.ChangePINRequest true "Change PIN Request"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /auth/change-pin [put]
func (h *AuthHandler) ChangePIN(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req request.ChangePINRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	if err := h.svc.ChangePIN(c.Context(), userID, &req); err != nil {
		return err
	}
	return c.JSON(response.Success("PIN berhasil diubah", nil))
}

// SetupPIN godoc
// @Summary Setup PIN
// @Description Set user PIN for the first time
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.SetPINRequest true "Set PIN Request"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /auth/pin/setup [post]
func (h *AuthHandler) SetupPIN(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req request.SetPINRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	if err := h.svc.SetupPIN(c.Context(), userID, &req); err != nil {
		return err
	}
	return c.JSON(response.Success("PIN berhasil diatur", nil))
}