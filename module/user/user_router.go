package user

import (
	"collapp/middleware"
	"collapp/module/user/controller"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Router(db *sql.DB, router *gin.RouterGroup) {

	userController := controller.NewUserController(db)
	users := router.Group("/users")

	users.POST("/login", userController.Login)
	users.Use(middleware.Auth(db))
	{
		users.GET("/", userController.FindAll)
		users.GET("/:userId", userController.FindById)
		users.POST("/", userController.Create)
		users.PUT("/:userId", userController.Update)
		users.DELETE("/:userId", userController.Delete)
	}

}
