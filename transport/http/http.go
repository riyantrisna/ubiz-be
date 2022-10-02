package http

import (
	"collapp/configs"
	"collapp/transport/http/middleware"
	"collapp/transport/http/router"

	"github.com/gin-gonic/gin"
)

// HTTP is the HTTP server.
type HTTP struct {
	Config         *configs.Config
	Router         router.Router
	AuthMiddleware middleware.AuthMiddleware
	routerEngine   *gin.Engine
}

// NewHTTP is the provider for HTTP.
func NewHTTP(config *configs.Config, router router.Router, authMiddleware middleware.AuthMiddleware) *HTTP {
	return &HTTP{
		Config:         config,
		Router:         router,
		AuthMiddleware: authMiddleware,
	}
}

func (h *HTTP) setupRoutes() {
	h.routerEngine = gin.Default()
	routerV1 := h.routerEngine.Group("/api/v1")
	h.Router.SetupRoutes(routerV1, h.AuthMiddleware)
}

func (h *HTTP) setupMiddleware() {
	h.routerEngine.Use(middleware.CORS())
}

// SetupAndServe will build a new router and prepare whatever the http router's need
func (h *HTTP) SetupAndServe() {
	h.setupRoutes()
	h.setupMiddleware()

	h.routerEngine.Run(h.Config.Address)
}
