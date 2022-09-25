package helper

import (
	"collapp/module/user/model"

	"github.com/gin-gonic/gin"
)

func PayloadJwt(context *gin.Context) model.User {
	payloadJwt := model.User{}

	user_id, ok := context.Get("user_id")
	if ok {
		payloadJwt.UserId = user_id.(int)
	} else {
		payloadJwt.UserId = 0
	}

	user_email, _ := context.Get("user_email")
	payloadJwt.UserEmail = user_email.(string)

	user_name, _ := context.Get("user_name")
	payloadJwt.UserEmail = user_name.(string)

	user_lang_code, _ := context.Get("user_lang_code")
	payloadJwt.UserLangCode = user_lang_code.(string)

	return payloadJwt
}
