package application

import (
	"context"
	"errors"

	"korp_backend/internal/modules/billing/domain"
	stockdomain "korp_backend/internal/modules/stock/domain"
)

var (
	ErrCloseInvoiceOwnerRequired = errors.New("owner id is required")
	ErrCloseInvoiceNotFound      = errors.New("invoice not found")
	ErrCloseInvoiceProductNotFound = errors.New("invoice product not found in stock")
)

type CloseInvoiceUseCase struct {
	repository domain.InvoiceRepository
	productRepository stockdomain.ProductRepository
}

type CloseInvoiceInput struct {
	OwnerID string
	Number  int
}

func NewCloseInvoiceUseCase(repository domain.InvoiceRepository, productRepository stockdomain.ProductRepository) CloseInvoiceUseCase {
	return CloseInvoiceUseCase{
		repository: repository,
		productRepository: productRepository,
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

	for _, item := range invoice.Items {
		product, err := uc.productRepository.GetByOwnerAndCode(ctx, input.OwnerID, item.ProductCode)
		if err != nil {
			return ErrCloseInvoiceProductNotFound
		}

		if err := product.DecreaseStock(item.Quantity); err != nil {
			return err
		}

		if err := uc.productRepository.UpdateStock(ctx, input.OwnerID, item.ProductCode, product.Stock); err != nil {
			return err
		}
	}

	newStatus := invoice.Status == domain.StatusOpen

	if err := uc.repository.UpdateStatus(ctx, invoice.Number, invoice.OwnerID, newStatus); err != nil {
		return err
	}

	return nil
}
