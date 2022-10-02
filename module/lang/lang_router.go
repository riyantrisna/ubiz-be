package lang

import (
	"collapp/configs"
	"collapp/module/lang/handler"
	"collapp/transport/http/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Router(db *sql.DB, router *gin.RouterGroup, cfg *configs.Config) {

	langHandler := handler.NewLangHandler(db, cfg)
	lang := router.Group("/lang")

	lang.Use(middleware.Auth(db, cfg))
	{
		lang.GET("/", langHandler.FindAll)
		lang.GET("/:langId", langHandler.FindById)
		lang.POST("/", langHandler.Create)
		lang.PUT("/:langId", langHandler.Update)
		lang.DELETE("/:langId", langHandler.Delete)
	}

}
