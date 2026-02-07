package service

import (
	"unicode"

	"github.com/thenopholo/my-money/internal/domain"
)

func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return domain.ErrPasswordTooShort
	}

	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return domain.ErrPasswordTooWeak
	}

	return nil
}
