package http

import (
	"collapp/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

// Setup will build a new router and prepare whatever the http router's need
func Setup() (*gin.Engine, *gin.RouterGroup) {
	router := gin.Default()
	router.Use(middleware.CORS())
	routerV1 := router.Group("/api/v1")

	return router, routerV1
}
