package main

import (
	"fmt"
	"time"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
	LastUpdated  time.Time
}

type CurrencyConverter struct {
	rates map[string]ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]ExchangeRate),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	key := from + "_" + to
	c.rates[key] = ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		LastUpdated:  time.Now(),
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}

	key := from + "_" + to
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
	}

	return amount * rate.Rate, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	currencies := make(map[string]bool)
	for _, rate := range c.rates {
		currencies[rate.FromCurrency] = true
		currencies[rate.ToCurrency] = true
	}

	result := make([]string, 0, len(currencies))
	for currency := range currencies {
		result = append(result, currency)
	}
	return result
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("EUR", "USD", 1.18)
	converter.AddRate("USD", "JPY", 110.5)
	converter.AddRate("GBP", "USD", 1.38)
	
	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	fmt.Printf("Supported currencies: %v\n", converter.GetSupportedCurrencies())
}package main

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

func convertAmount(amount float64, fromCurrency, toCurrency string) (float64, error) {
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

	result, err := convertAmount(amount, fromCurrency, toCurrency)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}
package main

import (
	"fmt"
	"time"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
	LastUpdated  time.Time
}

type CurrencyConverter struct {
	rates map[string]ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]ExchangeRate),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	key := from + "_" + to
	c.rates[key] = ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		LastUpdated:  time.Now(),
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}

	key := from + "_" + to
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
	}

	return amount * rate.Rate, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	currencies := make(map[string]bool)
	for _, rate := range c.rates {
		currencies[rate.FromCurrency] = true
		currencies[rate.ToCurrency] = true
	}

	result := make([]string, 0, len(currencies))
	for currency := range currencies {
		result = append(result, currency)
	}
	return result
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("EUR", "USD", 1.18)
	converter.AddRate("USD", "JPY", 110.5)
	
	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	fmt.Printf("Supported currencies: %v\n", converter.GetSupportedCurrencies())
}