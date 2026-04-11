package application

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"korp_backend/internal/modules/stock/domain"
)

type DeleteProductUseCase struct {
	repository domain.ProductRepository
}

func NewDeleteProductUseCase(repository domain.ProductRepository) DeleteProductUseCase {
	return DeleteProductUseCase{repository: repository}
}

func (uc DeleteProductUseCase) Execute(ctx context.Context, ownerID string, code string) error {
	if strings.TrimSpace(ownerID) == "" {
		return errors.New("owner id is required")
	}

	if strings.TrimSpace(code) == "" {
		return ErrProductCodeRequired
	}

	_, err := uc.repository.GetByOwnerAndCode(ctx, ownerID, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	return uc.repository.Delete(ctx, ownerID, code)
}
