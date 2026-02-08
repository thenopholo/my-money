package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	JWTDuration        time.Duration
	Port               string
	CORSAllowedOrigins []string
	LLMAgentURL        string
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

	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins != "" {
		cfg.CORSAllowedOrigins = strings.Split(corsOrigins, ",")
	} else {
		cfg.CORSAllowedOrigins = []string{"http://localhost:3000"}
	}

	cfg.LLMAgentURL = os.Getenv("LLM_AGENT_URL")
	if cfg.LLMAgentURL == "" {
		cfg.LLMAgentURL = "http://localhost:8001"
	}

	return cfg, nil
}
