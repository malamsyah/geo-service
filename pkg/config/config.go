package config

import (
	"fmt"
	"log"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	TZ         string
	Host       string
}

// nolint: gochecknoglobals
var configInstance *Config

func init() {
	Load(path.Join(".", ".env"))
}

func Load(path string) *Config {
	fmt.Println("Loading config from", path)

	viper.AutomaticEnv()
	viper.SetConfigFile(path)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("can't load config from `.env`. environment variables will be used. err: %v", err)
	}

	configInstance = &Config{
		AppPort:    viper.GetString("APP_PORT"),
		DBHost:     viper.GetString("DB_HOST"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBPort:     viper.GetString("DB_PORT"),
		TZ:         viper.GetString("TZ"),
		Host:       viper.GetString("HOST"),
	}

	return configInstance
}

func Instance() *Config {
	return configInstance
}
