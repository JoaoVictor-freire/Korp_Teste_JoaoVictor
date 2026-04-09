package memory

import (
	"context"
	"slices"
	"strconv"
	"sync"

	"korp_backend/internal/modules/billing/domain"
)

type InvoiceRepository struct {
	mu       sync.RWMutex
	invoices map[string]domain.Invoice
}

func NewInvoiceRepository() *InvoiceRepository {
	return &InvoiceRepository{
		invoices: make(map[string]domain.Invoice),
	}
}

func (r *InvoiceRepository) Create(_ context.Context, invoice domain.Invoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.invoices[ownerKey(invoice.OwnerID, invoice.Number)] = invoice
	return nil
}

func (r *InvoiceRepository) ListByOwner(_ context.Context, ownerID string) ([]domain.Invoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	invoices := make([]domain.Invoice, 0, len(r.invoices))
	for _, invoice := range r.invoices {
		if invoice.OwnerID == ownerID {
			invoices = append(invoices, invoice)
		}
	}

	slices.SortFunc(invoices, func(a, b domain.Invoice) int {
		switch {
		case a.Number < b.Number:
			return -1
		case a.Number > b.Number:
			return 1
		default:
			return 0
		}
	})

	return invoices, nil
}

func (r *InvoiceRepository) ExistsByOwnerAndNumber(_ context.Context, ownerID string, number int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.invoices[ownerKey(ownerID, number)]
	return exists, nil
}

func ownerKey(ownerID string, number int) string {
	return ownerID + "|" + strconv.Itoa(number)
}
