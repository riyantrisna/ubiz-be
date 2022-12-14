package handler

import (
	"collapp/configs"
	"collapp/helper"
	"collapp/module/lang/model"
	"collapp/module/lang/service"
	translationService "collapp/module/translation/service"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LangHandler struct {
	LangService        service.LangService
	Validate           *validator.Validate
	TranslationService translationService.TranslationService
	config             *configs.Config
}

func NewLangHandler(db *sql.DB, cfg *configs.Config, langService service.LangService, translationService translationService.TranslationService) LangHandler {
	validate := validator.New()
	return LangHandler{
		LangService:        langService,
		Validate:           validate,
		TranslationService: translationService,
		config:             cfg,
	}
}

func (h *LangHandler) Create(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	langCreateRequest := model.LangCreateRequest{}
	context.Bind(&langCreateRequest)

	langCreateRequest.CreatedBy = payloadJwt.UserId

	currentTime := time.Now()
	langCreateRequest.CreatedAt = currentTime.Format("2006-01-02 15:04:05")

	err := h.Validate.Struct(langCreateRequest)
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

	langResponse := h.LangService.Create(context, langCreateRequest)
	webResponse := helper.WebResponse{
		Code:   200,
		Status: h.TranslationService.Translation(context, "success_create_language", payloadJwt.UserLangCode),
		Data:   langResponse,
	}

	context.Writer.Header().Add("Content-Type", "application/json")
	context.JSON(200, webResponse)
}

func (h *LangHandler) Update(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	langUpdateRequest := model.LangUpdateRequest{}
	context.Bind(&langUpdateRequest)

	langUpdateRequest.UpdatedBy = payloadJwt.UserId

	currentTime := time.Now()
	langUpdateRequest.UpdatedAt = currentTime.Format("2006-01-02 15:04:05")

	langId := context.Param("langId")
	id, err := strconv.Atoi(langId)
	helper.IfError(err)

	langUpdateRequest.LangId = id

	err = h.Validate.Struct(langUpdateRequest)
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

	langResponse := h.LangService.Update(context, langUpdateRequest)

	if langResponse.LangId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_update_language", payloadJwt.UserLangCode),
			Data:   langResponse,
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

func (h *LangHandler) Delete(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	langId := context.Param("langId")
	id, err := strconv.Atoi(langId)
	helper.IfError(err)

	langResponse := h.LangService.Delete(context, id)

	if langResponse.LangId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_delete_language", payloadJwt.UserLangCode),
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

func (h *LangHandler) FindById(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	langId := context.Param("langId")
	id, err := strconv.Atoi(langId)
	helper.IfError(err)

	langResponse := h.LangService.FindById(context, id)

	if langResponse.LangId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_get_language", payloadJwt.UserLangCode),
			Data:   langResponse,
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

func (h *LangHandler) FindAll(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	langResponses := h.LangService.FindAll(context)

	if len(langResponses) > 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_get_language", payloadJwt.UserLangCode),
			Data:   langResponses,
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
