package application

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"korp_backend/internal/modules/billing/domain"
)

var ErrInvoiceCannotEditClosed = errors.New("closed invoices cannot be edited")

type UpdateInvoiceInput struct {
	OwnerID        string
	OriginalNumber int
	Number         int
	Items          []domain.InvoiceItem
}

type UpdateInvoiceUseCase struct {
	repository domain.InvoiceRepository
}

func NewUpdateInvoiceUseCase(repository domain.InvoiceRepository) UpdateInvoiceUseCase {
	return UpdateInvoiceUseCase{repository: repository}
}

func (uc UpdateInvoiceUseCase) Execute(ctx context.Context, input UpdateInvoiceInput) (domain.Invoice, error) {
	if err := validateInvoiceInput(input.OwnerID, input.Number, input.Items); err != nil {
		return domain.Invoice{}, err
	}

	if input.OriginalNumber <= 0 {
		return domain.Invoice{}, ErrInvoiceNumberInvalid
	}

	existingInvoice, err := uc.repository.GetByOwnerAndNumber(ctx, input.OwnerID, input.OriginalNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Invoice{}, ErrCloseInvoiceNotFound
		}
		return domain.Invoice{}, err
	}

	if existingInvoice.Status == domain.StatusClosed {
		return domain.Invoice{}, ErrInvoiceCannotEditClosed
	}

	if input.OriginalNumber != input.Number {
		exists, err := uc.repository.ExistsByOwnerAndNumber(ctx, input.OwnerID, input.Number)
		if err != nil {
			return domain.Invoice{}, err
		}

		if exists {
			return domain.Invoice{}, ErrInvoiceAlreadyExists
		}
	}

	updatedInvoice := domain.Invoice{
		OwnerID:   existingInvoice.OwnerID,
		Number:    input.Number,
		Status:    existingInvoice.Status,
		Items:     input.Items,
		CreatedAt: existingInvoice.CreatedAt,
	}

	if err := uc.repository.Update(ctx, input.OriginalNumber, updatedInvoice); err != nil {
		return domain.Invoice{}, err
	}

	return updatedInvoice, nil
}
