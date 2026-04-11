package application

import (
	"context"
	"errors"
	"time"

	"korp_backend/internal/modules/billing/domain"
)

var (
	ErrInvoiceNumberInvalid     = errors.New("invoice number must be greater than zero")
	ErrInvoiceItemsRequired     = errors.New("invoice must contain at least one item")
	ErrInvoiceItemCodeRequired  = errors.New("invoice item product code is required")
	ErrInvoiceItemQuantityError = errors.New("invoice item quantity must be greater than zero")
	ErrInvoiceAlreadyExists     = errors.New("invoice already exists")
)

type CreateInvoiceInput struct {
	OwnerID string
	Number  int
	Items   []domain.InvoiceItem
}

type CreateInvoiceUseCase struct {
	repository domain.InvoiceRepository
}

func NewCreateInvoiceUseCase(repository domain.InvoiceRepository) CreateInvoiceUseCase {
	return CreateInvoiceUseCase{repository: repository}
}

func (uc CreateInvoiceUseCase) Execute(ctx context.Context, input CreateInvoiceInput) (domain.Invoice, error) {
	if err := validateInvoiceInput(input.OwnerID, input.Number, input.Items); err != nil {
		return domain.Invoice{}, err
	}

	exists, err := uc.repository.ExistsByOwnerAndNumber(ctx, input.OwnerID, input.Number)
	if err != nil {
		return domain.Invoice{}, err
	}

	if exists {
		return domain.Invoice{}, ErrInvoiceAlreadyExists
	}

	invoice := domain.Invoice{
		OwnerID:   input.OwnerID,
		Number:    input.Number,
		Status:    domain.StatusOpen,
		Items:     input.Items,
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.repository.Create(ctx, invoice); err != nil {
		return domain.Invoice{}, err
	}

	return invoice, nil
}
