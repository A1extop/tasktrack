package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConf struct {
	Host     string
	Port     string
	Mode     string
	LogLevel string
}

type Config struct {
	App AppConf
}

func init() {
	println("Loading .env file...")
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: Could not load .env file")
		return
	}
	println("Successfully Loaded .env file")
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func New() *Config {
	return &Config{
		App: AppConf{
			Host:     getEnv("HTTP_HOST", "0.0.0.0"),
			Port:     getEnv("HTTP_PORT", "8080"),
			Mode:     getEnv("APP_MODE", "debug"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
	}
}
