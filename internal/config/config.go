package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port    string
	GinMode string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBParams   string
}

func Load() *Config {
	_ = godotenv.Load() // loads .env if present (noop in prod containers)

	cfg := &Config{
		Port:       get("PORT", "8080"),
		GinMode:    get("GIN_MODE", "debug"),
		DBHost:     get("DB_HOST", "db"),
		DBPort:     get("DB_PORT", "3306"),
		DBUser:     get("DB_USER", "root"),
		DBPassword: get("DB_PASSWORD", ""),
		DBName:     get("DB_NAME", "appdb"),
		DBParams:   get("DB_PARAMS", "charset=utf8mb4&parseTime=True&loc=Local"),
	}
	log.Printf("config loaded: port=%s gin_mode=%s db=%s@%s:%s/%s",
		cfg.Port, cfg.GinMode, cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
	return cfg
}

func get(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
