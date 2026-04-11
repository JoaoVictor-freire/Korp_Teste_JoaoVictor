package application

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"korp_backend/internal/modules/billing/domain"
)

var ErrInvoiceCannotDeleteClosed = errors.New("closed invoices cannot be deleted")

type DeleteInvoiceUseCase struct {
	repository domain.InvoiceRepository
}

func NewDeleteInvoiceUseCase(repository domain.InvoiceRepository) DeleteInvoiceUseCase {
	return DeleteInvoiceUseCase{repository: repository}
}

func (uc DeleteInvoiceUseCase) Execute(ctx context.Context, ownerID string, number int) error {
	if ownerID == "" {
		return errors.New("owner id is required")
	}

	if number <= 0 {
		return ErrInvoiceNumberInvalid
	}

	invoice, err := uc.repository.GetByOwnerAndNumber(ctx, ownerID, number)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCloseInvoiceNotFound
		}
		return err
	}

	if invoice.Status == domain.StatusClosed {
		return ErrInvoiceCannotDeleteClosed
	}

	return uc.repository.Delete(ctx, ownerID, number)
}
