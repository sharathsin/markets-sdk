# Markets SDK v2.0.0 Release Notes

We are thrilled to announce the release of **Markets SDK v2.0.0**. This major release transforms the library into a production-grade, enterprise-ready solution with a focus on stability, observability, and performance.

## 🌟 Key Highlights

### 🏗 Hexagonal Architecture
The codebase has been completely refactored into a **Clean/Hexagonal Architecture**.
- **Domain-Driven**: Core logic is isolated in `pkg/domain`, free from external dependencies.
- **Pluggable Providers**: Integration with CoinGecko and Yahoo Finance is now handled via flexible adapters in `pkg/providers`.
- **Extensible**: New data sources can be added by simply implementing the `ports.Provider` interface.

### 🛡 Resilience & Stability
We've introduced a suite of **Resilience Decorators** to ensuring your application stays up even when external APIs go down:
- **Circuit Breaker**: Automatically stops requests to failing providers to prevent cascading failures.
- **Smart Retries**: Implements exponential backoff for transient network errors.
- **Rate Limiting**: Built-in protection to respect API rate limits.

### 👁 Observability
The SDK is now opaque no more. Default hooks for:
- **Structured Logging**: Built on Go's `log/slog`.
- **Metrics**: Prometheus-ready interfaces for request counts and latency.
- **Tracing**: OpenTelemetry-compatible spans for request tracing.

### ⚡️ Performance
- **Zero-Allocation Decoding**: Implemented `sync.Pool` for JSON response processing, resulting in a **~10% performance boost** and reduced GC pressure.
- **Concurrency Best Practices**: Strict `context.Context` propagation throughout the call chain.

### 🛠 Developer Experience
- **CI/CD**: New GitHub Actions pipeline for automated linting, testing, and building.
- **Comprehensive Testing**: high test coverage across core clients and decorators.
- **Design Docs**: Added `DESIGN.md` for architectural transparency.

## 📦 Compatibility
- **Go Version**: Requires Go 1.21+
- **Public API**: Backward compatibility maintained via type aliases in `markets.go`.

## 🚀 Getting Started

```bash
go get github.com/sharathsin/markets-sdk@v2.0.0
```

```go
client := markets.NewMarketClient()
// ... use as before, but with superpowers!
```
