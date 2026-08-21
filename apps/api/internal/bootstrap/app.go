package bootstrap

import (
	"context"

	"github.com/golangnigeria/curexal/internal/kernel/router"
)

// App is the main application composition root runner.
type App struct {
	Container *Container
}

func NewApp(container *Container) *App {
	return &App{
		Container: container,
	}
}

func (a *App) Start() error {
	return a.Container.Server.Start()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.Container.Server.Shutdown(ctx)
}

func (a *App) Router() *router.Router {
	return a.Container.Router
}
