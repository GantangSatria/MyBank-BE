package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/handler"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
)

func RegisterMerchantRoutes(router fiber.Router, h *handler.MerchantHandler, auth middleware.AuthMiddleware) {
	merchants := router.Group("/merchants")

	// Public: list & detail
	merchants.Get("", h.GetMerchants)
	merchants.Get("/category", h.GetMerchantsByCategory)
	merchants.Get("/:merchant_id", h.GetMerchantByID)

	// Protected: CUD
	merchants.Post("", auth.Protect(), h.CreateMerchant)
	merchants.Put("/:id", auth.Protect(), h.UpdateMerchant)
	merchants.Delete("/:id", auth.Protect(), h.DeleteMerchant)
}
