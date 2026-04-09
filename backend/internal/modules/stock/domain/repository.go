package domain

import "context"

type ProductRepository interface {
	Create(ctx context.Context, product Product) error
	ListByOwner(ctx context.Context, ownerID string) ([]Product, error)
	ExistsByOwnerAndCode(ctx context.Context, ownerID string, code string) (bool, error)
}
