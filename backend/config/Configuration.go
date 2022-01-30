package config

import (
	"cultivator.wurmatron.io/backend/model"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

// Config values

var Host string
var Port int64
var StatusTimeoutInterval int64

func LoadOrSetupConfiguration(db *gorm.DB) {
	Host = findOrSetConfigValue(map[string]interface{}{"type": "backend", "key": "host"}, db, "host", "localhost").Value // TODO Lookup current ip, set as default
	Port = toInt(findOrSetConfigValue(map[string]interface{}{"type": "backend", "key": "port"}, db, "port", "8123"))
	StatusTimeoutInterval = toInt(findOrSetConfigValue(map[string]interface{}{"type": "backend", "key": "status_timeout_interval"}, db, "status_timeout_interval", "300"))
}

func toInt(cfg *model.Configuration) int64 {
	val, err := strconv.ParseInt(cfg.Value, 10, 64)
	if err != nil {
		log.Println("Failed to parse config entry for '" + cfg.Key + "' (" + cfg.Value + ")")
	}
	return val
}

func findOrSetConfigValue(conditions map[string]interface{}, db *gorm.DB, key string, value string) *model.Configuration {
	// Check if key exists
	count := int64(0)
	db.Table("configuration").Limit(1).Where(conditions).Count(&count)
	// Exists
	if count > 0 {
		var config model.Configuration
		db.Table("configuration").Where(conditions).FirstOrInit(&config)
		return &config
	} else {
		defaultConfig := &model.Configuration{
			Blockchain: "N/A",
			Type:       "backend",
			Key:        key,
			Value:      value,
			LastUpdate: time.Now().Unix(),
		}
		log.Println("Creating config entry for '" + defaultConfig.Key + "' (" + defaultConfig.Value + ") on " + defaultConfig.Type)
		db.Table("configuration").Create(defaultConfig)
		return defaultConfig
	}
}
