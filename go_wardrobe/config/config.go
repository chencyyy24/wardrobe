package config

import (
	"os"
)

type Config struct {
	ServerPort  string
	DatabaseDSN string
	UploadDir   string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", "root:123456@tcp(localhost:3306)/sql_wardrobe?charset=utf8mb4&parseTime=True&loc=Local"),
		UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
		JWTSecret:   getEnv("JWT_SECRET", "wardrobe-secret-key-2024"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
