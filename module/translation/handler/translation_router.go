package handler

import (
	"collapp/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func (h *TranslationHandler) Router(router *gin.RouterGroup, auth middleware.AuthMiddleware) {
	translation := router.Group("/translation")
	translation.Use(auth.Auth())
	{
		translation.GET("/", h.FindAll)
		translation.GET("/:translationId", h.FindById)
		translation.POST("/", h.Create)
		translation.PUT("/:translationId", h.Update)
		translation.DELETE("/:translationId", h.Delete)
	}

}
