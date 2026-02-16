package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	DBConn    string
	APIKey    string
	JWTSecret string
	JWTExpiry int
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

	apiKey := viper.GetString("API_KEY")

	jwtSecret := viper.GetString("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "kasir-secret-key-change-in-production"
		log.Println("Warning: JWT_SECRET not set, using default (not safe for production)")
	}

	jwtExpiry := viper.GetInt("JWT_EXPIRY_HOURS")
	if jwtExpiry == 0 {
		jwtExpiry = 24
	}

	return &Config{
		Port:      port,
		DBConn:    dbConn,
		APIKey:    apiKey,
		JWTSecret: jwtSecret,
		JWTExpiry: jwtExpiry,
	}
}
