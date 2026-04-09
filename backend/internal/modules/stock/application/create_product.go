package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"korp_backend/internal/modules/stock/domain"
)

var (
	ErrProductCodeRequired        = errors.New("product code is required")
	ErrProductDescriptionRequired = errors.New("product description is required")
	ErrProductStockInvalid        = errors.New("product stock must be zero or greater")
	ErrProductAlreadyExists       = errors.New("product already exists")
)

type CreateProductInput struct {
	OwnerID     string
	Code        string
	Description string
	Stock       int
}

type CreateProductUseCase struct {
	repository domain.ProductRepository
}

func NewCreateProductUseCase(repository domain.ProductRepository) CreateProductUseCase {
	return CreateProductUseCase{repository: repository}
}

func (uc CreateProductUseCase) Execute(ctx context.Context, input CreateProductInput) (domain.Product, error) {
	if strings.TrimSpace(input.OwnerID) == "" {
		return domain.Product{}, errors.New("owner id is required")
	}

	if strings.TrimSpace(input.Code) == "" {
		return domain.Product{}, ErrProductCodeRequired
	}

	if strings.TrimSpace(input.Description) == "" {
		return domain.Product{}, ErrProductDescriptionRequired
	}

	if input.Stock < 0 {
		return domain.Product{}, ErrProductStockInvalid
	}

	exists, err := uc.repository.ExistsByOwnerAndCode(ctx, input.OwnerID, input.Code)
	if err != nil {
		return domain.Product{}, err
	}

	if exists {
		return domain.Product{}, ErrProductAlreadyExists
	}

	product := domain.Product{
		OwnerID:     input.OwnerID,
		Code:        input.Code,
		Description: input.Description,
		Stock:       input.Stock,
		CreatedAt:   time.Now().UTC(),
	}

	if err := uc.repository.Create(ctx, product); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}
