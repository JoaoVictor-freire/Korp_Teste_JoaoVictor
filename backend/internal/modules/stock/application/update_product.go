package application

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"korp_backend/internal/modules/stock/domain"
)

var ErrProductNotFound = errors.New("product not found")

type UpdateProductInput struct {
	OwnerID      string
	OriginalCode string
	Code         string
	Description  string
	Stock        int
}

type UpdateProductUseCase struct {
	repository domain.ProductRepository
}

func NewUpdateProductUseCase(repository domain.ProductRepository) UpdateProductUseCase {
	return UpdateProductUseCase{repository: repository}
}

func (uc UpdateProductUseCase) Execute(ctx context.Context, input UpdateProductInput) (domain.Product, error) {
	if strings.TrimSpace(input.OwnerID) == "" {
		return domain.Product{}, errors.New("owner id is required")
	}

	if strings.TrimSpace(input.OriginalCode) == "" {
		return domain.Product{}, ErrProductCodeRequired
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

	existingProduct, err := uc.repository.GetByOwnerAndCode(ctx, input.OwnerID, input.OriginalCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Product{}, ErrProductNotFound
		}
		return domain.Product{}, err
	}

	if input.OriginalCode != input.Code {
		exists, err := uc.repository.ExistsByOwnerAndCode(ctx, input.OwnerID, input.Code)
		if err != nil {
			return domain.Product{}, err
		}

		if exists {
			return domain.Product{}, ErrProductAlreadyExists
		}
	}

	updatedProduct := domain.Product{
		OwnerID:     existingProduct.OwnerID,
		Code:        input.Code,
		Description: input.Description,
		Stock:       input.Stock,
		CreatedAt:   existingProduct.CreatedAt,
	}

	if err := uc.repository.Update(ctx, input.OriginalCode, updatedProduct); err != nil {
		return domain.Product{}, err
	}

	return updatedProduct, nil
}
