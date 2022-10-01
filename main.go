package main

import (
	"collapp/configs"
	"collapp/infras"
	"collapp/module/lang"
	"collapp/module/translation"
	"collapp/module/user"
	httpTransport "collapp/transport/http"

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

	router := httpTransport.Setup()

	db := infras.NewMysqlDB(config)
	api := router.Group("/api/v1")

	user.Router(db, api)
	lang.Router(db, api)
	translation.Router(db, api)

	router.Run(config.Address)
}
