package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AccessTokenSecret   string
	RefreshTokenSecret  string
	AccessTokenExpMin   int
	RefreshTokenExpDays int
}

var AppConfig *Config

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	accessExpMin, err := strconv.Atoi(getEnv("ACCESS_TOKEN_EXP_MIN", "15"))
	if err != nil {
		accessExpMin = 15
	}

	refreshExpDays, err := strconv.Atoi(getEnv("REFRESH_TOKEN_EXP_DAYS", "7"))
	if err != nil {
		refreshExpDays = 7
	}

	AppConfig = &Config{
		AccessTokenSecret:   getEnv("ACCESS_TOKEN_SECRET", "ACCESS_SECRET"),
		RefreshTokenSecret:  getEnv("REFRESH_TOKEN_SECRET", "REFRESH_SECRET"),
		AccessTokenExpMin:   accessExpMin,
		RefreshTokenExpDays: refreshExpDays,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}