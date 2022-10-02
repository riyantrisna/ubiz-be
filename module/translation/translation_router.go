package translation

import (
	"collapp/configs"
	"collapp/module/translation/handler"
	"collapp/transport/http/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Router(db *sql.DB, router *gin.RouterGroup, cfg *configs.Config) {

	translationHandler := handler.NewTranslationHandler(db, cfg)
	translation := router.Group("/translation")

	translation.Use(middleware.Auth(db, cfg))
	{
		translation.GET("/", translationHandler.FindAll)
		translation.GET("/:translationId", translationHandler.FindById)
		translation.POST("/", translationHandler.Create)
		translation.PUT("/:translationId", translationHandler.Update)
		translation.DELETE("/:translationId", translationHandler.Delete)
	}

}
