package middleware

import (
	"collapp/configs"
	"collapp/helper"
	translationService "collapp/module/translation/service"
	"collapp/module/user/service"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	UserId       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	UserEmail    string `json:"user_email"`
	UserLangCode string `json:"user_lang_code"`
	jwt.StandardClaims
}

type AuthMiddleware struct {
	config             *configs.Config
	translationService translationService.TranslationService
	userService        service.UserService
}

func NewAuthMiddleware(cfg *configs.Config, translationService translationService.TranslationService, userService service.UserService) AuthMiddleware {
	return AuthMiddleware{
		config:             cfg,
		translationService: translationService,
		userService:        userService,
	}
}

func (a *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(context *gin.Context) {

		jwtKey := []byte(a.config.JWT.Key)
		defaultLang := a.config.DefaultLang

		reqToken := context.Request.Header.Get("Authorization")
		if reqToken == "" {
			webResponse := helper.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: a.translationService.Translation(context, "unauthorized", defaultLang),
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusUnauthorized, webResponse)
			context.Abort()
			return
		}

		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]

		claims := &Claims{}

		tkn, _ := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if !tkn.Valid {
			webResponse := helper.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: a.translationService.Translation(context, "unauthorized", defaultLang),
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusUnauthorized, webResponse)
			context.Abort()
			return
		}

		userResponse := a.userService.FindById(context.Request.Context(), claims.UserId)

		if userResponse.UserToken == reqToken {
			context.Set("user_id", claims.UserId)
			context.Set("user_email", claims.UserEmail)
			context.Set("user_name", claims.UserName)
			context.Set("user_lang_code", claims.UserLangCode)
			context.Next()
		} else {
			webResponse := helper.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: a.translationService.Translation(context, "unauthorized", defaultLang),
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusUnauthorized, webResponse)
			context.Abort()
		}
	}
}
