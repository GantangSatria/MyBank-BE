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
// @Summary Get current user profile
// @Description Get profile of the currently logged in user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base{data=response.UserResponse}
// @Failure 401 {object} response.Base
// @Router /users/me [get]
func (h *UserHandler) GetMe(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	result, err := h.svc.GetMe(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("data profil berhasil diambil", result))
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update profile details of the currently logged in user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateProfileRequest true "Update Profile Request"
// @Success 200 {object} response.Base{data=response.UserResponse}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /users/me [patch]
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
// @Summary Update user phone number
// @Description Update phone number of the currently logged in user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdatePhoneRequest true "Update Phone Request"
// @Success 200 {object} response.Base{data=response.UserResponse}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 409 {object} response.Base
// @Router /users/me/phone [patch]
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
// @Summary Update personalization consent
// @Description Enable or disable data personalization for the user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdatePersonalizationRequest true "Update Personalization Consent Request"
// @Success 200 {object} response.Base{data=response.UserResponse}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /users/me/personalization [patch]
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
// @Summary Soft delete user account
// @Description Soft deletes the currently logged in user account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /users/me [delete]
func (h *UserHandler) DeleteAccount(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	if err := h.svc.DeleteAccount(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(response.Success("akun berhasil dihapus", nil))
}