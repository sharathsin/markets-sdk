package ports

import (
	"context"
	"markets-sdk/pkg/domain"
)

// Provider defines the interface for fetching market data
type Provider interface {
	// GetQuote returns the latest quote for a given symbol
	GetQuote(ctx context.Context, symbol string) (*domain.Quote, error)
}
