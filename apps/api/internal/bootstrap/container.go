package bootstrap

import (
	"github.com/golangnigeria/curexal/internal/kernel/database"
	"github.com/golangnigeria/curexal/internal/kernel/router"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/golangnigeria/curexal/internal/shared/logger"
)

// Container is the central application composition root struct holding global dependencies.
type Container struct {
	Config   *config.Config
	Logger   *logger.Logger
	Database *database.DB
	Server   *server.Server
	Router   *router.Router
}

// NewContainer initializes and composes the global dependency container.
func NewContainer(cfg *config.Config, log *logger.Logger, db *database.DB, srv *server.Server, rtr *router.Router) *Container {
	return &Container{
		Config:   cfg,
		Logger:   log,
		Database: db,
		Server:   srv,
		Router:   rtr,
	}
}
