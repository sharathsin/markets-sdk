package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"markets-sdk"
)

const baseURL = "https://query2.finance.yahoo.com/v8/finance/chart"

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

// chartResponse matches the structure of Yahoo Finance chart API
type chartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol             string  `json:"symbol"`
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				RegularMarketTime  int64   `json:"regularMarketTime"`
				PreviousClose      float64 `json:"chartPreviousClose"`
			} `json:"meta"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

func (p *Provider) GetQuote(ctx context.Context, symbol string) (*markets.Quote, error) {
	url := fmt.Sprintf("%s/%s?interval=1m&range=1d", baseURL, symbol)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Yahoo Finance often requires a User-Agent to not block the request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data chartResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if data.Chart.Error != nil {
		return nil, fmt.Errorf("yahoo api error: %s - %s", data.Chart.Error.Code, data.Chart.Error.Description)
	}

	if len(data.Chart.Result) == 0 {
		return nil, fmt.Errorf("symbol %s not found", symbol)
	}

	meta := data.Chart.Result[0].Meta

	// Calculate simple change
	change := meta.RegularMarketPrice - meta.PreviousClose

	// Note: Volume is in the 'indicators' part of the JSON which is more complex to parse for just a quote.
	// For this MVP, we will omit volume or could parse it if strictly needed.
	// Let's stick to price and change for now.

	return &markets.Quote{
		Symbol:      meta.Symbol,
		Price:       meta.RegularMarketPrice,
		Change24h:   change,
		Volume:      0, // Not easily available in Meta, needs parsing Quote array
		LastUpdated: time.Unix(meta.RegularMarketTime, 0),
		Source:      "yahoo",
	}, nil
}
