package main

// @title MyBank API
// @version 1.0
// @description REST API for MyBank Application
// @host mybank-be.fly.dev
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"github.com/GantangSatria/MyBank-BE/config"
	"github.com/GantangSatria/MyBank-BE/internal/bootstrap"
)

func main() {
	cfg := config.Load()
	db  := bootstrap.NewDB(cfg)
	app := bootstrap.NewApp(cfg, db)
	app.Start()
}