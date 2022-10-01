package lang

import (
	"collapp/configs"
	"collapp/module/lang/controller"
	"collapp/transport/http/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Router(db *sql.DB, router *gin.RouterGroup, cfg *configs.Config) {

	langController := controller.NewLangController(db, cfg)
	lang := router.Group("/lang")

	lang.Use(middleware.Auth(db, cfg))
	{
		lang.GET("/", langController.FindAll)
		lang.GET("/:langId", langController.FindById)
		lang.POST("/", langController.Create)
		lang.PUT("/:langId", langController.Update)
		lang.DELETE("/:langId", langController.Delete)
	}

}
