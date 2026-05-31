package bootstrap

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/GantangSatria/MyBank-BE/config"
	"github.com/GantangSatria/MyBank-BE/internal/handler"
	"github.com/GantangSatria/MyBank-BE/internal/middleware"
	"github.com/GantangSatria/MyBank-BE/internal/repository"
	"github.com/GantangSatria/MyBank-BE/internal/routes"
	"github.com/GantangSatria/MyBank-BE/internal/service"
)

type App struct {
	Fiber  *fiber.App
	Config *config.Config
	DB     *sql.DB
}

func NewApp(cfg *config.Config, db *sql.DB) *App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: customErrorHandler,
	})

	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}))

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "app": cfg.App.Name})
	})

	// di

	// Repository
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	recRepo := repository.NewRecommendationRepository(db)
	fcRepo := repository.NewFeatureClickRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)
	merchantRepo := repository.NewMerchantRepository(db)

	// Service
	authService := service.NewAuthService(userRepo, accountRepo, cfg)
	userService := service.NewUserService(userRepo)
	txService := service.NewTransactionService(txRepo, auditRepo, accountRepo)
	recService := service.NewRecommendationService(recRepo, txRepo, fcRepo, userRepo, auditRepo, cfg.App.MLServiceURL)
	merchantService := service.NewMerchantService(merchantRepo)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret, userRepo, )

	// Handler
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	txHandler := handler.NewTransactionHandler(txService)
	recHandler := handler.NewRecommendationHandler(recService)
	merchantHandler := handler.NewMerchantHandler(merchantService)

	// Routes
	routes.SetupRoutes(app, &routes.RouteConfig{
		App:                   app,
		AuthHandler:           authHandler,
		UserHandler:           userHandler,
		TransactionHandler:    txHandler,
		RecommendationHandler: recHandler,
		MerchantHandler:       merchantHandler,
		AuthMiddleware:        authMiddleware,
	})


	return &App{Fiber: app, Config: cfg, DB: db}
}

func (a *App) Start() {
	addr := fmt.Sprintf(":%s", a.Config.App.Port)
	log.Printf("[app] %s listening on %s (env: %s)", a.Config.App.Name, addr, a.Config.App.Env)
	if err := a.Fiber.Listen(addr); err != nil {
		log.Fatalf("[app] server error: %v", err)
	}
}