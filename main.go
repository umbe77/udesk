package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed icon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "udesk",
		Width:  800,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
		Frameless:         true,
		HideWindowOnClose: true,
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		Linux: &linux.Options{
			Icon: icon,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
