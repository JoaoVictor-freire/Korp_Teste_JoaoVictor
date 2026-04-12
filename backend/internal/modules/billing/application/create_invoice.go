package application

import (
	"context"
	"errors"
	"time"

	"korp_backend/internal/modules/billing/domain"
	stockdomain "korp_backend/internal/modules/stock/domain"
)

var (
	ErrInvoiceNumberInvalid       = errors.New("invoice number must be greater than zero")
	ErrInvoiceItemsRequired       = errors.New("invoice must contain at least one item")
	ErrInvoiceItemCodeRequired    = errors.New("invoice item product code is required")
	ErrInvoiceItemQuantityError   = errors.New("invoice item quantity must be greater than zero")
	ErrInvoiceAlreadyExists       = errors.New("invoice already exists")
	ErrInvoiceWithOutStockProduct = errors.New("one of your products hasnt stock for the bill")
)

type CreateInvoiceInput struct {
	OwnerID string
	Number  int
	Items   []domain.InvoiceItem
}

type CreateInvoiceUseCase struct {
	repository        domain.InvoiceRepository
	productRepository stockdomain.ProductRepository
}

func NewCreateInvoiceUseCase(repository domain.InvoiceRepository, productRepository stockdomain.ProductRepository) CreateInvoiceUseCase {
	return CreateInvoiceUseCase{
		repository:        repository,
		productRepository: productRepository,
	}
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

	for _, item := range input.Items {
		product, err := uc.productRepository.GetByOwnerAndCode(ctx, input.OwnerID, item.ProductCode)
		if err != nil {
			return domain.Invoice{}, err
		}

		if product.Stock < item.Quantity {
			return domain.Invoice{}, ErrInvoiceWithOutStockProduct
		}
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
