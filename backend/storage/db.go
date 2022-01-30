package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// TODO Temp Configuration
const (
	User     = "cultivator"
	password = "drowssap"
	Name     = "cultivator"
	host     = "0.0.0.0"
	port     = "5432"
	params   = "sslmode=disable"
)

func createDBConnectionStr() string {
	db := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s %s",
		host,
		port,
		User,
		Name,
		password,
		params)
	return db
}

var DB *gorm.DB

func NewDB() (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(postgres.Open(createDBConnectionStr()), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	return DB, err
}

func GetConnection() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}
	return NewDB()
}

func GetDBHost() string {
	return host + ":" + port
}
