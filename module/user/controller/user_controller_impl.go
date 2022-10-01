package controller

import (
	"collapp/helper"
	"collapp/middleware"
	translationService "collapp/module/translation/service"
	"collapp/module/user/model"
	"collapp/module/user/service"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type UserControllerImpl struct {
	UserService        service.UserService
	Validate           *validator.Validate
	TranslationService translationService.TranslationService
}

func NewUserController(db *sql.DB) UserController {
	validate := validator.New()
	userService := service.NewUserService(db)
	translationService := translationService.NewTranslationService(db)
	return &UserControllerImpl{
		UserService:        userService,
		Validate:           validate,
		TranslationService: translationService,
	}
}

func (controller *UserControllerImpl) Create(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	password := []byte(viper.GetString("defaultPassword"))

	userCreateRequest := model.UserCreateRequest{}
	context.Bind(&userCreateRequest)

	userCreateRequest.CreatedBy = payloadJwt.UserId

	currentTime := time.Now()
	userCreateRequest.CreatedAt = currentTime.Format("2006-01-02 15:04:05")

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	helper.PanicIfError(err)

	userCreateRequest.UserPassword = string(hashedPassword)

	err = controller.Validate.Struct(userCreateRequest)
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

	if userCreateRequest.UserPhoto != nil {
		var pathFile = viper.GetString("files.photo")
		var requestFile = "user_photo"

		fileName, err := helper.UploadFile(context, requestFile, pathFile)
		userCreateRequest.UserPhotoName = fileName
		if err != nil {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: controller.TranslationService.Translation(context, "internal_server_error", payloadJwt.UserLangCode) + " " + controller.TranslationService.Translation(context, "file_upload_failed", payloadJwt.UserLangCode),
				Data:   nil,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}
	}

	userResponse := controller.UserService.Create(context, userCreateRequest)
	webResponse := helper.WebResponse{
		Code:   200,
		Status: controller.TranslationService.Translation(context, "success_create_user", payloadJwt.UserLangCode),
		Data:   userResponse,
	}

	context.Writer.Header().Add("Content-Type", "application/json")
	context.JSON(200, webResponse)
}

func (controller *UserControllerImpl) Update(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	userUpdateRequest := model.UserUpdateRequest{}
	context.Bind(&userUpdateRequest)

	userUpdateRequest.UpdatedBy = payloadJwt.UserId

	currentTime := time.Now()
	userUpdateRequest.UpdatedAt = currentTime.Format("2006-01-02 15:04:05")

	userId := context.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	userUpdateRequest.UserId = id

	err = controller.Validate.Struct(userUpdateRequest)
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

	var pathFile = viper.GetString("files.photo")

	if userUpdateRequest.UserPhoto != nil {
		var requestFile = "user_photo"

		fileName, err := helper.UploadFile(context, requestFile, pathFile)
		userUpdateRequest.UserPhotoName = fileName
		if err != nil {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: controller.TranslationService.Translation(context, "internal_server_error", payloadJwt.UserLangCode) + " " + controller.TranslationService.Translation(context, "file_upload_failed", payloadJwt.UserLangCode),
				Data:   nil,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}
	}

	userResponse, oldPhoto := controller.UserService.Update(context, userUpdateRequest)

	if userUpdateRequest.UserPhoto != nil && oldPhoto != "" {
		helper.DeleteFile(oldPhoto, pathFile)
	}

	if userResponse.UserId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_update_user", payloadJwt.UserLangCode),
			Data:   userResponse,
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

func (controller *UserControllerImpl) Delete(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	userDeleteRequest := model.UserDeleteRequest{}
	context.Bind(&userDeleteRequest)

	userId := context.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	if userDeleteRequest.IsSoftDelete {
		userDeleteRequest.UserId = id

		userDeleteRequest.DeletedBy = payloadJwt.UserId

		currentTime := time.Now()
		userDeleteRequest.DeletedAt = currentTime.Format("2006-01-02 15:04:05")

		userResponse := controller.UserService.SoftDelete(context, userDeleteRequest)

		if userResponse.UserId != 0 {
			webResponse := helper.WebResponse{
				Code:   200,
				Status: controller.TranslationService.Translation(context, "success_delete_user", payloadJwt.UserLangCode),
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
	} else {
		userResponse := controller.UserService.Delete(context, id)

		if userResponse.UserId != 0 {
			if userResponse.UserPhoto != "" {
				var pathFile = viper.GetString("files.photo")
				helper.DeleteFile(userResponse.UserPhoto, pathFile)
			}
			webResponse := helper.WebResponse{
				Code:   200,
				Status: controller.TranslationService.Translation(context, "success_delete_user", payloadJwt.UserLangCode),
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
}

func (controller *UserControllerImpl) FindById(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	userId := context.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	userResponse := controller.UserService.FindById(context, id)

	if userResponse.UserId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_get_user", payloadJwt.UserLangCode),
			Data:   userResponse,
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

func (controller *UserControllerImpl) FindAll(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	userResponses := controller.UserService.FindAll(context)

	if len(userResponses) > 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_get_user", payloadJwt.UserLangCode),
			Data:   userResponses,
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

func (controller *UserControllerImpl) Login(context *gin.Context) {
	jwtKey := []byte(viper.GetString(`jwt.key`))
	defaultLang := viper.GetString(`defaultLang`)

	currentTime := time.Now()
	userLoginRequest := model.UserLoginRequest{}
	context.Bind(&userLoginRequest)

	err := controller.Validate.Struct(userLoginRequest)
	if err != nil {
		webResponse := helper.WebResponse{
			Code:   http.StatusBadRequest,
			Status: controller.TranslationService.Translation(context, "bad_request", defaultLang),
			Data:   err.Error(),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusBadRequest, webResponse)
		return
	}

	userCheck := controller.UserService.FindByEmail(context, userLoginRequest.UserEmail)

	err = bcrypt.CompareHashAndPassword([]byte(userCheck.UserPassword), []byte(userLoginRequest.UserPassword))

	if err == nil {

		userResponse := controller.UserService.FindById(context, userCheck.UserId)

		expired := viper.GetInt(`jwt.expired`)
		expiredRefresh := viper.GetInt(`jwt.expiredRefresh`)

		// start cretae JWT
		expirationTime := time.Now().Add(time.Duration(expired) * time.Minute)
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
				Status: controller.TranslationService.Translation(context, "internal_server_error", defaultLang),
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}

		expirationTimeRefresh := time.Now().Add(time.Duration(expiredRefresh) * time.Minute)
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
				Status: controller.TranslationService.Translation(context, "internal_server_error", defaultLang),
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
		userTokenUpdateResponse := controller.UserService.UpdateToken(context, userData)
		//end create JWT

		if userTokenUpdateResponse.UserEmail != "" {
			webResponse := helper.WebResponse{
				Code:   200,
				Status: controller.TranslationService.Translation(context, "success_login", defaultLang),
				Data:   userResponse,
			}
			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(200, webResponse)
		} else {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: controller.TranslationService.Translation(context, "internal_server_error", defaultLang),
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
		}
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: controller.TranslationService.Translation(context, "worng_email_or_password", defaultLang),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusUnauthorized, webResponse)
	}
}

func (controller *UserControllerImpl) RefreshToken(context *gin.Context) {
	jwtKey := []byte(viper.GetString(`jwt.key`))

	currentTime := time.Now()
	userRefreshToken := context.Param("userRefreshToken")

	claims := &middleware.Claims{}

	tkn, _ := jwt.ParseWithClaims(userRefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if !tkn.Valid {
		webResponse := helper.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: controller.TranslationService.Translation(context, "unauthorized", claims.UserLangCode),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusUnauthorized, webResponse)
		return
	}

	userCheck := controller.UserService.FindByTokenRefresh(context, userRefreshToken)

	if userCheck.UserId != 0 {

		userResponse := controller.UserService.FindById(context, userCheck.UserId)

		expired := viper.GetInt(`jwt.expired`)
		expiredRefresh := viper.GetInt(`jwt.expiredRefresh`)

		// start cretae JWT
		expirationTime := time.Now().Add(time.Duration(expired) * time.Minute)
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
				Status: controller.TranslationService.Translation(context, "internal_server_error", claims.UserLangCode),
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}

		expirationTimeRefresh := time.Now().Add(time.Duration(expiredRefresh) * time.Minute)
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
				Status: controller.TranslationService.Translation(context, "internal_server_error", claims.UserLangCode),
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
		userTokenUpdateResponse := controller.UserService.UpdateToken(context, userData)
		//end create JWT

		if userTokenUpdateResponse.UserEmail != "" {
			webResponse := helper.WebResponse{
				Code:   200,
				Status: controller.TranslationService.Translation(context, "refresh_token_success", claims.UserLangCode),
				Data:   userResponse,
			}
			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(200, webResponse)
		} else {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: controller.TranslationService.Translation(context, "internal_server_error", claims.UserLangCode),
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
		}
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: controller.TranslationService.Translation(context, "unauthorized", claims.UserLangCode),
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusUnauthorized, webResponse)
	}
}

func (controller *UserControllerImpl) Logout(context *gin.Context) {
	payloadJwt := helper.PayloadJwt(context)

	userId := 0
	value, ok := context.Get("user_id")
	if ok {
		userId = value.(int)
	}

	userResponse := controller.UserService.Logout(context, userId)

	if userResponse.UserId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: controller.TranslationService.Translation(context, "success_logout", payloadJwt.UserLangCode),
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
