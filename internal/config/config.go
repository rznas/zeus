package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// AppConfig holds application-level settings
type AppConfig struct {
	Port              string
	Env               string
	JWTSecret         string
	JWTExpiresMinutes int
	RateLimitPerMin   int
	OTPRatePerMin     int
	OTPTTLSeconds     int
}

// PostgresConfig holds Postgres settings
type PostgresConfig struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
}

// RedisConfig holds Redis settings
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// Config is the root configuration object
type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

// Load loads configuration from environment variables and optional .env file
func Load() *Config {
	// Try .env first, then sample.env as fallback
	if err := godotenv.Load(".env"); err != nil {
		if err2 := godotenv.Load("sample.env"); err2 == nil {
			log.Printf("loaded environment from sample.env")
		} else {
			log.Printf("no .env or sample.env found, using system environment")
		}
	} else {
		log.Printf("loaded environment from .env")
	}

	cfg := &Config{
		App: AppConfig{
			Port:              getenv("APP_PORT", "8080"),
			Env:               getenv("APP_ENV", "development"),
			JWTSecret:         getenv("JWT_SECRET", "supersecretjwt"),
			JWTExpiresMinutes: getenvInt("JWT_EXPIRES_MINUTES", 60),
			RateLimitPerMin:   getenvInt("RATE_LIMIT_PER_MINUTE", 60),
			OTPRatePerMin:     getenvInt("OTP_RATE_LIMIT_PER_MINUTE", 3),
			OTPTTLSeconds:     getenvInt("OTP_TTL_SECONDS", 300),
		},
		Postgres: PostgresConfig{
			Host:     getenv("POSTGRES_HOST", "localhost"),
			Port:     getenv("POSTGRES_PORT", "5432"),
			DB:       getenv("POSTGRES_DB", "zeus"),
			User:     getenv("POSTGRES_USER", "zeus"),
			Password: getenv("POSTGRES_PASSWORD", "zeus"),
		},
		Redis: RedisConfig{
			Addr:     getenv("REDIS_ADDR", "localhost:6379"),
			Password: getenv("REDIS_PASSWORD", ""),
			DB:       getenvInt("REDIS_DB", 0),
		},
	}

	log.Printf("config loaded: env=%s port=%s psql=%s:%s/%s redis=%s", cfg.App.Env, cfg.App.Port, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DB, cfg.Redis.Addr)
	return cfg
}
