package application

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"korp_backend/internal/modules/billing/domain"
)

type GetInvoiceUseCase struct {
	repository domain.InvoiceRepository
}

func NewGetInvoiceUseCase(repository domain.InvoiceRepository) GetInvoiceUseCase {
	return GetInvoiceUseCase{repository: repository}
}

func (uc GetInvoiceUseCase) Execute(ctx context.Context, ownerID string, number int) (domain.Invoice, error) {
	if ownerID == "" {
		return domain.Invoice{}, errors.New("owner id is required")
	}

	if number <= 0 {
		return domain.Invoice{}, ErrInvoiceNumberInvalid
	}

	invoice, err := uc.repository.GetByOwnerAndNumber(ctx, ownerID, number)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Invoice{}, ErrCloseInvoiceNotFound
		}
		return domain.Invoice{}, err
	}

	return invoice, nil
}
