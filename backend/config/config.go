package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port           string
	DBConn         string
	JWTSecret      string
	JWTExpiry      int
	AllowedOrigins string
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
		log.Fatal("DB_CONN is required")
	}

	jwtSecret := viper.GetString("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required (set via environment variable or .env)")
	}

	jwtExpiry := viper.GetInt("JWT_EXPIRY_HOURS")
	if jwtExpiry == 0 {
		jwtExpiry = 24
	}

	allowedOrigins := viper.GetString("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	return &Config{
		Port:           port,
		DBConn:         dbConn,
		JWTSecret:      jwtSecret,
		JWTExpiry:      jwtExpiry,
		AllowedOrigins: allowedOrigins,
	}
}
