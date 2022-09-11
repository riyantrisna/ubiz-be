package controller

import (
	"collapp/helper"
	"collapp/module/lang/model"
	"collapp/module/lang/service"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LangControllerImpl struct {
	LangService service.LangService
	Validate    *validator.Validate
}

func NewLangController(db *sql.DB) LangController {
	validate := validator.New()
	langService := service.NewLangService(db)
	return &LangControllerImpl{
		LangService: langService,
		Validate:    validate,
	}
}

func (controller *LangControllerImpl) Create(context *gin.Context) {
	langCreateRequest := model.LangCreateRequest{}
	context.Bind(&langCreateRequest)

	value, ok := context.Get("user_id")
	if ok {
		langCreateRequest.CreatedBy = value.(int)
	} else {
		langCreateRequest.CreatedBy = 0
	}

	currentTime := time.Now()
	langCreateRequest.CreatedAt = currentTime.Format("2006-01-02 15:04:05")

	err := controller.Validate.Struct(langCreateRequest)
	if err != nil {
		webResponse := helper.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   err.Error(),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	langResponse := controller.LangService.Create(context.Request.Context(), langCreateRequest)
	webResponse := helper.WebResponse{
		Code:   200,
		Status: "Success create language",
		Data:   langResponse,
	}

	context.Writer.Header().Add("Content-Type", "application/json")
	context.JSON(200, webResponse)
}

func (controller *LangControllerImpl) Update(context *gin.Context) {
	langUpdateRequest := model.LangUpdateRequest{}
	context.Bind(&langUpdateRequest)

	value, ok := context.Get("user_id")
	if ok {
		langUpdateRequest.UpdatedBy = value.(int)
	} else {
		langUpdateRequest.UpdatedBy = 0
	}

	currentTime := time.Now()
	langUpdateRequest.UpdatedAt = currentTime.Format("2006-01-02 15:04:05")

	langId := context.Param("langId")
	id, err := strconv.Atoi(langId)
	helper.PanicIfError(err)

	langUpdateRequest.LangId = id

	err = controller.Validate.Struct(langUpdateRequest)
	if err != nil {
		webResponse := helper.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   err.Error(),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	fmt.Println(langUpdateRequest)

	langResponse := controller.LangService.Update(context.Request.Context(), langUpdateRequest)

	if langResponse.LangId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: "Success update language",
			Data:   langResponse,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Data not found",
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (controller *LangControllerImpl) Delete(context *gin.Context) {
	langId := context.Param("langId")
	id, err := strconv.Atoi(langId)
	helper.PanicIfError(err)

	langResponse := controller.LangService.Delete(context.Request.Context(), id)

	if langResponse.LangId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: "Success delete language",
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Data not found",
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (controller *LangControllerImpl) FindById(context *gin.Context) {
	langId := context.Param("langId")
	id, err := strconv.Atoi(langId)
	helper.PanicIfError(err)

	langResponse := controller.LangService.FindById(context.Request.Context(), id)

	if langResponse.LangId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: "Success get language",
			Data:   langResponse,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Data not found",
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (controller *LangControllerImpl) FindAll(context *gin.Context) {
	langResponses := controller.LangService.FindAll(context.Request.Context())

	if len(langResponses) > 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: "Success get all langs",
			Data:   langResponses,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Data not found",
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}
