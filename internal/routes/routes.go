package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/handler"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
)

type RouteConfig struct {
	App            *fiber.App
	AuthHandler    *handler.AuthHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func SetupRoutes(app *fiber.App, cfg *RouteConfig) {
	api := app.Group("/api/v1")

	api.Get("", func (c fiber.Ctx) error  {
		return c.JSON(fiber.Map{"MyBank-BE": "Semoga Ready to use"})
	})

	RegisterAuthRoutes(api, cfg.AuthHandler, cfg.AuthMiddleware)

}