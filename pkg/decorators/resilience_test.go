package decorators_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"markets-sdk/pkg/decorators"
	"markets-sdk/pkg/domain"
)

func TestCircuitBreaker(t *testing.T) {
	ctx := context.Background()
	e := errors.New("failure")

	// 1. Setup: 2 failures to open
	mock := &MockProvider{
		QuoteFn: func(ctx context.Context, symbol string) (*domain.Quote, error) {
			return nil, e
		},
	}

	cb := decorators.NewCircuitBreaker(mock, 2, 100*time.Millisecond)

	// 2. Fail twice
	_, err := cb.GetQuote(ctx, "BTC")
	if err == nil {
		t.Error("expected error")
	}
	_, err = cb.GetQuote(ctx, "BTC")
	if err == nil {
		t.Error("expected error")
	}

	// 3. Should now be open
	_, err = cb.GetQuote(ctx, "BTC")
	if err == nil || err.Error() != "circuit breaker is open" {
		t.Errorf("expected circuit breaker open error, got %v", err)
	}

	// 4. Wait for reset
	time.Sleep(150 * time.Millisecond)

	// 5. Should allow one request (Half-Open)
	// Let's make it succeed this time
	mock.QuoteFn = func(ctx context.Context, symbol string) (*domain.Quote, error) {
		return &domain.Quote{Price: 100}, nil
	}

	q, err := cb.GetQuote(ctx, "BTC")
	if err != nil {
		t.Errorf("expected success in half-open state, got %v", err)
	}
	if q.Price != 100 {
		t.Error("expected price 100")
	}

	// 6. Should be closed now (allow more requests)
	mock.QuoteFn = func(ctx context.Context, symbol string) (*domain.Quote, error) {
		return &domain.Quote{Price: 200}, nil
	}
	q, err = cb.GetQuote(ctx, "BTC")
	if err != nil {
		t.Errorf("expected success in closed state, got %v", err)
	}
}

func TestRetry(t *testing.T) {
	ctx := context.Background()

	// 1. Fail twice, then succeed
	calls := 0
	mock := &MockProvider{
		QuoteFn: func(ctx context.Context, symbol string) (*domain.Quote, error) {
			calls++
			if calls <= 2 {
				return nil, errors.New("fail")
			}
			return &domain.Quote{Price: 100}, nil
		},
	}

	// Max 3 retries
	r := decorators.NewRetry(mock, 3, 1*time.Millisecond)

	q, err := r.GetQuote(ctx, "BTC")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if q.Price != 100 {
		t.Errorf("expected 100, got %v", q.Price)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}
