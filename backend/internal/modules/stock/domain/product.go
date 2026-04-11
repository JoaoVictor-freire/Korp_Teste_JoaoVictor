package domain

import (
	"errors"
	"time"
)

type Product struct {
	OwnerID     string    `json:"owner_id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}

var ErrInsufficientStock = errors.New("insufficient stock")

func (p *Product) DecreaseStock(quantity int) error {
	if quantity <= 0 {
		return nil
	}

	if p.Stock < quantity {
		return ErrInsufficientStock
	}

	p.Stock -= quantity
	return nil
}
