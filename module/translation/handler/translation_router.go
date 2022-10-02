package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *TranslationHandler) Router(router *gin.RouterGroup) {
	translation := router.Group("/translation")

	{
		translation.GET("/", h.FindAll)
		translation.GET("/:translationId", h.FindById)
		translation.POST("/", h.Create)
		translation.PUT("/:translationId", h.Update)
		translation.DELETE("/:translationId", h.Delete)
	}

}
