package markets

import (
	"context"
	"fmt"
	"markets-sdk/pkg/domain"
	"markets-sdk/pkg/ports"
)

// Ensure backward compatibility and ease of use
type Quote = domain.Quote
type AssetType = domain.AssetType

// MarketClient is the main entry point that can manage multiple providers
type MarketClient struct {
	providers map[string]ports.Provider
}

// NewMarketClient creates a new client
func NewMarketClient() *MarketClient {
	return &MarketClient{
		providers: make(map[string]ports.Provider),
	}
}

// RegisterProvider registers a provider with a specific name
func (c *MarketClient) RegisterProvider(name string, p ports.Provider) {
	c.providers[name] = p
}

// GetQuote fetches a quote from a specific provider
func (c *MarketClient) GetQuote(ctx context.Context, providerName string, symbol string) (*domain.Quote, error) {
	p, ok := c.providers[providerName]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", providerName)
	}
	return p.GetQuote(ctx, symbol)
}
