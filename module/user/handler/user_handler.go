package handler

import (
	"collapp/configs"
	"collapp/helper"
	translationService "collapp/module/translation/service"
	"collapp/module/user/model"
	"collapp/module/user/service"
	"collapp/transport/http/middleware"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserService        service.UserService
	Validate           *validator.Validate
	TranslationService translationService.TranslationService
	config             *configs.Config
}

func NewUserHandler(db *sql.DB, cfg *configs.Config, userSvc service.UserService, translationSvc translationService.TranslationService) UserHandler {
	validate := validator.New()
	return UserHandler{
		UserService:        userSvc,
		Validate:           validate,
		TranslationService: translationSvc,
		config:             cfg,
	}
}

func (h *UserHandler) Create(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	password := []byte(h.config.DefaultPassword)

	userCreateRequest := model.UserCreateRequest{}
	context.Bind(&userCreateRequest)

	userCreateRequest.CreatedBy = payloadJwt.UserId

	currentTime := time.Now()
	userCreateRequest.CreatedAt = currentTime.Format("2006-01-02 15:04:05")

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	helper.IfError(err)

	userCreateRequest.UserPassword = string(hashedPassword)

	err = h.Validate.Struct(userCreateRequest)
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

	if userCreateRequest.UserPhoto != nil {
		var pathFile = h.config.Files.Photo
		var requestFile = "user_photo"

		fileName, err := helper.UploadFile(context, requestFile, pathFile)
		userCreateRequest.UserPhotoName = fileName
		if err != nil {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: h.TranslationService.Translation(context, "internal_server_error", payloadJwt.UserLangCode) + " " + h.TranslationService.Translation(context, "file_upload_failed", payloadJwt.UserLangCode),
				Data:   nil,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}
	}

	userResponse := h.UserService.Create(context, userCreateRequest)
	webResponse := helper.WebResponse{
		Code:   200,
		Status: h.TranslationService.Translation(context, "success_create_user", payloadJwt.UserLangCode),
		Data:   userResponse,
	}

	context.Writer.Header().Add("Content-Type", "application/json")
	context.JSON(200, webResponse)
}

