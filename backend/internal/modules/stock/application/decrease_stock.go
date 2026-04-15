package application

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"korp_backend/internal/modules/stock/domain"
)

var ErrDecreaseStockQuantityInvalid = errors.New("quantity must be greater than zero")

type DecreaseStockInput struct {
	OwnerID  string
	Code     string
	Quantity int
}

type DecreaseStockUseCase struct {
	repository domain.ProductRepository
}

func NewDecreaseStockUseCase(repository domain.ProductRepository) DecreaseStockUseCase {
	return DecreaseStockUseCase{repository: repository}
}

func (uc DecreaseStockUseCase) Execute(ctx context.Context, input DecreaseStockInput) (domain.Product, error) {
	if strings.TrimSpace(input.OwnerID) == "" {
		return domain.Product{}, errors.New("owner id is required")
	}

	if strings.TrimSpace(input.Code) == "" {
		return domain.Product{}, ErrProductCodeRequired
	}

	if input.Quantity <= 0 {
		return domain.Product{}, ErrDecreaseStockQuantityInvalid
	}

	updated, err := uc.repository.DecreaseStock(ctx, input.OwnerID, input.Code, input.Quantity)
	if err != nil {
		return domain.Product{}, err
	}

	product, err := uc.repository.GetByOwnerAndCode(ctx, input.OwnerID, input.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Product{}, ErrProductNotFound
		}
		return domain.Product{}, err
	}

	if !updated {
		return domain.Product{}, domain.ErrInsufficientStock
	}

	return product, nil
}
