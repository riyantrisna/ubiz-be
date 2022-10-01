package controller

import (
	"collapp/configs"
	"collapp/helper"
	"collapp/module/translation/model"
	"collapp/module/translation/service"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TranslationControllerImpl struct {
	TranslationService service.TranslationService
	Validate           *validator.Validate
	config             *configs.Config
}

func NewTranslationController(db *sql.DB, cfg *configs.Config) TranslationController {
	validate := validator.New()
	translationService := service.NewTranslationService(db)
	return &TranslationControllerImpl{
		TranslationService: translationService,
		Validate:           validate,
		config:             cfg,
	}
}

func (controller *TranslationControllerImpl) Create(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	translationCreateRequest := model.TranslationCreateRequest{}
	context.Bind(&translationCreateRequest)

	translationCreateRequest.CreatedBy = payloadJwt.UserId

	currentTime := time.Now()
	translationCreateRequest.CreatedAt = currentTime.Format("2006-01-02 15:04:05")

	keyIsExist := controller.TranslationService.CheckKeyTranslationExist(context, translationCreateRequest.TranslationKey)
	if keyIsExist {
		webResponse := helper.WebResponse{
			Code:   http.StatusBadRequest,
			Status: controller.TranslationService.Translation(context, "bad_request", payloadJwt.UserLangCode),
			Data:   controller.TranslationService.Translation(context, "key_translation_is_exist", payloadJwt.UserLangCode) + " (" + translationCreateRequest.TranslationKey + ")",
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	err := controller.Validate.Struct(translationCreateRequest)
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

	translationResponse := controller.TranslationService.Create(context, translationCreateRequest)
	webResponse := helper.WebResponse{
		Code:   200,
		Status: controller.TranslationService.Translation(context, "success_create_translation", payloadJwt.UserLangCode),
		Data:   translationResponse,
	}

	context.Writer.Header().Add("Content-Type", "application/json")
	context.JSON(200, webResponse)
}

func (controller *TranslationControllerImpl) Update(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	translationUpdateRequest := model.TranslationUpdateRequest{}
	context.Bind(&translationUpdateRequest)

	translationUpdateRequest.UpdatedBy = payloadJwt.UserId

	currentTime := time.Now()
	translationUpdateRequest.UpdatedAt = currentTime.Format("2006-01-02 15:04:05")

	translationId := context.Param("translationId")
	id, err := strconv.Atoi(translationId)
	helper.PanicIfError(err)

	translationUpdateRequest.TranslationId = id

	err = controller.Validate.Struct(translationUpdateRequest)
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

	translationResponse := controller.TranslationService.Update(context, translationUpdateRequest)

	if translationResponse.TranslationId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_update_translation", payloadJwt.UserLangCode),
			Data:   translationResponse,
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

func (controller *TranslationControllerImpl) Delete(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	translationId := context.Param("translationId")
	id, err := strconv.Atoi(translationId)
	helper.PanicIfError(err)

	translationResponse := controller.TranslationService.Delete(context, id)

	if translationResponse.TranslationId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_delete_translation", payloadJwt.UserLangCode),
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

func (controller *TranslationControllerImpl) FindById(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	translationId := context.Param("translationId")
	id, err := strconv.Atoi(translationId)
	helper.PanicIfError(err)

	translationResponse := controller.TranslationService.FindById(context, id)

	if translationResponse.TranslationId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_get_translation", payloadJwt.UserLangCode),
			Data:   translationResponse,
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

func (controller *TranslationControllerImpl) FindAll(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	translationResponses := controller.TranslationService.FindAll(context)

	if len(translationResponses) > 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_get_translation", payloadJwt.UserLangCode),
			Data:   translationResponses,
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
