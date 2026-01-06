package main

import (
	"embed"

	"fyne.io/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed icon.png
var icon []byte

func initSysTray(app *App) (func(), func()) {

	start, end := systray.RunWithExternalLoop(
		func() { // onReady
			systray.SetIcon(icon)
			systray.SetTitle("udesk")

			showItem := systray.AddMenuItem("Show", "")
			go func() {
				for range showItem.ClickedCh {
					app.ShowWindow()
				}
			}()
			systray.AddSeparator()
			quitItem := systray.AddMenuItem("Quit", "")
			go func() {
				for range quitItem.ClickedCh {
					runtime.Quit(app.ctx)
				}
			}()
		}, func() { // onExit

		})

	return start, end
}

func main() {
	client := NewClient()

	// Create an instance of the app structure
	app := NewApp(client)

	start, end := initSysTray(app)
	defer end()

	start()
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
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "0ed7f376-9fb1-42f4-a4bf-71c84f3d3606",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				app.ShowWindow()
			},
		},
		Bind: []any{
			app,
		},
		Frameless:         true,
		LogLevel:          logger.DEBUG,
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
