package lang

import (
	"collapp/middleware"
	"collapp/module/lang/controller"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Router(db *sql.DB, router *gin.RouterGroup) {

	langController := controller.NewLangController(db)
	lang := router.Group("/lang")

	lang.Use(middleware.Auth(db))
	{
		lang.GET("/", langController.FindAll)
		lang.GET("/:langId", langController.FindById)
		lang.POST("/", langController.Create)
		lang.PUT("/:langId", langController.Update)
		lang.DELETE("/:langId", langController.Delete)
	}

}
