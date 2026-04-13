package application

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"korp_backend/internal/modules/stock/domain"
)

type GetProductUseCase struct {
	repository domain.ProductRepository
}

func NewGetProductUseCase(repository domain.ProductRepository) GetProductUseCase {
	return GetProductUseCase{repository: repository}
}

func (uc GetProductUseCase) Execute(ctx context.Context, ownerID string, code string) (domain.Product, error) {
	if strings.TrimSpace(ownerID) == "" {
		return domain.Product{}, errors.New("owner id is required")
	}

	if strings.TrimSpace(code) == "" {
		return domain.Product{}, ErrProductCodeRequired
	}

	product, err := uc.repository.GetByOwnerAndCode(ctx, ownerID, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Product{}, ErrProductNotFound
		}
		return domain.Product{}, err
	}

	return product, nil
}
