package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string
	DBConn string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file: %v", err)
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	dbConn := viper.GetString("DB_CONN")
	if dbConn == "" {
		log.Fatal("DB_CONN is required in .env file")
	}

	return &Config{
		Port:   port,
		DBConn: dbConn,
	}
}
