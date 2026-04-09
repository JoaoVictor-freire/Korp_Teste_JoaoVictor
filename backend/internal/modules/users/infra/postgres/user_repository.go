package postgres

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"korp_backend/internal/modules/users/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	model := UserModel{
		ID:           user.ID,
		Name:         user.Name,
		Email:        strings.ToLower(user.Email),
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	user.ID = model.ID
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, bool, error) {
	var model UserModel
	err := r.db.WithContext(ctx).
		Where("email = ?", strings.ToLower(email)).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, false, nil
		}
		return domain.User{}, false, err
	}

	return domain.User{
		ID:           model.ID,
		Name:         model.Name,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
		CreatedAt:    model.CreatedAt,
	}, true, nil
}
