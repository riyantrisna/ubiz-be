package handler

import (
	"collapp/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func (h *LangHandler) Router(router *gin.RouterGroup, auth middleware.AuthMiddleware) {
	lang := router.Group("/lang")
	lang.Use(auth.Auth())
	{
		lang.GET("/", h.FindAll)
		lang.GET("/:langId", h.FindById)
		lang.POST("/", h.Create)
		lang.PUT("/:langId", h.Update)
		lang.DELETE("/:langId", h.Delete)
	}
}
