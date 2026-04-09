package application

import (
	"context"

	"korp_backend/internal/modules/stock/domain"
)

type ListProductsUseCase struct {
	repository domain.ProductRepository
}

func NewListProductsUseCase(repository domain.ProductRepository) ListProductsUseCase {
	return ListProductsUseCase{repository: repository}
}

func (uc ListProductsUseCase) Execute(ctx context.Context, ownerID string) ([]domain.Product, error) {
	return uc.repository.ListByOwner(ctx, ownerID)
}
