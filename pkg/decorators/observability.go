package decorators

import (
	"context"
	"log/slog"
	"time"

	"markets-sdk/pkg/domain"
	"markets-sdk/pkg/ports"
)

// LoggingDecorator logs request details using slog
type LoggingDecorator struct {
	provider     ports.Provider
	logger       *slog.Logger
	providerName string
}

func NewLoggingDecorator(provider ports.Provider, logger *slog.Logger, providerName string) *LoggingDecorator {
	return &LoggingDecorator{
		provider:     provider,
		logger:       logger,
		providerName: providerName,
	}
}

func (l *LoggingDecorator) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
	start := time.Now()

	l.logger.Info("fetching quote", "provider", l.providerName, "symbol", symbol)

	quote, err := l.provider.GetQuote(ctx, symbol)

	duration := time.Since(start)

	if err != nil {
		l.logger.Error("failed to fetch quote",
			"provider", l.providerName,
			"symbol", symbol,
			"error", err,
			"duration", duration,
		)
		return nil, err
	}

	l.logger.Info("fetched quote",
		"provider", l.providerName,
		"symbol", symbol,
		"price", quote.Price,
		"duration", duration,
	)

	return quote, nil
}

// MetricsCollector defines an interface for collecting metrics, allowing Prometheus or other backends
type MetricsCollector interface {
	IncRequest(provider, status string)
	ObserveDuration(provider string, duration float64)
}

// MetricsDecorator captures metrics for requests
type MetricsDecorator struct {
	provider  ports.Provider
	collector MetricsCollector
	name      string
}

func NewMetricsDecorator(provider ports.Provider, collector MetricsCollector, name string) *MetricsDecorator {
	return &MetricsDecorator{
		provider:  provider,
		collector: collector,
		name:      name,
	}
}

func (m *MetricsDecorator) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
	start := time.Now()
	quote, err := m.provider.GetQuote(ctx, symbol)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	if m.collector != nil {
		m.collector.IncRequest(m.name, status)
		m.collector.ObserveDuration(m.name, duration)
	}

	return quote, err
}

// Tracer defines a minimal interface for distributed tracing (compatible with OpenTelemetry)
type Tracer interface {
	Start(ctx context.Context, spanName string) (context.Context, Span)
}

type Span interface {
	End()
	RecordError(err error)
}

// TracingDecorator wraps a provider with tracing spans
type TracingDecorator struct {
	provider ports.Provider
	tracer   Tracer
	name     string
}

func NewTracingDecorator(provider ports.Provider, tracer Tracer, name string) *TracingDecorator {
	return &TracingDecorator{
		provider: provider,
		tracer:   tracer,
		name:     name,
	}
}

func (t *TracingDecorator) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
	if t.tracer == nil {
		return t.provider.GetQuote(ctx, symbol)
	}

	ctx, span := t.tracer.Start(ctx, "GetQuote/"+t.name)
	defer span.End()

	quote, err := t.provider.GetQuote(ctx, symbol)
	if err != nil {
		span.RecordError(err)
	}
	return quote, err
}
