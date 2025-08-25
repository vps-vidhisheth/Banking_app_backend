package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	APIKey    string
	JWTSecret string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("API_KEY")
	jwtSecret := os.Getenv("JWT_SECRET")

	return &Config{
		Port:      port,
		APIKey:    apiKey,
		JWTSecret: jwtSecret,
	}
}
