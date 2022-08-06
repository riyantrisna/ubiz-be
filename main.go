package main

import (
	"collapp/app"
	"collapp/middleware"
	"collapp/module/lang"
	"collapp/module/user"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(middleware.CORS())

	db := app.NewDB()
	api := router.Group("/api/v1")

	user.Router(db, api)
	lang.Router(db, api)

	appPort := viper.GetString("address")
	if appPort == "" {
		appPort = ":9090"
	}
	router.Run(appPort)
}
