package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/handler"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
)

func RegisterAuthRoutes(app *fiber.App, h *handler.AuthHandler, auth *middleware.AuthMiddleware) {
	v1 := app.Group("/api/v1/auth")

	v1.Post("/register", h.Register)
	v1.Post("/login", h.Login)
	v1.Post("/refresh", h.RefreshToken)

}