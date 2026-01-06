package main

import (
	"context"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	client        *Client
	items         []ListItem
	filteredItems []ListItem
	currentAction string
}

// NewApp creates a new App application struct
func NewApp(client *Client) *App {
	return &App{
		client:        client,
		items:         []ListItem{},
		filteredItems: []ListItem{},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Hide window on startup - it will be shown via systray
	runtime.WindowHide(ctx)

	// Intercept window close event to hide instead of quit
	runtime.EventsOn(ctx, "window:close", func(optionalData ...any) {
		runtime.WindowHide(ctx)
	})
}

func (a *App) GetProcesses() {
	// get processes
	a.currentAction = "processes"
	var err error
	a.items, err = a.client.GetProcesses()
	if err != nil {
		// TODO: find a method to show an errore in communication
		runtime.LogError(a.ctx, err.Error())
	}

	runtime.EventsEmit(a.ctx, "ipc:results", a.items)
}

func (a *App) GetApplications() {
	// get applications
	a.currentAction = "applications"
	var err error
	a.items, err = a.client.GetApplications()
	if err != nil {
		// TODO: find a method to show an errore in communication
		runtime.LogError(a.ctx, err.Error())
	}

	runtime.EventsEmit(a.ctx, "ipc:results", a.items)
}

func (a *App) FilterItems(query string) []ListItem {
	items := make([]ListItem, 0)

	for _, item := range a.items {
		if fuzzy.MatchFold(query, item.Text) {
			items = append(items, item)
		}
	}
	a.filteredItems = items
	return items
}

func (a *App) SelectItem(idx int) {
	i := a.filteredItems[idx]
	a.client.selectedActions[a.currentAction](a.ctx, i)
}

func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
	runtime.WindowCenter(a.ctx)
}

func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
}
