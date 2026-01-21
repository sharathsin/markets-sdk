package decorators_test

import (
	"context"
	"errors"
	"markets-sdk/pkg/domain"
)

// MockProvider is a helper for testing decorators
type MockProvider struct {
	QuoteFn func(ctx context.Context, symbol string) (*domain.Quote, error)
}

func (m *MockProvider) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
	if m.QuoteFn != nil {
		return m.QuoteFn(ctx, symbol)
	}
	return nil, errors.New("not implemented")
}
