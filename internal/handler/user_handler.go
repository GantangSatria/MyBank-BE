package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/middleware"
	"github.com/GantangSatria/MyBank-BE/internal/service"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/request"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
	"github.com/GantangSatria/MyBank-BE/pkg/utils"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetMe godoc
// GET /api/v1/users/me
func (h *UserHandler) GetMe(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	result, err := h.svc.GetMe(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("data profil berhasil diambil", result))
}

// UpdateProfile godoc
// PATCH /api/v1/users/me
func (h *UserHandler) UpdateProfile(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req request.UpdateProfileRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	result, err := h.svc.UpdateProfile(c.Context(), userID, &req)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("profil berhasil diperbarui", result))
}

// UpdatePhone godoc
// PATCH /api/v1/users/me/phone
func (h *UserHandler) UpdatePhone(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req request.UpdatePhoneRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	result, err := h.svc.UpdatePhone(c.Context(), userID, &req)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("nomor HP berhasil diperbarui", result))
}

// UpdatePersonalizationConsent godoc
// PATCH /api/v1/users/me/personalization
func (h *UserHandler) UpdatePersonalizationConsent(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req request.UpdatePersonalizationRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	result, err := h.svc.UpdatePersonalizationConsent(c.Context(), userID, &req)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("preferensi personalisasi diperbarui", result))
}

// DeleteAccount godoc
// DELETE /api/v1/users/me
func (h *UserHandler) DeleteAccount(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	if err := h.svc.DeleteAccount(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(response.Success("akun berhasil dihapus", nil))
}