package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig
	DB  DBConfig
	JWT JWTConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Timezone string
}

type JWTConfig struct {
	Secret     string
	ExpiryHour int
	RefreshExpiryHour int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] .env file not found, reading from environment variables")
	}

	expiryHour, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOUR", "24"))
	refreshExpiryHour, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_HOUR", "168" ))

	return &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "MyBank"),
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "mybank_db"),
			Timezone: getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			ExpiryHour: expiryHour,
			RefreshExpiryHour: refreshExpiryHour,
		},
	}
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		c.User,
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Name,
		url.QueryEscape(c.Timezone),
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}