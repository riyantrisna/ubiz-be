package controller

import (
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
}

func NewTranslationController(db *sql.DB) TranslationController {
	validate := validator.New()
	translationService := service.NewTranslationService(db)
	return &TranslationControllerImpl{
		TranslationService: translationService,
		Validate:           validate,
	}
}

func (controller *TranslationControllerImpl) Create(context *gin.Context) {
	translationCreateRequest := model.TranslationCreateRequest{}
	context.Bind(&translationCreateRequest)

	value, _ := context.Get("user_lang_code")
	user_lang_code := value.(string)

	value, ok := context.Get("user_id")
	if ok {
		translationCreateRequest.CreatedBy = value.(int)
	} else {
		translationCreateRequest.CreatedBy = 0
	}

	currentTime := time.Now()
	translationCreateRequest.CreatedAt = currentTime.Format("2006-01-02 15:04:05")

	keyIsExist := controller.TranslationService.CheckKeyTranslationExist(context, translationCreateRequest.TranslationKey)
	if keyIsExist {
		webResponse := helper.WebResponse{
			Code:   http.StatusBadRequest,
			Status: controller.TranslationService.Translation(context, "bad_request", user_lang_code),
			Data:   controller.TranslationService.Translation(context, "key_translation_is_exist", user_lang_code) + " (" + translationCreateRequest.TranslationKey + ")",
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	err := controller.Validate.Struct(translationCreateRequest)
	if err != nil {
		webResponse := helper.WebResponse{
			Code:   http.StatusBadRequest,
			Status: controller.TranslationService.Translation(context, "bad_request", user_lang_code),
			Data:   err.Error(),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	translationResponse := controller.TranslationService.Create(context, translationCreateRequest)
	webResponse := helper.WebResponse{
		Code:   200,
		Status: controller.TranslationService.Translation(context, "success_create_translation", user_lang_code),
		Data:   translationResponse,
	}

	context.Writer.Header().Add("Content-Type", "application/json")
	context.JSON(200, webResponse)
}

func (controller *TranslationControllerImpl) Update(context *gin.Context) {
	translationUpdateRequest := model.TranslationUpdateRequest{}
	context.Bind(&translationUpdateRequest)

	value, _ := context.Get("user_lang_code")
	user_lang_code := value.(string)

	value, ok := context.Get("user_id")
	if ok {
		translationUpdateRequest.UpdatedBy = value.(int)
	} else {
		translationUpdateRequest.UpdatedBy = 0
	}

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
			Status: controller.TranslationService.Translation(context, "bad_request", user_lang_code),
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
			Status: controller.TranslationService.Translation(context, "success_update_translation", user_lang_code),
			Data:   translationResponse,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: controller.TranslationService.Translation(context, "data_not_found", user_lang_code),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (controller *TranslationControllerImpl) Delete(context *gin.Context) {
	translationId := context.Param("translationId")
	id, err := strconv.Atoi(translationId)
	helper.PanicIfError(err)

	value, _ := context.Get("user_lang_code")
	user_lang_code := value.(string)

	translationResponse := controller.TranslationService.Delete(context, id)

	if translationResponse.TranslationId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_delete_translation", user_lang_code),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: controller.TranslationService.Translation(context, "data_not_found", user_lang_code),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (controller *TranslationControllerImpl) FindById(context *gin.Context) {
	translationId := context.Param("translationId")
	id, err := strconv.Atoi(translationId)
	helper.PanicIfError(err)

	value, _ := context.Get("user_lang_code")
	user_lang_code := value.(string)

	translationResponse := controller.TranslationService.FindById(context, id)

	if translationResponse.TranslationId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_get_translation", user_lang_code),
			Data:   translationResponse,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: controller.TranslationService.Translation(context, "data_not_found", user_lang_code),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}

func (controller *TranslationControllerImpl) FindAll(context *gin.Context) {
	translationResponses := controller.TranslationService.FindAll(context)

	value, _ := context.Get("user_lang_code")
	user_lang_code := value.(string)

	if len(translationResponses) > 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_get_translation", user_lang_code),
			Data:   translationResponses,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(200, webResponse)
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusNotFound,
			Status: controller.TranslationService.Translation(context, "data_not_found", user_lang_code),
			Data:   nil,
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusNotFound, webResponse)
	}
}
