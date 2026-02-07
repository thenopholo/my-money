package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/thenopholo/my-money/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (us *UserService) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	if err := validatePasswordStrength(password); err != nil {
		return nil, err
	}

	_, err := us.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrEmailInUse
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(name, email, string(hash))
	if err != nil {
		return nil, err
	}

	if err := us.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := us.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, err
	}

	return user, nil
}

func (us *UserService) UpdatePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	if err := validatePasswordStrength(newPassword); err != nil {
		return err
	}

	user, err := us.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return domain.ErrInvalidCredentials
		}

		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(newPassword)); err == nil {
		return domain.ErrSamePassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)

	return us.userRepo.Update(ctx, user)
}
