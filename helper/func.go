package helper

import (
	"collapp/module/user/model"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func PayloadJwt(context *gin.Context) model.User {
	payloadJwt := model.User{}

	user_id, ok := context.Get("user_id")
	if ok {
		payloadJwt.UserId = user_id.(int)
	}

	user_email, ok := context.Get("user_email")
	if ok {
		payloadJwt.UserEmail = user_email.(string)
	}

	user_name, ok := context.Get("user_name")
	if ok {
		payloadJwt.UserEmail = user_name.(string)
	}

	user_lang_code, ok := context.Get("user_lang_code")
	if ok {
		payloadJwt.UserLangCode = user_lang_code.(string)
	}

	return payloadJwt
}

func UploadFile(context *gin.Context, requestName string, destination string) (string, error) {
	currentTime := time.Now()

	file, handler, err := context.Request.FormFile(requestName)
	if err != nil {
		PanicIfError(err)
	}
	defer file.Close()

	var fileName = currentTime.Format("20060102150405") + handler.Filename
	f, err := os.OpenFile(destination+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		PanicIfError(err)
	}
	defer f.Close()

	_, err = io.Copy(f, file)

	return fileName, err
}

func DeleteFile(fileName string, destination string) {
	dir, _ := os.Getwd()
	fileLocation := filepath.Join(dir, destination+fileName)
	err := os.Remove(fileLocation)
	PanicIfError(err)
}
