package handler

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

type TranslationHandlerImpl struct {
	TranslationService service.TranslationService
	Validate           *validator.Validate
	config             *configs.Config
}

func NewTranslationHandler(db *sql.DB, cfg *configs.Config) TranslationHandler {
	validate := validator.New()
	translationService := service.NewTranslationService(db)
	return &TranslationHandlerImpl{
		TranslationService: translationService,
		Validate:           validate,
		config:             cfg,
	}
}

func (h *TranslationHandlerImpl) Create(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	translationCreateRequest := model.TranslationCreateRequest{}
	context.Bind(&translationCreateRequest)

	translationCreateRequest.CreatedBy = payloadJwt.UserId

	currentTime := time.Now()
	translationCreateRequest.CreatedAt = currentTime.Format("2006-01-02 15:04:05")

	keyIsExist := h.TranslationService.CheckKeyTranslationExist(context, translationCreateRequest.TranslationKey)
	if keyIsExist {
		webResponse := helper.WebResponse{
			Code:   http.StatusBadRequest,
			Status: h.TranslationService.Translation(context, "bad_request", payloadJwt.UserLangCode),
			Data:   h.TranslationService.Translation(context, "key_translation_is_exist", payloadJwt.UserLangCode) + " (" + translationCreateRequest.TranslationKey + ")",
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	err := h.Validate.Struct(translationCreateRequest)
	if err != nil {
		webResponse := helper.WebResponse{
			Code:   http.StatusBadRequest,
			Status: h.TranslationService.Translation(context, "bad_request", payloadJwt.UserLangCode),
			Data:   err.Error(),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	translationResponse := h.TranslationService.Create(context, translationCreateRequest)
	webResponse := helper.WebResponse{
		Code:   200,
		Status: h.TranslationService.Translation(context, "success_create_translation", payloadJwt.UserLangCode),
		Data:   translationResponse,
	}

	context.Writer.Header().Add("Content-Type", "application/json")
	context.JSON(200, webResponse)
}

func (h *TranslationHandlerImpl) Update(context *gin.Context) {
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

	err = h.Validate.Struct(translationUpdateRequest)
	if err != nil {
		webResponse := helper.WebResponse{
			Code:   http.StatusBadRequest,
			Status: h.TranslationService.Translation(context, "bad_request", payloadJwt.UserLangCode),
			Data:   err.Error(),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	translationResponse := h.TranslationService.Update(context, translationUpdateRequest)

	if translationResponse.TranslationId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_update_translation", payloadJwt.UserLangCode),
			Data:   translationResponse,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: h.TranslationService.Translation(context, "data_not_found", payloadJwt.UserLangCode),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (h *TranslationHandlerImpl) Delete(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	translationId := context.Param("translationId")
	id, err := strconv.Atoi(translationId)
	helper.PanicIfError(err)

	translationResponse := h.TranslationService.Delete(context, id)

	if translationResponse.TranslationId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_delete_translation", payloadJwt.UserLangCode),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: h.TranslationService.Translation(context, "data_not_found", payloadJwt.UserLangCode),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (h *TranslationHandlerImpl) FindById(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	translationId := context.Param("translationId")
	id, err := strconv.Atoi(translationId)
	helper.PanicIfError(err)

	translationResponse := h.TranslationService.FindById(context, id)

	if translationResponse.TranslationId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_get_translation", payloadJwt.UserLangCode),
			Data:   translationResponse,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: h.TranslationService.Translation(context, "data_not_found", payloadJwt.UserLangCode),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (h *TranslationHandlerImpl) FindAll(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	translationResponses := h.TranslationService.FindAll(context)

	if len(translationResponses) > 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_get_translation", payloadJwt.UserLangCode),
			Data:   translationResponses,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: h.TranslationService.Translation(context, "data_not_found", payloadJwt.UserLangCode),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}
