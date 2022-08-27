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
	maxIdleConns := viper.GetInt(`database.mysql.maxIdleConns`)
	maxOpenConns := viper.GetInt(`database.mysql.maxOpenConns`)
	connMaxLifetime := viper.GetDuration(`database.mysql.connMaxLifetime`)
	connMaxIdleTime := viper.GetDuration(`database.mysql.connMaxIdleTime`)

	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	helper.PanicIfError(err)

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Minute)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Minute)

	return db
}