func (h *UserHandler) Update(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	userUpdateRequest := model.UserUpdateRequest{}
	context.Bind(&userUpdateRequest)

	userUpdateRequest.UpdatedBy = payloadJwt.UserId

	currentTime := time.Now()
	userUpdateRequest.UpdatedAt = currentTime.Format("2006-01-02 15:04:05")

	userId := context.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.IfError(err)

	userUpdateRequest.UserId = id

	err = h.Validate.Struct(userUpdateRequest)
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

	var pathFile = h.config.Files.Photo

	if userUpdateRequest.UserPhoto != nil {
		var requestFile = "user_photo"

		fileName, err := helper.UploadFile(context, requestFile, pathFile)
		userUpdateRequest.UserPhotoName = fileName
		if err != nil {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: h.TranslationService.Translation(context, "internal_server_error", payloadJwt.UserLangCode) + " " + h.TranslationService.Translation(context, "file_upload_failed", payloadJwt.UserLangCode),
				Data:   nil,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}
	}

	userResponse, oldPhoto := h.UserService.Update(context, userUpdateRequest)

	if userUpdateRequest.UserPhoto != nil && oldPhoto != "" {
		helper.DeleteFile(oldPhoto, pathFile)
	}

	if userResponse.UserId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_update_user", payloadJwt.UserLangCode),
			Data:   userResponse,
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

func (h *UserHandler) Delete(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	userDeleteRequest := model.UserDeleteRequest{}
	context.Bind(&userDeleteRequest)

	userId := context.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.IfError(err)

	if userDeleteRequest.IsSoftDelete {
		userDeleteRequest.UserId = id

		userDeleteRequest.DeletedBy = payloadJwt.UserId

		currentTime := time.Now()
		userDeleteRequest.DeletedAt = currentTime.Format("2006-01-02 15:04:05")

		userResponse := h.UserService.SoftDelete(context, userDeleteRequest)

		if userResponse.UserId != 0 {
			webResponse := helper.WebResponse{
				Code:   200,
				Status: h.TranslationService.Translation(context, "success_delete_user", payloadJwt.UserLangCode),
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
	} else {
		userResponse := h.UserService.Delete(context, id)

		if userResponse.UserId != 0 {
			if userResponse.UserPhoto != "" {
				var pathFile = h.config.Files.Photo
				helper.DeleteFile(userResponse.UserPhoto, pathFile)
			}
			webResponse := helper.WebResponse{
				Code:   200,
				Status: h.TranslationService.Translation(context, "success_delete_user", payloadJwt.UserLangCode),
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
}

func (h *UserHandler) FindById(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)
	fmt.Println(payloadJwt)

	userId := context.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.IfError(err)

	userResponse := h.UserService.FindById(context, id)

	if userResponse.UserId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_get_user", payloadJwt.UserLangCode),
			Data:   userResponse,
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

func (h *UserHandler) FindAll(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	userResponses := h.UserService.FindAll(context)

	if len(userResponses) > 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_get_user", payloadJwt.UserLangCode),
			Data:   userResponses,
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

func (h *UserHandler) Login(context *gin.Context) {
	jwtKey := []byte(h.config.JWT.Key)
	defaultLang := h.config.DefaultLang

	currentTime := time.Now()
	userLoginRequest := model.UserLoginRequest{}
	context.Bind(&userLoginRequest)

	err := h.Validate.Struct(userLoginRequest)
	if err != nil {
		webResponse := helper.WebResponse{
			Code:   http.StatusBadRequest,
			Status: h.TranslationService.Translation(context, "bad_request", defaultLang),
			Data:   err.Error(),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	userCheck := h.UserService.FindByEmail(context, userLoginRequest.UserEmail)

	err = bcrypt.CompareHashAndPassword([]byte(userCheck.UserPassword), []byte(userLoginRequest.UserPassword))

	if err == nil {

		userResponse := h.UserService.FindById(context, userCheck.UserId)

		expired := h.config.JWT.Expired
		expiredRefresh := h.config.JWT.ExpiredRefresh

		// start cretae JWT
		expirationTime := time.Now().Add(expired)
		claims := middleware.Claims{
			UserId:       userResponse.UserId,
			UserName:     userResponse.UserName,
			UserEmail:    userResponse.UserEmail,
			UserLangCode: userResponse.UserLangCode,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: h.TranslationService.Translation(context, "internal_server_error", defaultLang),
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}

		expirationTimeRefresh := time.Now().Add(expiredRefresh)
		claimsRefresh := middleware.Claims{
			UserId:       userResponse.UserId,
			UserName:     userResponse.UserName,
			UserEmail:    userResponse.UserEmail,
			UserLangCode: userResponse.UserLangCode,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTimeRefresh.Unix(),
			},
		}
		tokenRefresgh := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
		tokenStringRefresh, err := tokenRefresgh.SignedString(jwtKey)

		if err != nil {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: h.TranslationService.Translation(context, "internal_server_error", defaultLang),
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}

		userResponse.UserToken = tokenString
		userResponse.UserTokenRefresh = tokenStringRefresh

		userData := model.UserUpdateTokenRequest{}

		userData.UserId = userResponse.UserId
		userData.UserToken = tokenString
		userData.UserTokenRefresh = tokenStringRefresh
		userData.UserLastLogin = currentTime.Format("2006-01-02 15:04:05")
		userTokenUpdateResponse := h.UserService.UpdateToken(context, userData)
		//end create JWT

		if userTokenUpdateResponse.UserEmail != "" {
			webResponse := helper.WebResponse{
				Code:   200,
				Status: h.TranslationService.Translation(context, "success_login", defaultLang),
				Data:   userResponse,
			}
			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(200, webResponse)
		} else {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: h.TranslationService.Translation(context, "internal_server_error", defaultLang),
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
		}
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: h.TranslationService.Translation(context, "worng_email_or_password", defaultLang),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusUnauthorized, webResponse)
	}
}

func (h *UserHandler) RefreshToken(context *gin.Context) {
	jwtKey := []byte(h.config.JWT.Key)

	currentTime := time.Now()
	userRefreshToken := context.Param("userRefreshToken")

	claims := &middleware.Claims{}

	tkn, _ := jwt.ParseWithClaims(userRefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if !tkn.Valid {
		webResponse := helper.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: h.TranslationService.Translation(context, "unauthorized", claims.UserLangCode),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusUnauthorized, webResponse)
		return
	}

	userCheck := h.UserService.FindByTokenRefresh(context, userRefreshToken)

	if userCheck.UserId != 0 {

		userResponse := h.UserService.FindById(context, userCheck.UserId)

		expired := h.config.JWT.Expired
		expiredRefresh := h.config.JWT.ExpiredRefresh

		// start cretae JWT
		expirationTime := time.Now().Add(expired)
		claims := middleware.Claims{
			UserId:       userResponse.UserId,
			UserName:     userResponse.UserName,
			UserEmail:    userResponse.UserEmail,
			UserLangCode: userResponse.UserLangCode,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: h.TranslationService.Translation(context, "internal_server_error", claims.UserLangCode),
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}

		expirationTimeRefresh := time.Now().Add(expiredRefresh)
		claimsRefresh := middleware.Claims{
			UserId:       userResponse.UserId,
			UserName:     userResponse.UserName,
			UserEmail:    userResponse.UserEmail,
			UserLangCode: userResponse.UserLangCode,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTimeRefresh.Unix(),
			},
		}
		tokenRefresgh := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
		tokenStringRefresh, err := tokenRefresgh.SignedString(jwtKey)

		if err != nil {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: h.TranslationService.Translation(context, "internal_server_error", claims.UserLangCode),
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}

		userResponse.UserToken = tokenString
		userResponse.UserTokenRefresh = tokenStringRefresh

		userData := model.UserUpdateTokenRequest{}

		userData.UserId = userResponse.UserId
		userData.UserToken = tokenString
		userData.UserTokenRefresh = tokenStringRefresh
		userData.UserLastLogin = currentTime.Format("2006-01-02 15:04:05")
		userTokenUpdateResponse := h.UserService.UpdateToken(context, userData)
		//end create JWT

		if userTokenUpdateResponse.UserEmail != "" {
			webResponse := helper.WebResponse{
				Code:   200,
				Status: h.TranslationService.Translation(context, "refresh_token_success", claims.UserLangCode),
				Data:   userResponse,
			}
			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(200, webResponse)
		} else {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: h.TranslationService.Translation(context, "internal_server_error", claims.UserLangCode),
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
		}
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: h.TranslationService.Translation(context, "unauthorized", claims.UserLangCode),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusUnauthorized, webResponse)
	}
}

func (h *UserHandler) Logout(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	userId := 0
	value, ok := context.Get("user_id")
	if ok {
		userId = value.(int)
	}

	userResponse := h.UserService.Logout(context, userId)

	if userResponse.UserId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: h.TranslationService.Translation(context, "success_logout", payloadJwt.UserLangCode),
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
