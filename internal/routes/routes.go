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

func SetupRoutes(cfg *RouteConfig) {
	RegisterAuthRoutes(cfg.App, cfg.AuthHandler, cfg.AuthMiddleware)

}