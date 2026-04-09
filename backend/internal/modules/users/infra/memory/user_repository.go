package memory

import (
	"context"
	"strings"
	"sync"

	"korp_backend/internal/modules/users/domain"
)

type UserRepository struct {
	mu      sync.RWMutex
	byEmail map[string]domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		byEmail: make(map[string]domain.User),
	}
}

func (r *UserRepository) Create(_ context.Context, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byEmail[strings.ToLower(user.Email)] = user
	return nil
}

func (r *UserRepository) FindByEmail(_ context.Context, email string) (domain.User, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byEmail[strings.ToLower(email)]
	return user, ok, nil
}
