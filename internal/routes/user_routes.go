package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/handler"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
)

func RegisterUserRoutes(router fiber.Router, h *handler.UserHandler, auth middleware.AuthMiddleware){
	user := router.Group("/users/me", auth.Protect())

	user.Get("", h.GetMe)
	user.Patch("", h.UpdateProfile)
	user.Patch("/phone", h.UpdatePhone)
	user.Patch("/personalization", h.UpdatePersonalizationConsent)
	user.Delete("", h.DeleteAccount)
}