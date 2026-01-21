package markets

import (
    "context"
    "fmt"
    "time"
)

// AssetType represents the type of financial asset
type AssetType string

const (
    AssetTypeStock  AssetType = "STOCK"
    AssetTypeCrypto AssetType = "CRYPTO"
)

// Quote represents a simplified pricing quote for an asset
type Quote struct {
    Symbol      string    `json:"symbol"`
    Price       float64   `json:"price"`
    Change24h   float64   `json:"change_24h,omitempty"`
    Volume      float64   `json:"volume,omitempty"`
    LastUpdated time.Time `json:"last_updated"`
    Source      string    `json:"source"`
}

// Provider defines the interface for fetching market data
type Provider interface {
    // GetQuote returns the latest quote for a given symbol
    GetQuote(ctx context.Context, symbol string) (*Quote, error)
}

// MarketClient is the main entry point that can manage multiple providers
type MarketClient struct {
    providers map[string]Provider
}

// NewMarketClient creates a new client
func NewMarketClient() *MarketClient {
    return &MarketClient{
        providers: make(map[string]Provider),
    }
}

// RegisterProvider registers a provider with a specific name
func (c *MarketClient) RegisterProvider(name string, p Provider) {
    c.providers[name] = p
}

// GetQuote fetches a quote from a specific provider
func (c *MarketClient) GetQuote(ctx context.Context, providerName string, symbol string) (*Quote, error) {
    p, ok := c.providers[providerName]
    if !ok {
        return nil, fmt.Errorf("provider %s not found", providerName)
    }
    return p.GetQuote(ctx, symbol)
}
