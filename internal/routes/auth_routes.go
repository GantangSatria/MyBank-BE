package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/handler"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
)

func RegisterAuthRoutes(router fiber.Router,  h *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	auth := router.Group("/auth")

	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.RefreshToken)

	// Protected routes
	protected := auth.Group("", authMiddleware.Protect())
	protected.Put("/password", h.ChangePassword)
	protected.Put("/change-pin", h.ChangePIN)
	protected.Post("/pin/setup", h.SetupPIN)

}