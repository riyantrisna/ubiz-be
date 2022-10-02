package handler

import (
	"collapp/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Router(router *gin.RouterGroup, auth middleware.AuthMiddleware) {
	users := router.Group("/users")

	users.POST("/login", h.Login)
	users.GET("/refresh-token/:userRefreshToken", h.RefreshToken)

	usersAuth := users.Group("")
	usersAuth.Use(auth.Auth())
	{
		usersAuth.GET("/", h.FindAll)
		usersAuth.GET("/:userId", h.FindById)
		usersAuth.POST("/", h.Create)
		usersAuth.PUT("/:userId", h.Update)
		usersAuth.DELETE("/:userId", h.Delete)
		usersAuth.PUT("/logout", h.Logout)
	}

}
