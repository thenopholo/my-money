package domain

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(name, email, passwordHash string) (*User, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	if len(passwordHash) < 50 {
		return nil, ErrPasswordTooWeak
	}


	return &User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}