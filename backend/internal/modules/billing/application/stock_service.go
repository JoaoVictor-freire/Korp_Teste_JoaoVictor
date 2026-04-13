package application

import (
	"context"
	"errors"
)

var (
	ErrStockProductNotFound = errors.New("product not found in stock")
	ErrStockUnavailable     = errors.New("stock service unavailable")
	ErrStockUnauthorized    = errors.New("stock service unauthorized")
	ErrStockInsufficient    = errors.New("insufficient stock")
	ErrStockCircuitOpen     = errors.New("stock service circuit breaker is open")
)

type StockProduct struct {
	Code        string
	Description string
	Stock       int
}

type StockService interface {
	GetProduct(ctx context.Context, code string) (StockProduct, error)
	DecreaseStock(ctx context.Context, code string, quantity int) error
}
