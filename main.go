package main

import (
	"coin-control/backend/auth"
	"coin-control/backend/bybit"
	"coin-control/backend/database"
	"coin-control/backend/user"
	"context"
	"embed"
	"log"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Load .env if present (dev convenience). Ignore errors silently.
	_ = godotenv.Load()
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.DB.Close()

	// Create an instance of the app structure
	app := NewApp()
	userService := user.NewUserService()
	authService := auth.NewAuthService()
	bybitService := bybit.NewBybitService()

	// Create application with options
	err := wails.Run(&options.App{
		Title:    "coin-control",
		Width:    1024,
		Height:   768,
		LogLevel: logger.ERROR,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			// pass runtime ctx to bybit package for EventsEmit
			bybit.SetRuntimeCtx(ctx)
		},
		Bind: []interface{}{
			app,
			userService,
			authService,
			bybitService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
