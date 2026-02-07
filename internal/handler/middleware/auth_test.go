package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/auth"
)

func TestAuth_Sucesso(t *testing.T) {
	t.Parallel()

	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour)
	userID := uuid.New()

	token, err := jwtManager.Generate(userID, "test@email.com")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	var capturedUserID uuid.UUID
	handler := Auth(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, ok := UserIDFromContext(r.Context())
		if !ok {
			t.Error("UserIDFromContext() retornou false")
		}
		capturedUserID = uid
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if capturedUserID != userID {
		t.Errorf("userID = %v, want %v", capturedUserID, userID)
	}
}

func TestAuth_SemHeader(t *testing.T) {
	t.Parallel()

	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour)
	handler := Auth(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler não deveria ser chamado")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	assertErrorJSON(t, w, "missing authorization header")
}

func TestAuth_HeaderInvalido(t *testing.T) {
	t.Parallel()

	jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour)

	tests := []struct {
		name   string
		header string
	}{
		{"sem Bearer", "token-sem-bearer"},
		{"Basic em vez de Bearer", "Basic abc123"},
		{"Bearer sem token", "Bearer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			handler := Auth(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("handler não deveria ser chamado")
			}))

			req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
			req.Header.Set("Authorization", tt.header)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestAuth_TokenExpirado(t *testing.T) {
	t.Parallel()

	jwtManager := auth.NewJWTManager("test-secret", -1*time.Hour)
	token, _ := jwtManager.Generate(uuid.New(), "test@email.com")

	verifier := auth.NewJWTManager("test-secret", 1*time.Hour)
	handler := Auth(verifier)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler não deveria ser chamado")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_JWTManagerNil(t *testing.T) {
	t.Parallel()

	handler := Auth(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler não deveria ser chamado")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestUserIDFromContext_SemValor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	_, ok := UserIDFromContext(ctx)
	if ok {
		t.Error("UserIDFromContext() deveria retornar false para contexto vazio")
	}
}

func assertErrorJSON(t *testing.T, w *httptest.ResponseRecorder, expectedMsg string) {
	t.Helper()

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("resposta não é JSON válido: %v", err)
	}
	if body["error"] != expectedMsg {
		t.Errorf("error = %q, want %q", body["error"], expectedMsg)
	}
}
