package markets_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"markets-sdk"
	"markets-sdk/pkg/domain"
	// Use mock provider from internal tests manually or define a simple one here
	// Since markets_test is external, we need a local mock
)

type mockProvider struct {
	price float64
	err   error
}

func (m *mockProvider) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &domain.Quote{
		Symbol:      symbol,
		Price:       m.price,
		LastUpdated: time.Now(),
		Source:      "mock",
	}, nil
}

func TestMarketClient(t *testing.T) {
	client := markets.NewMarketClient()

	// 1. Test Register
	mock := &mockProvider{price: 150.0}
	client.RegisterProvider("test-provider", mock)

	// 2. Test GetQuote Success
	q, err := client.GetQuote(context.Background(), "test-provider", "ABC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Price != 150.0 {
		t.Errorf("expected price 150.0, got %f", q.Price)
	}
	if q.Source != "mock" {
		t.Errorf("expected source 'mock', got %s", q.Source)
	}

	// 3. Test Provider Not Found
	_, err = client.GetQuote(context.Background(), "missing-provider", "ABC")
	if err == nil {
		t.Error("expected error for missing provider")
	}

	// 4. Test Provider Error
	failMock := &mockProvider{err: errors.New("network error")}
	client.RegisterProvider("fail-provider", failMock)
	_, err = client.GetQuote(context.Background(), "fail-provider", "ABC")
	if err == nil {
		t.Error("expected error from provider")
	}
}
