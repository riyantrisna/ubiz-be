package translation

import (
	"collapp/middleware"
	"collapp/module/translation/controller"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Router(db *sql.DB, router *gin.RouterGroup) {

	translationController := controller.NewTranslationController(db)
	translation := router.Group("/translation")

	translation.Use(middleware.Auth(db))
	{
		translation.GET("/", translationController.FindAll)
		translation.GET("/:translationId", translationController.FindById)
		translation.POST("/", translationController.Create)
		translation.PUT("/:translationId", translationController.Update)
		translation.DELETE("/:translationId", translationController.Delete)
	}

}
