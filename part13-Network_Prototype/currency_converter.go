package main

import (
	"fmt"
	"os"
	"strconv"
)

type ExchangeRate struct {
	Currency string
	Rate     float64
}

var rates = []ExchangeRate{
	{"USD", 1.0},
	{"EUR", 0.85},
	{"GBP", 0.73},
	{"JPY", 110.0},
	{"CAD", 1.25},
}

func convertCurrency(amount float64, fromCurrency, toCurrency string) (float64, error) {
	var fromRate, toRate float64
	foundFrom, foundTo := false, false

	for _, rate := range rates {
		if rate.Currency == fromCurrency {
			fromRate = rate.Rate
			foundFrom = true
		}
		if rate.Currency == toCurrency {
			toRate = rate.Rate
			foundTo = true
		}
	}

	if !foundFrom {
		return 0, fmt.Errorf("unsupported source currency: %s", fromCurrency)
	}
	if !foundTo {
		return 0, fmt.Errorf("unsupported target currency: %s", toCurrency)
	}

	return amount * (toRate / fromRate), nil
}

func listSupportedCurrencies() {
	fmt.Println("Supported currencies:")
	for _, rate := range rates {
		fmt.Printf("  %s (rate: %.4f)\n", rate.Currency, rate.Rate)
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: currency_converter <amount> <from_currency> <to_currency>")
		fmt.Println("Example: currency_converter 100 USD EUR")
		listSupportedCurrencies()
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		os.Exit(1)
	}

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	result, err := convertCurrency(amount, fromCurrency, toCurrency)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}package main

import (
	"fmt"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
}

type CurrencyConverter struct {
	rates []ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: []ExchangeRate{
			{"USD", "EUR", 0.92},
			{"EUR", "USD", 1.09},
			{"USD", "GBP", 0.79},
			{"GBP", "USD", 1.27},
			{"USD", "JPY", 155.0},
			{"JPY", "USD", 0.0065},
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	for _, rate := range c.rates {
		if rate.FromCurrency == fromCurrency && rate.ToCurrency == toCurrency {
			return amount * rate.Rate, nil
		}
	}

	return 0, fmt.Errorf("conversion rate not found for %s to %s", fromCurrency, toCurrency)
}

func (c *CurrencyConverter) AddRate(fromCurrency, toCurrency string, rate float64) {
	c.rates = append(c.rates, ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
	})
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)

	converter.AddRate("EUR", "GBP", 0.86)
	result2, err := converter.Convert(50.0, "EUR", "GBP")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f EUR = %.2f GBP\n", 50.0, result2)
}