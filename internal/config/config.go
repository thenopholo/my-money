package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	JWTDuration time.Duration
	Port        string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        os.Getenv("PORT"),
		JWTDuration: 24 * time.Hour,
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("missing required env var: DATABASE_URL")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("missing required env var: JWT_SECRET")
	}
	if cfg.Port == "" {
		cfg.Port = "4235"
	}

	return cfg, nil
}
