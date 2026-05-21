package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/middleware"
	"github.com/GantangSatria/MyBank-BE/internal/service"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/request"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
	"github.com/GantangSatria/MyBank-BE/pkg/utils"
)

type RecommendationHandler struct {
	svc service.RecommendationService
}

func NewRecommendationHandler(svc service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{svc: svc}
}

// GetRecommendations godoc
// GET /api/v1/recommendations
func (h *RecommendationHandler) GetRecommendations(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	result, err := h.svc.GetRecommendations(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("rekomendasi berhasil diambil", result))
}

// TrackClick godoc
// POST /api/v1/recommendations/:id/click
func (h *RecommendationHandler) TrackClick(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	recID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return apperrors.BadRequest("ID rekomendasi tidak valid")
	}

	if err := h.svc.TrackClick(c.Context(), userID, recID); err != nil {
		return err
	}
	return c.JSON(response.Success("klik rekomendasi tercatat", nil))
}

// GetReason godoc
// GET /api/v1/recommendations/:id/reason — "Kenapa saya melihat ini?"
func (h *RecommendationHandler) GetReason(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	recID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return apperrors.BadRequest("ID rekomendasi tidak valid")
	}

	result, err := h.svc.GetReason(c.Context(), userID, recID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("alasan rekomendasi berhasil diambil", result))
}

// TrackFeatureClick godoc
// POST /api/v1/features/click — track klik fitur mobile banking
func (h *RecommendationHandler) TrackFeatureClick(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req request.TrackFeatureClickRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	if err := h.svc.TrackFeatureClick(c.Context(), userID, req.FeatureName); err != nil {
		return err
	}
	return c.JSON(response.Success("klik fitur tercatat", nil))
}

// GetFeatureClicks godoc
// GET /api/v1/features/clicks — get semua klik fitur user
func (h *RecommendationHandler) GetFeatureClicks(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	result, err := h.svc.GetFeatureClicks(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("data klik fitur berhasil diambil", result))
}
