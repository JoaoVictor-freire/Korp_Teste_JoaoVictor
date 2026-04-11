package domain

import (
	"errors"
	"time"
)

const (
	StatusOpen   = "OPEN"
	StatusClosed = "CLOSED"
)

type Invoice struct {
	OwnerID   string        `json:"owner_id"`
	Number    int           `json:"number"`
	Status    string        `json:"status"`
	Items     []InvoiceItem `json:"items"`
	CreatedAt time.Time     `json:"created_at"`
}

type InvoiceItem struct {
	ProductCode string `json:"product_code"`
	Quantity    int    `json:"quantity"`
}

var ErrInvoiceAlreadyClosed = errors.New("invoice already closed")

func (i *Invoice) Close() error {
	if i.Status == StatusClosed {
		return ErrInvoiceAlreadyClosed
	}
	i.Status = StatusClosed

	return nil
}
