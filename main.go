package main

import (
	"collapp/configs"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize config
	_ = configs.Get()

	http := InitializeService()
	http.SetupAndServe()
}
