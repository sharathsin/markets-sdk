package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"markets-sdk"
)

const baseURL = "https://api.coingecko.com/api/v3"

type Provider struct {
	client *http.Client
}

func NewProvider() *Provider {
	return &Provider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// simplePriceResponse matches the structure returned by /simple/price
// e.g., {"bitcoin": {"usd": 50000.0, "usd_24h_change": 2.5, "usd_24h_vol": 100000}}
type simplePriceResponse map[string]struct {
	USD           float64 `json:"usd"`
	USD24hChange  float64 `json:"usd_24h_change"`
	USD24hVol     float64 `json:"usd_24h_vol"`
	LastUpdatedAt float64 `json:"last_updated_at"` // sometimes present in other endpoints, but simple/price is simpler
}

func (p *Provider) GetQuote(ctx context.Context, symbol string) (*markets.Quote, error) {
	// CoinGecko uses IDs (e.g., "bitcoin") rather than symbols (e.g., "BTC") for the main query,
	// but often users might pass "bitcoin". For this simple SDK, we'll assume the user passes the ID.
	// In a production app, we'd need a mapping symbol->id.

	id := strings.ToLower(symbol)
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd&include_24hr_vol=true&include_24hr_change=true", baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data simplePriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	item, ok := data[id]
	if !ok {
		return nil, fmt.Errorf("symbol %s not found in response", symbol)
	}

	return &markets.Quote{
		Symbol:      symbol,
		Price:       item.USD,
		Change24h:   item.USD24hChange,
		Volume:      item.USD24hVol,
		LastUpdated: time.Now(), // Simple price doesn't return timestamp, so we use Now()
		Source:      "coingecko",
	}, nil
}
