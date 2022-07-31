package controller

import (
	"collapp/helper"
	"collapp/middleware"
	"collapp/module/user/model/web"
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

var password = []byte(viper.GetString(`defaultPassword`))
var jwtKey = []byte(viper.GetString(`jwt.key`))

type UserControllerImpl struct {
	UserService service.UserService
	Validate    *validator.Validate
}

func NewUserController(db *sql.DB) UserController {
	validate := validator.New()
	userService := service.NewUserService(db)
	return &UserControllerImpl{
		UserService: userService,
		Validate:    validate,
	}
}

func (controller *UserControllerImpl) Create(context *gin.Context) {
	userCreateRequest := web.UserCreateRequest{}
	context.Bind(&userCreateRequest)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	helper.PanicIfError(err)

	userCreateRequest.UserPassword = string(hashedPassword)

	err = controller.Validate.Struct(userCreateRequest)
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

	userResponse := controller.UserService.Create(context.Request.Context(), userCreateRequest)
	webResponse := helper.WebResponse{
		Code:   200,
		Status: "Success create user",
		Data:   userResponse,
	}

	context.Writer.Header().Add("Content-Type", "application/json")
	context.JSON(200, webResponse)
}

func (controller *UserControllerImpl) Update(context *gin.Context) {
	userUpdateRequest := web.UserUpdateRequest{}
	context.Bind(&userUpdateRequest)

	userId := context.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	userUpdateRequest.UserId = id

	err = controller.Validate.Struct(userUpdateRequest)
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

	userResponse := controller.UserService.Update(context.Request.Context(), userUpdateRequest)

	if userResponse.UserId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: "Success update user",
			Data:   userResponse,
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

func (controller *UserControllerImpl) Delete(context *gin.Context) {
	userId := context.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	userResponse := controller.UserService.Delete(context.Request.Context(), id)

	if userResponse.UserId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: "Success delete user",
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

func (controller *UserControllerImpl) FindById(context *gin.Context) {
	userId := context.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.PanicIfError(err)

	userResponse := controller.UserService.FindById(context.Request.Context(), id)

	if userResponse.UserId != 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: "Success get user",
			Data:   userResponse,
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

func (controller *UserControllerImpl) FindAll(context *gin.Context) {
	userResponses := controller.UserService.FindAll(context.Request.Context())

	if len(userResponses) > 0 {
		webResponse := helper.WebResponse{
			Code:   200,
			Status: "Success get all users",
			Data:   userResponses,
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

func (controller *UserControllerImpl) Login(context *gin.Context) {
	userLoginRequest := web.UserLoginRequest{}
	context.Bind(&userLoginRequest)

	err := controller.Validate.Struct(userLoginRequest)
	helper.PanicIfError(err)

	userCheck := controller.UserService.FindByEmail(context.Request.Context(), userLoginRequest.UserEmail)

	err = bcrypt.CompareHashAndPassword([]byte(userCheck.UserPassword), []byte(userLoginRequest.UserPassword))

	if err == nil {

		userResponse := controller.UserService.FindById(context.Request.Context(), userCheck.UserId)

		expired := viper.GetInt(`jwt.expired`)
		expiredRefresh := viper.GetInt(`jwt.expiredRefresh`)

		// start cretae JWT
		expirationTime := time.Now().Add(time.Duration(expired) * time.Minute)
		claims := middleware.Claims{
			UserId:    userResponse.UserId,
			UserName:  userResponse.UserName,
			UserEmail: userResponse.UserEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: "Internal Server Error",
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}

		expirationTimeRefresh := time.Now().Add(time.Duration(expiredRefresh) * time.Minute)
		claimsRefresh := middleware.Claims{
			UserId:    userResponse.UserId,
			UserName:  userResponse.UserName,
			UserEmail: userResponse.UserEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTimeRefresh.Unix(),
			},
		}
		tokenRefresgh := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
		tokenStringRefresh, err := tokenRefresgh.SignedString(jwtKey)

		if err != nil {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: "Internal Server Error",
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
			return
		}

		userResponse.UserToken = tokenString
		userResponse.UserTokenRefresh = tokenStringRefresh

		userTokenUpdateRequest := web.UserTokenUpdateRequest{}

		userTokenUpdateRequest.UserId = userResponse.UserId
		userTokenUpdateRequest.UserToken = tokenString
		userTokenUpdateRequest.UserTokenRefresh = tokenStringRefresh
		userTokenUpdateResponse := controller.UserService.UpdateToken(context.Request.Context(), userTokenUpdateRequest)
		//end create JWT

		if userTokenUpdateResponse.UserEmail != "" {
			webResponse := helper.WebResponse{
				Code:   200,
				Status: "Login success",
				Data:   userResponse,
			}
			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(200, webResponse)
		} else {
			webResponse := helper.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: "Internal Server Error",
				Data:   err,
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusInternalServerError, webResponse)
		}
	} else {
		webResponse := helper.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "Worng email or password",
		}

		context.Writer.Header().Add("Content-Type", "application/json")
		context.JSON(http.StatusUnauthorized, webResponse)
	}

}
