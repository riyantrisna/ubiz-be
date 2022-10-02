package handler

import (
	"github.com/gin-gonic/gin"
)

type LangHandler interface {
	Create(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
	FindById(context *gin.Context)
	FindAll(context *gin.Context)
}
