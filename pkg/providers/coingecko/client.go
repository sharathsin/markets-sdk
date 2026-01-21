package coingecko

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"markets-sdk/pkg/domain"
)

const baseURL = "https://api.coingecko.com/api/v3"

// Buffer pool to reduce GC pressure when reading response bodies
var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

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
type simplePriceResponse map[string]struct {
	USD           float64 `json:"usd"`
	USD24hChange  float64 `json:"usd_24h_change"`
	USD24hVol     float64 `json:"usd_24h_vol"`
	LastUpdatedAt float64 `json:"last_updated_at"`
}

func (p *Provider) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
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

	// Use pooled buffer to read body. This allows us to have the body ensuring
	// we can log it if needed, or re-process, while minimizing allocations.
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	if _, err := io.Copy(buf, resp.Body); err != nil {
		return nil, err
	}

	var data simplePriceResponse
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w. body: %s", err, buf.String())
	}

	item, ok := data[id]
	if !ok {
		return nil, fmt.Errorf("symbol %s not found in response", symbol)
	}

	return &domain.Quote{
		Symbol:      symbol,
		Price:       item.USD,
		Change24h:   item.USD24hChange,
		Volume:      item.USD24hVol,
		LastUpdated: time.Now(),
		Source:      "coingecko",
	}, nil
}
