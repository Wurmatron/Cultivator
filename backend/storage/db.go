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
	host     = "postgres"
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

func CreateTablesIfNeeded() {
	DB.Raw("CREATE TABLE IF NOT EXISTS public.configuration(blockchain text COLLATE pg_catalog.\"default\" NOT NULL, type text COLLATE pg_catalog.\"default\" NOT NULL, key text COLLATE pg_catalog.\"default\" NOT NULL, value text COLLATE pg_catalog.\"default\" NOT NULL, last_update bigint DEFAULT 0);")
	DB.Raw("CREATE TABLE IF NOT EXISTS public.metrics(entry_type text COLLATE pg_catalog.\"default\" NOT NULL, \"timestamp\" bigint NOT NULL, type text COLLATE pg_catalog.\"default\", value text COLLATE pg_catalog.\"default\",blockchain text COLLATE pg_catalog.\"default\", style text COLLATE pg_catalog.\"default\");")
}
