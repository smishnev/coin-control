package main

import (
	"coin-control/backend/auth"
	"coin-control/backend/database"
	"coin-control/backend/user"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.DB.Close()

	// Create an instance of the app structure
	app := NewApp()
	userService := user.NewUserService()
	authService := auth.NewAuthService()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "coin-control",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			userService,
			authService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
