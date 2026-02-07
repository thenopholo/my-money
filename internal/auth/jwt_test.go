package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTManager_GenerateAndVerify(t *testing.T) {
	t.Parallel()

	manager := NewJWTManager("test-secret-key-super-safe", 1*time.Hour)
	userID := uuid.New()
	email := "joao@email.com"

	token, err := manager.Generate(userID, email)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if token == "" {
		t.Fatal("Generate() retornou token vazio")
	}

	claims, err := manager.Verify(token)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if claims.UserID != userID.String() {
		t.Errorf("UserID = %q, want %q", claims.UserID, userID.String())
	}
	if claims.Email != email {
		t.Errorf("Email = %q, want %q", claims.Email, email)
	}
}

func TestJWTManager_Verify_TokenExpirado(t *testing.T) {
	t.Parallel()

	manager := NewJWTManager("test-secret", -1*time.Hour) // duração negativa = já expirado
	userID := uuid.New()

	token, err := manager.Generate(userID, "test@email.com")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	_, err = manager.Verify(token)
	if err == nil {
		t.Error("Verify() deveria retornar erro para token expirado")
	}
}

func TestJWTManager_Verify_TokenInvalido(t *testing.T) {
	t.Parallel()

	manager := NewJWTManager("test-secret", 1*time.Hour)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "token completamente inválido",
			token: "not-a-jwt-token",
		},
		{
			name:  "token vazio",
			token: "",
		},
		{
			name:  "token manipulado",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.invalidsignature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := manager.Verify(tt.token)
			if err == nil {
				t.Error("Verify() deveria retornar erro para token inválido")
			}
		})
	}
}

func TestJWTManager_Verify_SecretDiferente(t *testing.T) {
	t.Parallel()

	manager1 := NewJWTManager("secret-1", 1*time.Hour)
	manager2 := NewJWTManager("secret-2", 1*time.Hour)

	token, err := manager1.Generate(uuid.New(), "test@email.com")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	_, err = manager2.Verify(token)
	if err == nil {
		t.Error("Verify() deveria retornar erro para secret diferente")
	}
}

func TestJWTManager_Generate_ClaimsCompletas(t *testing.T) {
	t.Parallel()

	manager := NewJWTManager("test-secret", 2*time.Hour)
	userID := uuid.New()

	token, err := manager.Generate(userID, "test@email.com")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	claims, err := manager.Verify(token)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	if claims.Subject != userID.String() {
		t.Errorf("Subject = %q, want %q", claims.Subject, userID.String())
	}
	if claims.IssuedAt == nil {
		t.Error("IssuedAt é nil")
	}
	if claims.ExpiresAt == nil {
		t.Error("ExpiresAt é nil")
	}
}
