package main

import (
	"collapp/configs"

	_ "github.com/go-sql-driver/mysql"
)

var config *configs.Config

func main() {
	// Initialize config
	config = configs.Get()

	http := InitializeEvent()
	http.SetupAndServe()
}
