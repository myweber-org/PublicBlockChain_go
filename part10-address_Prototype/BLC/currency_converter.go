
package main

import (
	"fmt"
	"os"
)

type ExchangeRate struct {
	Currency string
	Rate     float64
}

type CurrencyConverter struct {
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: map[string]float64{
			"USD": 1.0,
			"EUR": 0.92,
			"GBP": 0.79,
			"JPY": 149.5,
			"CAD": 1.36,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	fromRate, fromExists := c.rates[fromCurrency]
	toRate, toExists := c.rates[toCurrency]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("unsupported currency")
	}

	baseAmount := amount / fromRate
	return baseAmount * toRate, nil
}

func (c *CurrencyConverter) AddRate(currency string, rate float64) {
	c.rates[currency] = rate
}

func main() {
	converter := NewCurrencyConverter()

	if len(os.Args) < 4 {
		fmt.Println("Usage: go run currency_converter.go <amount> <from_currency> <to_currency>")
		fmt.Println("Example: go run currency_converter.go 100 USD EUR")
		return
	}

	var amount float64
	_, err := fmt.Sscanf(os.Args[1], "%f", &amount)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		return
	}

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	result, err := converter.Convert(amount, fromCurrency, toCurrency)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRates struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	rates     map[string]float64
	lastFetch time.Time
	cacheTTL  time.Duration
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates:    make(map[string]float64),
		cacheTTL: 30 * time.Minute,
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastFetch) < c.cacheTTL && len(c.rates) > 0 {
		return nil
	}

	resp, err := http.Get("https://api.exchangerate-api.com/v4/latest/USD")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var rates ExchangeRates
	if err := json.Unmarshal(body, &rates); err != nil {
		return err
	}

	c.rates = rates.Rates
	c.lastFetch = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if err := c.fetchRates(); err != nil {
		return 0, err
	}

	fromRate, fromExists := c.rates[from]
	toRate, toExists := c.rates[to]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("unsupported currency code")
	}

	if from == "USD" {
		return amount * toRate, nil
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	from := "EUR"
	to := "JPY"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
}