package router

import (
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
)

// Router wraps Echo web engine for platform API routing.
type Router struct {
	Engine *echo.Echo
	server *server.Server
}

// NewRouter initializes the central router engine.
func NewRouter(s *server.Server) *Router {
	e := echo.New()
	return &Router{
		Engine: e,
		server: s,
	}
}

// PlatformRouter is an alias for Router.
type PlatformRouter = Router

func NewPlatformRouter(r *Router) *PlatformRouter {
	return r
}
