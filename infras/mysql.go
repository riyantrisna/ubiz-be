package infras

import (
	"collapp/configs"
	"collapp/helper"
	"database/sql"
	"fmt"
)

// NewMysqlDB will build a new mysql database connection
func NewMysqlDB(cfg *configs.Config) *sql.DB {
	dbHost := cfg.Database.Mysql.DbHost
	dbPort := cfg.Database.Mysql.DbPort
	dbUser := cfg.Database.Mysql.DbUser
	dbPass := cfg.Database.Mysql.DbPass
	dbName := cfg.Database.Mysql.DbName
	maxIdleConns := cfg.Database.Mysql.MaxIdleConns
	maxOpenConns := cfg.Database.Mysql.MaxOpenConns
	connMaxLifetime := cfg.Database.Mysql.ConnMaxLifeTime
	connMaxIdleTime := cfg.Database.Mysql.ConnMaxIdleTime

	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", url)
	helper.PanicIfError(err)

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	return db
}
