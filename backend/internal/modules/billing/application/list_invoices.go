package application

import (
	"context"

	"korp_backend/internal/modules/billing/domain"
)

type ListInvoicesUseCase struct {
	repository domain.InvoiceRepository
}

func NewListInvoicesUseCase(repository domain.InvoiceRepository) ListInvoicesUseCase {
	return ListInvoicesUseCase{repository: repository}
}

func (uc ListInvoicesUseCase) Execute(ctx context.Context, ownerID string) ([]domain.Invoice, error) {
	return uc.repository.ListByOwner(ctx, ownerID)
}
