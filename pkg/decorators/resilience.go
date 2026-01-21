package decorators

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"markets-sdk/pkg/domain"
	"markets-sdk/pkg/ports"
)

// CircuitBreaker is a decorator that implements the Circuit Breaker pattern
type CircuitBreaker struct {
	provider         ports.Provider
	failureThreshold int
	resetTimeout     time.Duration

	mu              sync.Mutex
	failures        int
	state           state
	lastFailureTime time.Time
}

type state int

const (
	stateClosed state = iota
	stateOpen
	stateHalfOpen
)

func NewCircuitBreaker(provider ports.Provider, failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		provider:         provider,
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
		state:            stateClosed,
	}
}

func (cb *CircuitBreaker) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
	cb.mu.Lock()
	if cb.state == stateOpen {
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = stateHalfOpen
		} else {
			cb.mu.Unlock()
			return nil, fmt.Errorf("circuit breaker is open")
		}
	}
	cb.mu.Unlock()

	quote, err := cb.provider.GetQuote(ctx, symbol)

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailureTime = time.Now()
		if cb.failures >= cb.failureThreshold {
			cb.state = stateOpen
		}
		return nil, err
	}

	// Success
	if cb.state == stateHalfOpen {
		cb.state = stateClosed
		cb.failures = 0
	} else if cb.state == stateClosed {
		cb.failures = 0
	}

	return quote, nil
}

// Retry is a decorator that implements retries with exponential backoff
type Retry struct {
	provider   ports.Provider
	maxRetries int
	baseDelay  time.Duration
}

func NewRetry(provider ports.Provider, maxRetries int, baseDelay time.Duration) *Retry {
	return &Retry{
		provider:   provider,
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
	}
}

func (r *Retry) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
	var err error
	var quote *domain.Quote

	for i := 0; i <= r.maxRetries; i++ {
		if i > 0 {
			// Exponential backoff: baseDelay * 2^(i-1)
			delay := r.baseDelay * time.Duration(math.Pow(2, float64(i-1)))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		quote, err = r.provider.GetQuote(ctx, symbol)
		if err == nil {
			return quote, nil
		}
	}
	return nil, fmt.Errorf("after %d retries: %w", r.maxRetries, err)
}

// RateLimit is a decorator that implements simple token bucket rate limiting
type RateLimit struct {
	provider ports.Provider
	tokens   chan struct{}
}

func NewRateLimit(provider ports.Provider, rps int) *RateLimit {
	rl := &RateLimit{
		provider: provider,
		tokens:   make(chan struct{}, rps),
	}
	// Fill bucket initially
	for i := 0; i < rps; i++ {
		rl.tokens <- struct{}{}
	}

	// Leak tokens (refill)
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(rps))
		defer ticker.Stop()
		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
				// Bucket full
			}
		}
	}()

	return rl
}

func (rl *RateLimit) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-rl.tokens:
		return rl.provider.GetQuote(ctx, symbol)
	}
}
