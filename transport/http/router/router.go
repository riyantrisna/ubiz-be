package router

import (
	langHandler "collapp/module/lang/handler"
	translationHandler "collapp/module/translation/handler"
	userHandler "collapp/module/user/handler"
	"collapp/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

// ModuleHandlers is a struct that contains all module-specific handlers.
type ModuleHandlers struct {
	UserHandler        userHandler.UserHandler
	LangHandler        langHandler.LangHandler
	TranslationHandler translationHandler.TranslationHandler
}

// Router is the router struct containing handlers.
type Router struct {
	ModuleHandlers ModuleHandlers
}

// NewRouter is the provider function for this router.
func NewRouter(handlers ModuleHandlers) Router {
	return Router{
		ModuleHandlers: handlers,
	}
}

// SetupRoutes sets up all routing for this server.
func (r *Router) SetupRoutes(routerGroup *gin.RouterGroup, auth middleware.AuthMiddleware) {
	r.ModuleHandlers.UserHandler.Router(routerGroup, auth)
	r.ModuleHandlers.TranslationHandler.Router(routerGroup, auth)
	r.ModuleHandlers.LangHandler.Router(routerGroup, auth)
}
