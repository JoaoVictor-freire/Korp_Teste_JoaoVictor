package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"korp_backend/internal/modules/users/domain"
)

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrUserExists       = errors.New("user already exists")
)

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterUseCase struct {
	repository domain.UserRepository
}

func NewRegisterUseCase(repository domain.UserRepository) RegisterUseCase {
	return RegisterUseCase{repository: repository}
}

func (uc RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" {
		return domain.User{}, ErrEmailRequired
	}

	if strings.TrimSpace(input.Password) == "" {
		return domain.User{}, ErrPasswordRequired
	}

	_, exists, err := uc.repository.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	if exists {
		return domain.User{}, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		ID:           newID(),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC(),
	}

	if err := uc.repository.Create(ctx, user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}
