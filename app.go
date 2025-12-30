package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Hide window on startup - it will be shown via systray
	// runtime.WindowHide(ctx)

	// Intercept window close event to hide instead of quit
	runtime.EventsOn(ctx, "window:close", func(optionalData ...any) {
		runtime.WindowHide(ctx)
	})
}

func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
	runtime.WindowCenter(a.ctx)
}

func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
}
