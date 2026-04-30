package main

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