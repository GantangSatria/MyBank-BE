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
// @Summary Get user recommendations
// @Description Get personalized recommendations for the logged-in user
// @Tags recommendations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base{data=[]response.RecommendationResponse}
// @Failure 401 {object} response.Base
// @Router /recommendations [get]
func (h *RecommendationHandler) GetRecommendations(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	result, err := h.svc.GetRecommendations(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("rekomendasi berhasil diambil", result))
}

// TrackClick godoc
// @Summary Track recommendation click
// @Description Record a click on a specific recommendation
// @Tags recommendations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Recommendation ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 403 {object} response.Base
// @Router /recommendations/{id}/click [post]
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
// @Summary Get recommendation reason
// @Description Get the reason why a specific recommendation was shown to the user
// @Tags recommendations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Recommendation ID"
// @Success 200 {object} response.Base{data=response.RecommendationReasonResponse}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 403 {object} response.Base
// @Router /recommendations/{id}/reason [get]
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
// @Summary Track feature click
// @Description Record a user's click on a mobile banking feature
// @Tags features
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.TrackFeatureClickRequest true "Feature Click Request"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /features/click [post]
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
// @Summary Get feature clicks
// @Description Get a list of all features clicked by the user
// @Tags features
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base{data=[]response.FeatureClickResponse}
// @Failure 401 {object} response.Base
// @Router /features/clicks [get]
func (h *RecommendationHandler) GetFeatureClicks(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	result, err := h.svc.GetFeatureClicks(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("data klik fitur berhasil diambil", result))
}
