package domain

import "context"

type InvoiceRepository interface {
	Create(ctx context.Context, invoice Invoice) error
	ListByOwner(ctx context.Context, ownerID string) ([]Invoice, error)
	ExistsByOwnerAndNumber(ctx context.Context, ownerID string, number int) (bool, error)
	GetByOwnerAndNumber(ctx context.Context, ownerID string, number int) (Invoice, error)
	UpdateStatus(ctx context.Context, number int, ownerID string, newStatus bool) error
}
