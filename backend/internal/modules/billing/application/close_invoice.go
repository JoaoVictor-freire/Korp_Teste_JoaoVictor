package application

import (
	"context"
	"errors"
	"fmt"

	"korp_backend/internal/modules/billing/domain"
)

var (
	ErrCloseInvoiceOwnerRequired   = errors.New("owner id is required")
	ErrCloseInvoiceNotFound        = errors.New("invoice not found")
	ErrCloseInvoiceProductNotFound = errors.New("invoice product not found in stock")
)

type CloseInvoiceUseCase struct {
	repository   domain.InvoiceRepository
	stockService StockService
}

type CloseInvoiceInput struct {
	OwnerID string
	Number  int
}

type InsufficientStockError struct {
	ProductCode string
}

func (e InsufficientStockError) Error() string {
	return fmt.Sprintf("product %s is out of stock", e.ProductCode)
}

func NewCloseInvoiceUseCase(repository domain.InvoiceRepository, stockService StockService) CloseInvoiceUseCase {
	return CloseInvoiceUseCase{
		repository:   repository,
		stockService: stockService,
	}
}

func (uc CloseInvoiceUseCase) Execute(ctx context.Context, input CloseInvoiceInput) error {
	if input.Number <= 0 {
		return ErrInvoiceNumberInvalid
	}

	if input.OwnerID == "" {
		return ErrCloseInvoiceOwnerRequired
	}

	invoice, err := uc.repository.GetByOwnerAndNumber(ctx, input.OwnerID, input.Number)

	if err != nil {
		return ErrCloseInvoiceNotFound
	}

	if err := invoice.Close(); err != nil {
		return err
	}

	claimed, err := uc.repository.UpdateStatusIfCurrent(ctx, invoice.Number, invoice.OwnerID, true, false)
	if err != nil {
		return err
	}
	if !claimed {
		return domain.ErrInvoiceAlreadyClosed
	}

	for _, item := range invoice.Items {
		err := uc.stockService.DecreaseStock(ctx, item.ProductCode, item.Quantity)
		if errors.Is(err, ErrStockProductNotFound) {
			_, _ = uc.repository.UpdateStatusIfCurrent(ctx, invoice.Number, invoice.OwnerID, false, true)
			return ErrCloseInvoiceProductNotFound
		}
		if errors.Is(err, ErrStockInsufficient) {
			_, _ = uc.repository.UpdateStatusIfCurrent(ctx, invoice.Number, invoice.OwnerID, false, true)
			return InsufficientStockError{ProductCode: item.ProductCode}
		}
		if err != nil {
			_, _ = uc.repository.UpdateStatusIfCurrent(ctx, invoice.Number, invoice.OwnerID, false, true)
			return err
		}
	}

	return nil
}
