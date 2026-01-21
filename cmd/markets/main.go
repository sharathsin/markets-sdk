package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"markets-sdk"
	"markets-sdk/pkg/providers/coingecko"
	"markets-sdk/pkg/providers/yahoo"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

func main() {
	// Defines flags
	providerFlag := flag.String("provider", "", "Provider to use: 'crypto' or 'stock'")
	symbolFlag := flag.String("symbol", "", "Symbol to fetch (e.g., 'bitcoin', 'AAPL')")
	flag.Parse()

	if *providerFlag == "" || *symbolFlag == "" {
		printUsage()
		os.Exit(1)
	}

	// Initialize Client
	client := markets.NewMarketClient()
	client.RegisterProvider("crypto", coingecko.NewProvider())
	client.RegisterProvider("stock", yahoo.NewProvider())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("%sFetching data for %s...%s\n", ColorCyan, *symbolFlag, ColorReset)

	quote, err := client.GetQuote(ctx, *providerFlag, *symbolFlag)
	if err != nil {
		fmt.Printf("%sError: %v%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}

	printStylish(quote)
}

func printUsage() {
	fmt.Printf("%sMarkets CLI%s\n", ColorBold, ColorReset)
	fmt.Println("Usage:")
	fmt.Println("  markets -provider <crypto|stock> -symbol <name>")
	fmt.Println("\nExamples:")
	fmt.Println("  markets -provider crypto -symbol bitcoin")
	fmt.Println("  markets -provider stock -symbol AAPL")
}

func printStylish(q *markets.Quote) {
	fmt.Println("\n" + strings.Repeat("-", 30))
	fmt.Printf("%s%s (%s)%s\n", ColorBold, strings.ToUpper(q.Symbol), strings.ToUpper(q.Source), ColorReset)
	fmt.Println(strings.Repeat("-", 30))

	fmt.Printf("Price:      %s$%.2f%s\n", ColorBlue, q.Price, ColorReset)

	// Colorize change
	changeColor := ColorGreen
	if q.Change24h < 0 {
		changeColor = ColorRed
	}
	fmt.Printf("Change 24h: %s$%.2f%s\n", changeColor, q.Change24h, ColorReset)

	fmt.Printf("Updated:    %s\n", q.LastUpdated.Format(time.Kitchen))
	fmt.Println(strings.Repeat("-", 30))
}
