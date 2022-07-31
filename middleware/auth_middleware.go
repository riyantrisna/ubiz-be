package middleware

import (
	"collapp/helper"
	"collapp/module/user/service"
	"database/sql"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var jwtKey = []byte(viper.GetString(`jwt.key`))

type Claims struct {
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	jwt.StandardClaims
}

func Auth(db *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		reqToken := context.Request.Header.Get("Authorization")
		if reqToken == "" {
			webResponse := helper.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: "Unauthorized",
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
				Status: "Unauthorized",
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusUnauthorized, webResponse)
			context.Abort()
			return
		}

		userService := service.NewUserService(db)
		userResponse := service.UserService.FindById(userService, context.Request.Context(), claims.UserId)

		if userResponse.UserToken == reqToken {
			context.Next()
		} else {
			webResponse := helper.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: "Unauthorized",
			}

			context.Writer.Header().Add("Content-Type", "application/json")
			context.JSON(http.StatusUnauthorized, webResponse)
			context.Abort()
		}
	}
}