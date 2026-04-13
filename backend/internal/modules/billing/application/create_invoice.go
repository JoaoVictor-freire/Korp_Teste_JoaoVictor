package application

import (
	"context"
	"errors"
	"fmt"
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
	repository   domain.InvoiceRepository
	stockService StockService
}

type InvoiceOutOfStockError struct {
	ProductCode string
}

func (e InvoiceOutOfStockError) Error() string {
	return fmt.Sprintf("product %s is out of stock", e.ProductCode)
}

type InvoiceProductNotFoundError struct {
	ProductCode string
}

func (e InvoiceProductNotFoundError) Error() string {
	return fmt.Sprintf("product %s was not found in stock", e.ProductCode)
}

func NewCreateInvoiceUseCase(repository domain.InvoiceRepository, stockService StockService) CreateInvoiceUseCase {
	return CreateInvoiceUseCase{
		repository:   repository,
		stockService: stockService,
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
		product, err := uc.stockService.GetProduct(ctx, item.ProductCode)
		if err != nil {
			if errors.Is(err, ErrStockProductNotFound) {
				return domain.Invoice{}, InvoiceProductNotFoundError{ProductCode: item.ProductCode}
			}
			return domain.Invoice{}, err
		}

		if product.Stock < item.Quantity {
			return domain.Invoice{}, InvoiceOutOfStockError{ProductCode: item.ProductCode}
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
