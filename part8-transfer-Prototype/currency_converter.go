package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	apiURL     string
	rates      map[string]float64
	lastUpdate time.Time
	cacheTTL   time.Duration
}

func NewCurrencyConverter(apiKey string) *CurrencyConverter {
	return &CurrencyConverter{
		apiURL:     fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/USD"),
		rates:      make(map[string]float64),
		cacheTTL:   30 * time.Minute,
		lastUpdate: time.Now().Add(-1 * time.Hour),
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastUpdate) < c.cacheTTL {
		return nil
	}

	resp, err := http.Get(c.apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var data ExchangeRateResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	c.rates = data.Rates
	c.lastUpdate = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if err := c.fetchRates(); err != nil {
		return 0, err
	}

	fromRate, fromExists := c.rates[fromCurrency]
	toRate, toExists := c.rates[toCurrency]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("unsupported currency: %s or %s", fromCurrency, toCurrency)
	}

	if fromCurrency == "USD" {
		return amount * toRate, nil
	}

	amountInUSD := amount / fromRate
	return amountInUSD * toRate, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	if err := c.fetchRates(); err != nil {
		return []string{}
	}

	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter("")

	amount := 100.0
	from := "EUR"
	to := "JPY"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
	
	fmt.Println("Supported currencies:")
	for _, currency := range converter.GetSupportedCurrencies() {
		fmt.Printf("- %s\n", currency)
	}
}