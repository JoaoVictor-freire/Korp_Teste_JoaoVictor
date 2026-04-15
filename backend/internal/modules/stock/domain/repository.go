package domain

import "context"

type ProductRepository interface {
	Create(ctx context.Context, product Product) error
	ListByOwner(ctx context.Context, ownerID string) ([]Product, error)
	ExistsByOwnerAndCode(ctx context.Context, ownerID string, code string) (bool, error)
	GetByOwnerAndCode(ctx context.Context, ownerID string, code string) (Product, error)
	Update(ctx context.Context, originalCode string, product Product) error
	UpdateStock(ctx context.Context, ownerID string, code string, newStock int) error
	DecreaseStock(ctx context.Context, ownerID string, code string, quantity int) (bool, error)
	Delete(ctx context.Context, ownerID string, code string) error
}
