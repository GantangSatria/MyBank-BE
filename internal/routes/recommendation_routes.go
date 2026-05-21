package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/handler"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
)

func RegisterRecommendationRoutes(router fiber.Router, h *handler.RecommendationHandler, auth middleware.AuthMiddleware) {
	rec := router.Group("/recommendations", auth.Protect())

	rec.Get("", h.GetRecommendations)
	rec.Post("/:id/click", h.TrackClick)
	rec.Get("/:id/reason", h.GetReason)

	// Feature click tracking
	features := router.Group("/features", auth.Protect())
	features.Post("/click", h.TrackFeatureClick)
	features.Get("/clicks", h.GetFeatureClicks)
}
