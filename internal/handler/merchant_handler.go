package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/service"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/request"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
	"github.com/GantangSatria/MyBank-BE/pkg/utils"
)

type MerchantHandler struct {
	svc service.MerchantService
}

func NewMerchantHandler(svc service.MerchantService) *MerchantHandler {
	return &MerchantHandler{svc: svc}
}

// CreateMerchant godoc
// @Summary Create a new merchant
// @Description Create a new merchant entry
// @Tags merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateMerchantRequest true "Create Merchant Request"
// @Success 201 {object} response.Base{data=response.MerchantResponse}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /merchants [post]
func (h *MerchantHandler) CreateMerchant(c fiber.Ctx) error {
	var req request.CreateMerchantRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	result, err := h.svc.Create(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success("merchant berhasil dibuat", result))
}

// GetMerchants godoc
// @Summary Get all merchants
// @Description Get a list of all merchants
// @Tags merchants
// @Accept json
// @Produce json
// @Success 200 {object} response.Base{data=[]response.MerchantResponse}
// @Router /merchants [get]
func (h *MerchantHandler) GetMerchants(c fiber.Ctx) error {
	result, err := h.svc.GetAll(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(response.Success("daftar merchant berhasil diambil", result))
}

// GetMerchantByID godoc
// @Summary Get merchant by merchant_id
// @Description Get a single merchant by its merchant_id (e.g. M001)
// @Tags merchants
// @Accept json
// @Produce json
// @Param merchant_id path string true "Merchant ID (e.g. M001)"
// @Success 200 {object} response.Base{data=response.MerchantResponse}
// @Failure 404 {object} response.Base
// @Router /merchants/{merchant_id} [get]
func (h *MerchantHandler) GetMerchantByID(c fiber.Ctx) error {
	merchantID := c.Params("merchant_id")
	if merchantID == "" {
		return apperrors.BadRequest("merchant_id wajib diisi")
	}

	result, err := h.svc.GetByMerchantID(c.Context(), merchantID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("merchant berhasil diambil", result))
}

// GetMerchantsByCategory godoc
// @Summary Get merchants by category
// @Description Get all active merchants filtered by category
// @Tags merchants
// @Accept json
// @Produce json
// @Param category query string true "Merchant Category (e.g. F&B, Retail)"
// @Success 200 {object} response.Base{data=[]response.MerchantResponse}
// @Router /merchants/category [get]
func (h *MerchantHandler) GetMerchantsByCategory(c fiber.Ctx) error {
	category := c.Query("category")
	if category == "" {
		return apperrors.BadRequest("parameter category wajib diisi")
	}

	result, err := h.svc.GetByCategory(c.Context(), category)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("merchant berhasil diambil", result))
}

// UpdateMerchant godoc
// @Summary Update a merchant
// @Description Update merchant details by ID
// @Tags merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Merchant DB ID"
// @Param request body request.UpdateMerchantRequest true "Update Merchant Request"
// @Success 200 {object} response.Base{data=response.MerchantResponse}
// @Failure 400 {object} response.Base
// @Failure 404 {object} response.Base
// @Router /merchants/{id} [put]
func (h *MerchantHandler) UpdateMerchant(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return apperrors.BadRequest("ID merchant tidak valid")
	}

	var req request.UpdateMerchantRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	result, err := h.svc.Update(c.Context(), id, &req)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("merchant berhasil diupdate", result))
}

// DeleteMerchant godoc
// @Summary Delete a merchant
// @Description Delete a merchant by ID
// @Tags merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Merchant DB ID"
// @Success 200 {object} response.Base
// @Failure 404 {object} response.Base
// @Router /merchants/{id} [delete]
func (h *MerchantHandler) DeleteMerchant(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return apperrors.BadRequest("ID merchant tidak valid")
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(response.Success("merchant berhasil dihapus", nil))
}
