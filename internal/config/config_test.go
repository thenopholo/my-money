package config

import (
	"testing"
)

func TestLoad_Sucesso(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/mydb")
	t.Setenv("JWT_SECRET", "super-secret")
	t.Setenv("PORT", "8080")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.DatabaseURL != "postgres://localhost/mydb" {
		t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, "postgres://localhost/mydb")
	}
	if cfg.JWTSecret != "super-secret" {
		t.Errorf("JWTSecret = %q, want %q", cfg.JWTSecret, "super-secret")
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
}

func TestLoad_PortaPadrao(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/mydb")
	t.Setenv("JWT_SECRET", "super-secret")
	t.Setenv("PORT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "4235" {
		t.Errorf("Port = %q, want padrão %q", cfg.Port, "4235")
	}
}

func TestLoad_SemDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "super-secret")

	_, err := Load()
	if err == nil {
		t.Error("Load() deveria retornar erro para DATABASE_URL vazio")
	}
}

func TestLoad_SemJWTSecret(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/mydb")
	t.Setenv("JWT_SECRET", "")

	_, err := Load()
	if err == nil {
		t.Error("Load() deveria retornar erro para JWT_SECRET vazio")
	}
}

func TestLoad_JWTDurationPadrao(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/mydb")
	t.Setenv("JWT_SECRET", "super-secret")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.JWTDuration.Hours() != 24 {
		t.Errorf("JWTDuration = %v, want 24h", cfg.JWTDuration)
	}
}
