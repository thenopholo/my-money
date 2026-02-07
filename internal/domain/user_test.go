package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewUser_Validations(t *testing.T) {
	t.Parallel()

	validHash := strings.Repeat("a", 60) // simula bcrypt hash

	tests := []struct {
		name      string
		inputName string
		inputEmail string
		inputHash string
		wantErr   error
	}{
		{
			name:       "deve criar usuário com dados válidos",
			inputName:  "João",
			inputEmail: "joao@email.com",
			inputHash:  validHash,
			wantErr:    nil,
		},
		{
			name:       "deve retornar erro para nome vazio",
			inputName:  "",
			inputEmail: "joao@email.com",
			inputHash:  validHash,
			wantErr:    ErrEmptyName,
		},
		{
			name:       "deve retornar erro para email inválido sem domínio",
			inputName:  "João",
			inputEmail: "invalido",
			inputHash:  validHash,
			wantErr:    ErrInvalidEmail,
		},
		{
			name:       "deve retornar erro para email vazio",
			inputName:  "João",
			inputEmail: "",
			inputHash:  validHash,
			wantErr:    ErrInvalidEmail,
		},
		{
			name:       "deve retornar erro para email sem @",
			inputName:  "João",
			inputEmail: "joaoemail.com",
			inputHash:  validHash,
			wantErr:    ErrInvalidEmail,
		},
		{
			name:       "deve retornar erro para email sem TLD",
			inputName:  "João",
			inputEmail: "joao@email",
			inputHash:  validHash,
			wantErr:    ErrInvalidEmail,
		},
		{
			name:       "deve retornar erro para password hash curto",
			inputName:  "João",
			inputEmail: "joao@email.com",
			inputHash:  "short",
			wantErr:    ErrPasswordTooWeak,
		},
		{
			name:       "deve retornar erro para password hash com 49 caracteres",
			inputName:  "João",
			inputEmail: "joao@email.com",
			inputHash:  strings.Repeat("a", 49),
			wantErr:    ErrPasswordTooWeak,
		},
		{
			name:       "deve aceitar hash com exatamente 50 caracteres",
			inputName:  "João",
			inputEmail: "joao@email.com",
			inputHash:  strings.Repeat("a", 50),
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			user, err := NewUser(tt.inputName, tt.inputEmail, tt.inputHash)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if user == nil {
					t.Error("NewUser() retornou nil sem erro")
					return
				}
				if user.Name != tt.inputName {
					t.Errorf("NewUser() Name = %q, want %q", user.Name, tt.inputName)
				}
				if user.Email != tt.inputEmail {
					t.Errorf("NewUser() Email = %q, want %q", user.Email, tt.inputEmail)
				}
				if user.PasswordHash != tt.inputHash {
					t.Error("NewUser() PasswordHash não corresponde")
				}
				if user.ID.String() == "" {
					t.Error("NewUser() ID não foi gerado")
				}
				if user.CreatedAt.IsZero() {
					t.Error("NewUser() CreatedAt é zero")
				}
				if user.UpdatedAt.IsZero() {
					t.Error("NewUser() UpdatedAt é zero")
				}
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"email válido simples", "user@example.com", true},
		{"email com subdomain", "user@sub.example.com", true},
		{"email com ponto no local", "user.name@example.com", true},
		{"email com + no local", "user+tag@example.com", true},
		{"email com - no local", "user-name@example.com", true},
		{"email vazio", "", false},
		{"sem @", "userexample.com", false},
		{"sem domínio", "user@", false},
		{"sem TLD", "user@example", false},
		{"@ no início", "@example.com", false},
		{"espaço no email", "user @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isValidEmail(tt.email); got != tt.want {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}
