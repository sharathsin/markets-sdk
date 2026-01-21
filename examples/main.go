package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"markets-sdk"
	"markets-sdk/pkg/providers/coingecko"
	"markets-sdk/pkg/providers/yahoo"
)

func main() {
	client := markets.NewMarketClient()

	// Register providers
	client.RegisterProvider("crypto", coingecko.NewProvider())
	client.RegisterProvider("stock", yahoo.NewProvider())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Fetch Crypto
	cryptoSymbol := "bitcoin"
	fmt.Printf("Fetching %s...\n", cryptoSymbol)
	q1, err := client.GetQuote(ctx, "crypto", cryptoSymbol)
	if err != nil {
		log.Printf("Error fetching crypto: %v\n", err)
	} else {
		printQuote(q1)
	}

	fmt.Println("--------------------")

	// 2. Fetch Stock
	stockSymbol := "AAPL"
	fmt.Printf("Fetching %s...\n", stockSymbol)
	q2, err := client.GetQuote(ctx, "stock", stockSymbol)
	if err != nil {
		log.Printf("Error fetching stock: %v\n", err)
	} else {
		printQuote(q2)
	}
}

func printQuote(q *markets.Quote) {
	fmt.Printf("Symbol: %s\n", q.Symbol)
	fmt.Printf("Price: %.2f\n", q.Price)
	fmt.Printf("Change: %.2f\n", q.Change24h)
	fmt.Printf("Source: %s\n", q.Source)
	fmt.Printf("Updated: %s\n", q.LastUpdated.Format(time.RFC3339))
}
