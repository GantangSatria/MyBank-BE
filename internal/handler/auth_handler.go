package handler

import (
	"github.com/gofiber/fiber/v3"

	// "github.com/GantangSatria/MyBank-BE/internal/middleware"
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
// POST /api/v1/auth/register
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
// POST /api/v1/auth/login
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
// POST /api/v1/auth/refresh
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

// ChangePIN godoc
// PUT /api/v1/auth/pin  (protected)
// func (h *AuthHandler) ChangePIN(c fiber.Ctx) error {
// 	userID := middleware.GetUserID(c)

// 	var req request.ChangePINRequest
// 	if err := c.Bind().JSON(&req); err != nil {
// 		return apperrors.BadRequest("body request tidak valid")
// 	}
// 	if err := utils.ValidateStruct(&req); err != nil {
// 		return apperrors.BadRequest(err.Error())
// 	}

// 	if err := h.svc.ChangePIN(c.Context(), userID, &req); err != nil {
// 		return err
// 	}
// 	return c.JSON(response.Success("PIN berhasil diubah", nil))
// }