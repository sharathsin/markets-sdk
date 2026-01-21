# 📈 Markets SDK

![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)

A modular, extensible Go SDK for fetching real-time financial data for **Stocks** and **Cryptocurrencies**. Designed for simplicity, clean architecture, and ease of use.

## 🚀 Features

- **Unified Interface**: Fetch price quotes for any asset type using a single `Provider` interface.
- **Crypto Support**: Built-in support for **CoinGecko** API.
- **Stocks Support**: Built-in support for **Yahoo Finance** API.
- **CLI Tool**: Includes a sleek command-line interface for quick lookups.
- **Zero Heavy Dependencies**: Built primarily with the Go standard library.

## 📦 Installation

To usage the library in your project:

```bash
go get github.com/your-username/markets-sdk
```

## 🛠 Usage

### Library

```go
package main

import (
	"context"
	"fmt"
	"markets-sdk"
	"markets-sdk/providers/coingecko"
)

func main() {
	client := markets.NewMarketClient()
	client.RegisterProvider("crypto", coingecko.NewProvider())

	quote, _ := client.GetQuote(context.Background(), "crypto", "bitcoin")
	fmt.Printf("BTC: $%.2f\n", quote.Price)
}
```

### CLI Tool

Build the tool:
```bash
make build
```

Run it:
```bash
./bin/markets -provider crypto -symbol ethereum
```
*Output:*
```text
------------------------------
ETHEREUM (COINGECKO)
------------------------------
Price:      $2650.00
Change 24h: $120.50
Updated:    3:45PM
------------------------------
```

## 🏗 Architecture

The project follows a clean, interface-driven design:

- **`Provider` Interface**: The contract that all data sources must implement.
- **`MarketClient`**: The high-level orchestrator that manages providers.
- **`providers/`**: detailed implementations for specific APIs (CoinGecko, Yahoo, etc.).

## 🧪 Testing

Run the test suite:

```bash
make test
```

## 📄 License

MIT
