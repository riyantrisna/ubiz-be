package main

import (
	"collapp/configs"
	"collapp/infras"
	"collapp/middleware"
	"collapp/module/lang"
	"collapp/module/translation"
	"collapp/module/user"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var config *configs.Config

func main() {
	// Initialize config
	config = configs.Get()

	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(middleware.CORS())

	db := infras.NewMysqlDB(config)
	api := router.Group("/api/v1")

	user.Router(db, api)
	lang.Router(db, api)
	translation.Router(db, api)

	appPort := viper.GetString("address")
	if appPort == "" {
		appPort = ":9090"
	}
	router.Run(appPort)
}
