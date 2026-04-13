package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	nethttp "net/http"
	"strings"
	"sync"
	"time"

	"korp_backend/internal/modules/billing/application"
)

type StockClient struct {
	baseURL string
	client  *nethttp.Client
	breaker *circuitBreaker
}

type StockClientConfig struct {
	RequestTimeout   time.Duration
	FailureThreshold int
	ResetTimeout     time.Duration
}

type circuitState string

const (
	circuitClosed   circuitState = "closed"
	circuitOpen     circuitState = "open"
	circuitHalfOpen circuitState = "half-open"
)

type circuitBreaker struct {
	mu                  sync.Mutex
	state               circuitState
	consecutiveFailures int
	openedAt            time.Time
	probeInFlight       bool
	failureThreshold    int
	resetTimeout        time.Duration
}

type stockEnvelope[T any] struct {
	Data T `json:"data"`
}

type stockErrorEnvelope struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type stockProductResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
}

type decreaseStockRequest struct {
	Quantity int `json:"quantity"`
}

func NewStockClient(baseURL string, cfg StockClientConfig) *StockClient {
	requestTimeout := cfg.RequestTimeout
	if requestTimeout <= 0 {
		requestTimeout = 5 * time.Second
	}

	failureThreshold := cfg.FailureThreshold
	if failureThreshold <= 0 {
		failureThreshold = 3
	}

	resetTimeout := cfg.ResetTimeout
	if resetTimeout <= 0 {
		resetTimeout = 15 * time.Second
	}

	return &StockClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &nethttp.Client{
			Timeout: requestTimeout,
		},
		breaker: &circuitBreaker{
			state:            circuitClosed,
			failureThreshold: failureThreshold,
			resetTimeout:     resetTimeout,
		},
	}
}

func (c *StockClient) GetProduct(ctx context.Context, code string) (application.StockProduct, error) {
	var envelope stockEnvelope[stockProductResponse]
	if err := c.doJSON(ctx, nethttp.MethodGet, fmt.Sprintf("/api/v1/products/%s", code), nil, &envelope); err != nil {
		return application.StockProduct{}, err
	}

	return application.StockProduct{
		Code:        envelope.Data.Code,
		Description: envelope.Data.Description,
		Stock:       envelope.Data.Stock,
	}, nil
}

func (c *StockClient) DecreaseStock(ctx context.Context, code string, quantity int) error {
	return c.doJSON(ctx, nethttp.MethodPatch, fmt.Sprintf("/api/v1/products/%s/decrease", code), decreaseStockRequest{
		Quantity: quantity,
	}, nil)
}

func (c *StockClient) doJSON(ctx context.Context, method string, path string, body any, out any) error {
	authHeader, ok := application.AuthorizationHeaderFromContext(ctx)
	if !ok {
		return application.ErrStockUnauthorized
	}

	if err := c.breaker.allowRequest(); err != nil {
		return err
	}

	var payload io.Reader
	if body != nil {
		rawBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewReader(rawBody)
	}

	req, err := nethttp.NewRequestWithContext(ctx, method, c.baseURL+path, payload)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", authHeader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.breaker.recordFailure()
		return application.ErrStockUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		mappedErr := c.mapError(resp)
		if c.shouldTripBreaker(mappedErr) {
			c.breaker.recordFailure()
		} else {
			c.breaker.recordSuccess()
		}
		return mappedErr
	}

	if out == nil {
		c.breaker.recordSuccess()
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		c.breaker.recordFailure()
		return application.ErrStockUnavailable
	}

	c.breaker.recordSuccess()
	return nil
}

func (c *StockClient) mapError(resp *nethttp.Response) error {
	var envelope stockErrorEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		switch resp.StatusCode {
		case nethttp.StatusUnauthorized:
			return application.ErrStockUnauthorized
		case nethttp.StatusNotFound:
			return application.ErrStockProductNotFound
		case nethttp.StatusConflict:
			return application.ErrStockInsufficient
		default:
			return application.ErrStockUnavailable
		}
	}

	message := strings.ToLower(strings.TrimSpace(envelope.Error.Message))

	switch resp.StatusCode {
	case nethttp.StatusUnauthorized:
		return application.ErrStockUnauthorized
	case nethttp.StatusBadRequest:
		return errors.New(envelope.Error.Message)
	case nethttp.StatusNotFound:
		return application.ErrStockProductNotFound
	case nethttp.StatusConflict:
		if strings.Contains(message, "stock") {
			return application.ErrStockInsufficient
		}
		return errors.New(envelope.Error.Message)
	case nethttp.StatusServiceUnavailable,
		nethttp.StatusBadGateway,
		nethttp.StatusGatewayTimeout,
		nethttp.StatusInternalServerError:
		return application.ErrStockUnavailable
	default:
		return errors.New(envelope.Error.Message)
	}
}

func (c *StockClient) shouldTripBreaker(err error) bool {
	return errors.Is(err, application.ErrStockUnavailable)
}

func (b *circuitBreaker) allowRequest() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case circuitClosed:
		return nil
	case circuitOpen:
		if time.Since(b.openedAt) < b.resetTimeout {
			return application.ErrStockCircuitOpen
		}

		b.state = circuitHalfOpen
		b.probeInFlight = true
		return nil
	case circuitHalfOpen:
		if b.probeInFlight {
			return application.ErrStockCircuitOpen
		}

		b.probeInFlight = true
		return nil
	default:
		return nil
	}
}

func (b *circuitBreaker) recordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.state = circuitClosed
	b.consecutiveFailures = 0
	b.openedAt = time.Time{}
	b.probeInFlight = false
}

func (b *circuitBreaker) recordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.probeInFlight = false

	if b.state == circuitHalfOpen {
		b.state = circuitOpen
		b.openedAt = time.Now()
		b.consecutiveFailures = b.failureThreshold
		return
	}

	b.consecutiveFailures++
	if b.consecutiveFailures < b.failureThreshold {
		return
	}

	b.state = circuitOpen
	b.openedAt = time.Now()
}
