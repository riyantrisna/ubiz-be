package user

import (
	"collapp/configs"
	"collapp/module/user/handler"
	"collapp/transport/http/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Router(db *sql.DB, router *gin.RouterGroup, cfg *configs.Config) {

	userHandler := handler.NewUserHandler(db, cfg)
	users := router.Group("/users")

	users.POST("/login", userHandler.Login)
	users.GET("/refresh-token/:userRefreshToken", userHandler.RefreshToken)
	users.Use(middleware.Auth(db, cfg))
	{
		users.GET("/", userHandler.FindAll)
		users.GET("/:userId", userHandler.FindById)
		users.POST("/", userHandler.Create)
		users.PUT("/:userId", userHandler.Update)
		users.DELETE("/:userId", userHandler.Delete)
		users.PUT("/logout", userHandler.Logout)
	}

}
