package app

import (
	"collapp/helper"
	"database/sql"
	"time"

	"github.com/spf13/viper"
)

func NewDB() *sql.DB {
	dbHost := viper.GetString(`database.mysql.dbHost`)
	dbPort := viper.GetString(`database.mysql.dbPort`)
	dbUser := viper.GetString(`database.mysql.dbUser`)
	dbPass := viper.GetString(`database.mysql.dbPass`)
	dbName := viper.GetString(`database.mysql.dbName`)

	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
