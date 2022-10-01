package http

import (
	"collapp/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

// Setup will build a new router and prepare whatever the http router's need
func Setup() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORS())

	return router
}
