package user

import (
	"collapp/configs"
	"collapp/module/user/controller"
	"collapp/transport/http/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Router(db *sql.DB, router *gin.RouterGroup, cfg *configs.Config) {

	userController := controller.NewUserController(db, cfg)
	users := router.Group("/users")

	users.POST("/login", userController.Login)
	users.GET("/refresh-token/:userRefreshToken", userController.RefreshToken)
	users.Use(middleware.Auth(db, cfg))
	{
		users.GET("/", userController.FindAll)
		users.GET("/:userId", userController.FindById)
		users.POST("/", userController.Create)
		users.PUT("/:userId", userController.Update)
		users.DELETE("/:userId", userController.Delete)
		users.PUT("/logout", userController.Logout)
	}

}
