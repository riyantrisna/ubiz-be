package controller

import (
	"collapp/configs"
	"collapp/helper"
	"collapp/module/lang/model"
	"collapp/module/lang/service"
	translationService "collapp/module/translation/service"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LangControllerImpl struct {
	LangService        service.LangService
	Validate           *validator.Validate
	TranslationService translationService.TranslationService
	config             *configs.Config
}

func NewLangController(db *sql.DB, cfg *configs.Config) LangController {
	validate := validator.New()
	langService := service.NewLangService(db)
	translationService := translationService.NewTranslationService(db)
	return &LangControllerImpl{
		LangService:        langService,
		Validate:           validate,
		TranslationService: translationService,
		config:             cfg,
	}
}

func (controller *LangControllerImpl) Create(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	langCreateRequest := model.LangCreateRequest{}
	context.Bind(&langCreateRequest)

	langCreateRequest.CreatedBy = payloadJwt.UserId

	currentTime := time.Now()
	langCreateRequest.CreatedAt = currentTime.Format("2006-01-02 15:04:05")

	err := controller.Validate.Struct(langCreateRequest)
	if err != nil {
		webResponse := helper.WebResponse{
			Code:   http.StatusBadRequest,
			Status: controller.TranslationService.Translation(context, "bad_request", payloadJwt.UserLangCode),
			Data:   err.Error(),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	langResponse := controller.LangService.Create(context, langCreateRequest)
	webResponse := helper.WebResponse{
		Code:   200,
		Status: controller.TranslationService.Translation(context, "success_create_language", payloadJwt.UserLangCode),
		Data:   langResponse,
	}

	context.Writer.Header().Add("Content-Type", "application/json")
	context.JSON(200, webResponse)
}

func (controller *LangControllerImpl) Update(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	langUpdateRequest := model.LangUpdateRequest{}
	context.Bind(&langUpdateRequest)

	langUpdateRequest.UpdatedBy = payloadJwt.UserId

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
			Status: controller.TranslationService.Translation(context, "bad_request", payloadJwt.UserLangCode),
			Data:   err.Error(),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	fmt.Println(langUpdateRequest)

	langResponse := controller.LangService.Update(context, langUpdateRequest)

	if langResponse.LangId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_update_language", payloadJwt.UserLangCode),
			Data:   langResponse,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: controller.TranslationService.Translation(context, "data_not_found", payloadJwt.UserLangCode),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (controller *LangControllerImpl) Delete(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	langId := context.Param("langId")
	id, err := strconv.Atoi(langId)
	helper.PanicIfError(err)

	langResponse := controller.LangService.Delete(context, id)

	if langResponse.LangId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_delete_language", payloadJwt.UserLangCode),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: controller.TranslationService.Translation(context, "data_not_found", payloadJwt.UserLangCode),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (controller *LangControllerImpl) FindById(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	langId := context.Param("langId")
	id, err := strconv.Atoi(langId)
	helper.PanicIfError(err)

	langResponse := controller.LangService.FindById(context, id)

	if langResponse.LangId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_get_language", payloadJwt.UserLangCode),
			Data:   langResponse,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: controller.TranslationService.Translation(context, "data_not_found", payloadJwt.UserLangCode),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (controller *LangControllerImpl) FindAll(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	langResponses := controller.LangService.FindAll(context)

	if len(langResponses) > 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_get_language", payloadJwt.UserLangCode),
			Data:   langResponses,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: controller.TranslationService.Translation(context, "data_not_found", payloadJwt.UserLangCode),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}
