package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		inputName string
		inputEmail string
		inputPass  string
		mockRepo   func() *mockUserRepository
		wantErr  error
	}{
		{
			name:       "deve registrar usuário com sucesso",
			inputName:  "João",
			inputEmail: "joao@email.com",
			inputPass:  "Str0ng!Pass",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
						return nil, domain.ErrUserNotFound
					},
					createFn: func(_ context.Context, _ *domain.User) error {
						return nil
					},
				}
			},
			wantErr: nil,
		},
		{
			name:       "deve retornar erro para senha fraca",
			inputName:  "João",
			inputEmail: "joao@email.com",
			inputPass:  "weak",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{}
			},
			wantErr: domain.ErrPasswordTooShort,
		},
		{
			name:       "deve retornar erro para senha sem especial",
			inputName:  "João",
			inputEmail: "joao@email.com",
			inputPass:  "Str0ngPass",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{}
			},
			wantErr: domain.ErrPasswordTooWeak,
		},
		{
			name:       "deve retornar erro para email já em uso",
			inputName:  "João",
			inputEmail: "joao@email.com",
			inputPass:  "Str0ng!Pass",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
						return &domain.User{}, nil
					},
				}
			},
			wantErr: domain.ErrEmailInUse,
		},
		{
			name:       "deve retornar erro de repo no Create",
			inputName:  "João",
			inputEmail: "joao@email.com",
			inputPass:  "Str0ng!Pass",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
						return nil, domain.ErrUserNotFound
					},
					createFn: func(_ context.Context, _ *domain.User) error {
						return errors.New("db error")
					},
				}
			},
			wantErr: nil, // não é erro sentinela
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.mockRepo()
			svc := NewUserService(repo)

			user, err := svc.Register(context.Background(), tt.inputName, tt.inputEmail, tt.inputPass)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if tt.name == "deve retornar erro de repo no Create" {
				if err == nil {
					t.Error("Register() deveria retornar erro do repo")
				}
				return
			}

			if err != nil {
				t.Errorf("Register() unexpected error = %v", err)
				return
			}

			if user == nil {
				t.Fatal("Register() retornou nil sem erro")
			}
			if user.Name != tt.inputName {
				t.Errorf("Name = %q, want %q", user.Name, tt.inputName)
			}
			if user.Email != tt.inputEmail {
				t.Errorf("Email = %q, want %q", user.Email, tt.inputEmail)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	t.Parallel()

	password := "Str0ng!Pass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name     string
		email    string
		password string
		mockRepo func() *mockUserRepository
		wantErr  error
	}{
		{
			name:     "deve fazer login com sucesso",
			email:    "joao@email.com",
			password: password,
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
						return &domain.User{
							ID:           uuid.New(),
							Name:         "João",
							Email:        "joao@email.com",
							PasswordHash: string(hash),
						}, nil
					},
				}
			},
			wantErr: nil,
		},
		{
			name:     "deve retornar erro para usuário não encontrado",
			email:    "naoexiste@email.com",
			password: password,
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
						return nil, domain.ErrUserNotFound
					},
				}
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name:     "deve retornar erro para senha incorreta",
			email:    "joao@email.com",
			password: "WrongPass1!",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
						return &domain.User{
							PasswordHash: string(hash),
						}, nil
					},
				}
			},
			wantErr: domain.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.mockRepo()
			svc := NewUserService(repo)

			user, err := svc.Login(context.Background(), tt.email, tt.password)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && user == nil {
				t.Error("Login() retornou nil sem erro")
			}
		})
	}
}

func TestUserService_UpdatePassword(t *testing.T) {
	t.Parallel()

	currentPassword := "Str0ng!Pass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(currentPassword), bcrypt.DefaultCost)
	userID := uuid.New()

	tests := []struct {
		name        string
		currentPass string
		newPass     string
		mockRepo    func() *mockUserRepository
		wantErr     error
	}{
		{
			name:        "deve atualizar senha com sucesso",
			currentPass: currentPassword,
			newPass:     "N3w!Passw0rd",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
						return &domain.User{
							ID:           userID,
							PasswordHash: string(hash),
						}, nil
					},
					updateFn: func(_ context.Context, _ *domain.User) error {
						return nil
					},
				}
			},
			wantErr: nil,
		},
		{
			name:        "deve retornar erro para nova senha fraca",
			currentPass: currentPassword,
			newPass:     "weak",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{}
			},
			wantErr: domain.ErrPasswordTooShort,
		},
		{
			name:        "deve retornar erro para senha atual incorreta",
			currentPass: "WrongPass1!",
			newPass:     "N3w!Passw0rd",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
						return &domain.User{
							ID:           userID,
							PasswordHash: string(hash),
						}, nil
					},
				}
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name:        "deve retornar erro para mesma senha",
			currentPass: currentPassword,
			newPass:     currentPassword,
			mockRepo: func() *mockUserRepository {
				// A nova senha é igual à atual, precisa retornar o user com o hash correto
				return &mockUserRepository{
					getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
						return &domain.User{
							ID:           userID,
							PasswordHash: string(hash),
						}, nil
					},
				}
			},
			wantErr: domain.ErrSamePassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.mockRepo()
			svc := NewUserService(repo)

			err := svc.UpdatePassword(context.Background(), userID, tt.currentPass, tt.newPass)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Garante que o hash gerado é bcrypt válido
func TestUserService_Register_GeneratesBcryptHash(t *testing.T) {
	t.Parallel()

	var savedUser *domain.User
	repo := &mockUserRepository{
		getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		},
		createFn: func(_ context.Context, u *domain.User) error {
			savedUser = u
			return nil
		},
	}

	svc := NewUserService(repo)
	_, err := svc.Register(context.Background(), "João", "joao@email.com", "Str0ng!Pass")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if savedUser == nil {
		t.Fatal("user não foi salvo")
	}

	if !strings.HasPrefix(savedUser.PasswordHash, "$2") {
		t.Error("PasswordHash não é bcrypt")
	}
}
