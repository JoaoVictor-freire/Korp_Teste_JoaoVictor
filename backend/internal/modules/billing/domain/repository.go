package domain

import "context"

type InvoiceRepository interface {
	Create(ctx context.Context, invoice Invoice) error
	ListByOwner(ctx context.Context, ownerID string) ([]Invoice, error)
	ExistsByOwnerAndNumber(ctx context.Context, ownerID string, number int) (bool, error)
	GetByOwnerAndNumber(ctx context.Context, ownerID string, number int) (Invoice, error)
	Update(ctx context.Context, originalNumber int, invoice Invoice) error
	UpdateStatusIfCurrent(ctx context.Context, number int, ownerID string, currentStatus bool, newStatus bool) (bool, error)
	Delete(ctx context.Context, ownerID string, number int) error
}
