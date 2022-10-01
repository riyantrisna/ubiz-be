package translation

import (
	"collapp/configs"
	"collapp/module/translation/controller"
	"collapp/transport/http/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Router(db *sql.DB, router *gin.RouterGroup, cfg *configs.Config) {

	translationController := controller.NewTranslationController(db, cfg)
	translation := router.Group("/translation")

	translation.Use(middleware.Auth(db, cfg))
	{
		translation.GET("/", translationController.FindAll)
		translation.GET("/:translationId", translationController.FindById)
		translation.POST("/", translationController.Create)
		translation.PUT("/:translationId", translationController.Update)
		translation.DELETE("/:translationId", translationController.Delete)
	}

}
