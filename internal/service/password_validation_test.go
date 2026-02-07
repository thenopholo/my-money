package service

import (
	"errors"
	"testing"

	"github.com/thenopholo/my-money/internal/domain"
)

func TestValidatePasswordStrength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "deve aceitar senha forte",
			password: "Str0ng!Pass",
			wantErr:  nil,
		},
		{
			name:     "deve aceitar senha com caracteres especiais variados",
			password: "Aa1@bcde",
			wantErr:  nil,
		},
		{
			name:     "deve retornar erro para senha curta (menos de 8)",
			password: "Aa1@bc",
			wantErr:  domain.ErrPasswordTooShort,
		},
		{
			name:     "deve retornar erro para senha com 7 caracteres",
			password: "Aa1@bcd",
			wantErr:  domain.ErrPasswordTooShort,
		},
		{
			name:     "deve retornar erro para senha sem maiúscula",
			password: "aaa1@bcde",
			wantErr:  domain.ErrPasswordTooWeak,
		},
		{
			name:     "deve retornar erro para senha sem minúscula",
			password: "AAA1@BCDE",
			wantErr:  domain.ErrPasswordTooWeak,
		},
		{
			name:     "deve retornar erro para senha sem número",
			password: "Aaaa@bcde",
			wantErr:  domain.ErrPasswordTooWeak,
		},
		{
			name:     "deve retornar erro para senha sem caractere especial",
			password: "Aaaa1bcde",
			wantErr:  domain.ErrPasswordTooWeak,
		},
		{
			name:     "deve retornar erro para senha vazia",
			password: "",
			wantErr:  domain.ErrPasswordTooShort,
		},
		{
			name:     "deve retornar erro para senha apenas com números",
			password: "12345678",
			wantErr:  domain.ErrPasswordTooWeak,
		},
		{
			name:     "deve retornar erro para senha apenas com minúsculas",
			password: "abcdefgh",
			wantErr:  domain.ErrPasswordTooWeak,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validatePasswordStrength(tt.password)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("validatePasswordStrength(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}
