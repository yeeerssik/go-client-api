package models

import (
	"fmt"
	"go_client_service/config"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var db *gorm.DB
var Layout = "2006-01-02 15:44:05"
var DOBLayout = "2006-01-02"
var Location *time.Location

func Init() error {
	Location = time.Local
	return initializeDBConn()
}

// Db connection
func initializeDBConn() (err error) {
	connectionStr := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s&connection+timeout=30",
		config.Config.DataBase.Username,
		config.Config.DataBase.Password,
		config.Config.DataBase.Host,
		config.Config.DataBase.Port,
		config.Config.DataBase.Name,
	)
	fmt.Printf("%s", connectionStr)
	db, err = gorm.Open(sqlserver.Open(connectionStr), &gorm.Config{})
	return err
}
