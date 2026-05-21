package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/handler"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
)

type RouteConfig struct {
	App                   *fiber.App
	AuthHandler           *handler.AuthHandler
	UserHandler	          *handler.UserHandler
	TransactionHandler    *handler.TransactionHandler
	RecommendationHandler *handler.RecommendationHandler
	AuthMiddleware        *middleware.AuthMiddleware
}

func SetupRoutes(app *fiber.App, cfg *RouteConfig) {
	app.Get("/", func (c fiber.Ctx) error {
		return c.JSON(fiber.Map{"MyBank-BE": "Pakai prefix /api/v1"})
	})

	api := app.Group("/api/v1")

	api.Get("/", func (c fiber.Ctx) error  {
		return c.JSON(fiber.Map{"MyBank-BE": "Semoga Ready to use"})
	})

	RegisterAuthRoutes(api, cfg.AuthHandler, cfg.AuthMiddleware)
	RegisterUserRoutes(api, cfg.UserHandler, *cfg.AuthMiddleware)
	RegisterTransactionRoutes(api, cfg.TransactionHandler, *cfg.AuthMiddleware)
	RegisterRecommendationRoutes(api, cfg.RecommendationHandler, *cfg.AuthMiddleware)
}