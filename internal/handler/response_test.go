package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestJSON(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	data := map[string]string{"message": "ok"}
	JSON(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	if result["message"] != "ok" {
		t.Errorf("message = %q, want %q", result["message"], "ok")
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	Error(w, http.StatusBadRequest, "bad request")

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	if result["error"] != "bad request" {
		t.Errorf("error = %q, want %q", result["error"], "bad request")
	}
}

func TestDecodeJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{
			name:    "JSON válido",
			body:    `{"name":"teste"}`,
			wantErr: false,
		},
		{
			name:    "JSON inválido",
			body:    `{invalid}`,
			wantErr: true,
		},
		{
			name:    "campo desconhecido",
			body:    `{"name":"teste","unknown_field":"value"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			var dst struct {
				Name string `json:"name"`
			}
			err := decodeJSON(r, &dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"UUID válido", "550e8400-e29b-41d4-a716-446655440000", false},
		{"UUID inválido", "not-a-uuid", true},
		{"UUID vazio", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := parseUUID(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUUID(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestParseDecimal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"decimal válido", "123.45", false},
		{"decimal inteiro", "100", false},
		{"decimal inválido", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := parseDecimal(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDecimal(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    string
		wantErr  bool
		wantDate time.Time
	}{
		{
			name:     "RFC3339 válido",
			value:    "2025-06-15T10:30:00Z",
			wantErr:  false,
			wantDate: time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			name:     "formato yyyy-mm-dd válido",
			value:    "2025-06-15",
			wantErr:  false,
			wantDate: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "data vazia",
			value:   "",
			wantErr: true,
		},
		{
			name:    "data inválida",
			value:   "not-a-date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := parseDate(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDate(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !result.Equal(tt.wantDate) {
				t.Errorf("parseDate(%q) = %v, want %v", tt.value, result, tt.wantDate)
			}
		})
	}
}

func TestParseDateOrNow(t *testing.T) {
	t.Parallel()

	t.Run("vazio retorna agora", func(t *testing.T) {
		t.Parallel()
		before := time.Now().Add(-1 * time.Second)
		result, err := parseDateOrNow("")
		after := time.Now().Add(1 * time.Second)
		if err != nil {
			t.Fatalf("parseDateOrNow(\"\") error = %v", err)
		}
		if result.Before(before) || result.After(after) {
			t.Errorf("parseDateOrNow(\"\") = %v, esperava ~now", result)
		}
	})

	t.Run("data válida retorna data", func(t *testing.T) {
		t.Parallel()
		result, err := parseDateOrNow("2025-06-15")
		if err != nil {
			t.Fatalf("parseDateOrNow(\"2025-06-15\") error = %v", err)
		}
		expected := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("parseDateOrNow(\"2025-06-15\") = %v, want %v", result, expected)
		}
	})
}
