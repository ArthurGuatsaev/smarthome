package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr        string
	LogLevel        string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	DBPath          string
	APIKey          string
}

func Load() Config {
	return Config{
		HTTPAddr:        getenv("HTTP_ADDR", ":8080"),
		LogLevel:        getenv("LOG_LEVEL", "info"),
		ShutdownTimeout: getenvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),

		ReadTimeout:  getenvDuration("HTTP_READ_TIMEOUT", 5*time.Second),
		WriteTimeout: getenvDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getenvDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
		DBPath:       getenv("DB_PATH", "./data/smarthome.db"),
		APIKey:       getenv("API_KEY", "devkey"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Поддержка форматов: "1500ms", "10s", "2m" (как time.ParseDuration),
// а также просто число секунд: "10"
func getenvDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	if d, err := time.ParseDuration(v); err == nil {
		return d
	}
	if secs, err := strconv.Atoi(v); err == nil {
		return time.Duration(secs) * time.Second
	}
	return def
}
