package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Router(router *gin.RouterGroup) {
	users := router.Group("/users")

	users.POST("/login", h.Login)
	users.GET("/refresh-token/:userRefreshToken", h.RefreshToken)
	// users.Use(authMiddleware.Auth())
	{
		users.GET("/", h.FindAll)
		users.GET("/:userId", h.FindById)
		users.POST("/", h.Create)
		users.PUT("/:userId", h.Update)
		users.DELETE("/:userId", h.Delete)
		users.PUT("/logout", h.Logout)
	}

}
