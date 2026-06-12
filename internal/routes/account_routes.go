package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/GantangSatria/MyBank-BE/internal/handler"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
)

func RegisterAccountRoutes(router fiber.Router, accountHandler *handler.AccountHandler, auth middleware.AuthMiddleware) {
	accounts := router.Group("/accounts", auth.Protect())

	accounts.Post("/", accountHandler.CreateAccount)
	accounts.Get("/", accountHandler.GetAccounts)
}
