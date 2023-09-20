package config

import (
	"fmt"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	config = viper.Init("db")
)

func GetDsn(key string) string {
	host := config.GetString(fmt.Sprintf("%s.host", key))
	port := config.GetInt(fmt.Sprintf("%s.port", key))
	database := config.GetString(fmt.Sprintf("%s.database", key))
	username := config.GetString(fmt.Sprintf("%s.username", key))
	password := config.GetString(fmt.Sprintf("%s.password", key))
	// data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	return dsn
}

func GetDB() *gorm.DB {
	return DB
}
