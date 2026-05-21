package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/handler"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
)

func RegisterTransactionRoutes(router fiber.Router, h *handler.TransactionHandler, auth middleware.AuthMiddleware) {
	tx := router.Group("/transactions", auth.Protect())

	tx.Post("", h.CreateTransaction)
	tx.Get("", h.GetTransactions)
	tx.Get("/summary", h.GetSpendingSummary)
	tx.Get("/:id", h.GetTransactionDetail)
}
