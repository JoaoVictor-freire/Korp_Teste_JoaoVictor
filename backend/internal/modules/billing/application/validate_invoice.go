package application

import (
	"errors"

	"korp_backend/internal/modules/billing/domain"
)

func validateInvoiceInput(ownerID string, number int, items []domain.InvoiceItem) error {
	if ownerID == "" {
		return errors.New("owner id is required")
	}

	if number <= 0 {
		return ErrInvoiceNumberInvalid
	}

	if len(items) == 0 {
		return ErrInvoiceItemsRequired
	}

	for _, item := range items {
		if item.ProductCode == "" {
			return ErrInvoiceItemCodeRequired
		}

		if item.Quantity <= 0 {
			return ErrInvoiceItemQuantityError
		}
	}

	return nil
}
