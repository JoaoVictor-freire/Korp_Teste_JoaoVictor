package memory

import (
	"context"
	"slices"
	"sync"

	"korp_backend/internal/modules/stock/domain"
)

type ProductRepository struct {
	mu       sync.RWMutex
	products map[string]domain.Product
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make(map[string]domain.Product),
	}
}

func (r *ProductRepository) Create(_ context.Context, product domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[ownerKey(product.OwnerID, product.Code)] = product
	return nil
}

func (r *ProductRepository) ListByOwner(_ context.Context, ownerID string) ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]domain.Product, 0, len(r.products))
	for _, product := range r.products {
		if product.OwnerID == ownerID {
			products = append(products, product)
		}
	}

	slices.SortFunc(products, func(a, b domain.Product) int {
		switch {
		case a.Code < b.Code:
			return -1
		case a.Code > b.Code:
			return 1
		default:
			return 0
		}
	})

	return products, nil
}

func (r *ProductRepository) ExistsByOwnerAndCode(_ context.Context, ownerID string, code string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.products[ownerKey(ownerID, code)]
	return exists, nil
}

func ownerKey(ownerID string, code string) string {
	return ownerID + "|" + code
}
