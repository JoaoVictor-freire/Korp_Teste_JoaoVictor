package application

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"korp_backend/internal/modules/users/domain"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginInput struct {
	Email    string
	Password string
}

type LoginUseCase struct {
	repository domain.UserRepository
}

func NewLoginUseCase(repository domain.UserRepository) LoginUseCase {
	return LoginUseCase{repository: repository}
}

func (uc LoginUseCase) Execute(ctx context.Context, input LoginInput) (domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.TrimSpace(input.Password)
	if email == "" || password == "" {
		return domain.User{}, ErrInvalidCredentials
	}

	user, exists, err := uc.repository.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	if !exists {
		return domain.User{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.User{}, ErrInvalidCredentials
	}

	return user, nil
}
