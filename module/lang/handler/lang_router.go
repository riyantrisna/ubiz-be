package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *LangHandler) Router(router *gin.RouterGroup) {
	lang := router.Group("/lang")

	{
		lang.GET("/", h.FindAll)
		lang.GET("/:langId", h.FindById)
		lang.POST("/", h.Create)
		lang.PUT("/:langId", h.Update)
		lang.DELETE("/:langId", h.Delete)
	}
}
