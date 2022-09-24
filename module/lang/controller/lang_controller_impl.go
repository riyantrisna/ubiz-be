package controller

import (
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
}

func NewLangController(db *sql.DB) LangController {
	validate := validator.New()
	langService := service.NewLangService(db)
	translationService := translationService.NewTranslationService(db)
	return &LangControllerImpl{
		LangService:        langService,
		Validate:           validate,
		TranslationService: translationService,
	}
}

func (controller *LangControllerImpl) Create(context *gin.Context) {
	langCreateRequest := model.LangCreateRequest{}
	context.Bind(&langCreateRequest)

	value, _ := context.Get("user_lang_code")
	user_lang_code := value.(string)

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
			Status: controller.TranslationService.Translation(context, "bad_request", user_lang_code),
			Data:   err.Error(),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	langResponse := controller.LangService.Create(context, langCreateRequest)
	webResponse := helper.WebResponse{
		Code:   200,
		Status: controller.TranslationService.Translation(context, "success_create_language", user_lang_code),
		Data:   langResponse,
	}

	context.Writer.Header().Add("Content-Type", "application/json")
	context.JSON(200, webResponse)
}

func (controller *LangControllerImpl) Update(context *gin.Context) {
	langUpdateRequest := model.LangUpdateRequest{}
	context.Bind(&langUpdateRequest)

	value, _ := context.Get("user_lang_code")
	user_lang_code := value.(string)

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
			Status: controller.TranslationService.Translation(context, "bad_request", user_lang_code),
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
			Status: controller.TranslationService.Translation(context, "success_update_language", user_lang_code),
			Data:   langResponse,
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

func (controller *LangControllerImpl) Delete(context *gin.Context) {
	langId := context.Param("langId")
	id, err := strconv.Atoi(langId)
	helper.PanicIfError(err)

	value, _ := context.Get("user_lang_code")
	user_lang_code := value.(string)

	langResponse := controller.LangService.Delete(context, id)

	if langResponse.LangId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_delete_language", user_lang_code),
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

func (controller *LangControllerImpl) FindById(context *gin.Context) {
	langId := context.Param("langId")
	id, err := strconv.Atoi(langId)
	helper.PanicIfError(err)

	value, _ := context.Get("user_lang_code")
	user_lang_code := value.(string)

	langResponse := controller.LangService.FindById(context, id)

	if langResponse.LangId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_get_language", user_lang_code),
			Data:   langResponse,
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

func (controller *LangControllerImpl) FindAll(context *gin.Context) {
	langResponses := controller.LangService.FindAll(context)

	value, _ := context.Get("user_lang_code")
	user_lang_code := value.(string)

	if len(langResponses) > 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_get_language", user_lang_code),
			Data:   langResponses,
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
