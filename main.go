package main

import (
	"collapp/configs"
	"collapp/infras"
	"collapp/module/lang"
	"collapp/module/translation"
	"collapp/module/user"
	httpTransport "collapp/transport/http"

	_ "github.com/go-sql-driver/mysql"
)

var config *configs.Config

func main() {
	// Initialize config
	config = configs.Get()

	router := httpTransport.Setup()

	db := infras.NewMysqlDB(config)
	api := router.Group("/api/v1")

	user.Router(db, api, config)
	lang.Router(db, api, config)
	translation.Router(db, api, config)

	router.Run(config.Address)
}
