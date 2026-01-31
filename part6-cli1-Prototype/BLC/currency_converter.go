package main

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
	ttl       time.Duration
}

func NewCurrencyConverter(ttl time.Duration) *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]float64),
		ttl:   ttl,
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastFetch) < c.ttl && len(c.rates) > 0 {
		return nil
	}

	resp, err := http.Get("https://api.exchangerate-api.com/v4/latest/USD")
	if err != nil {
		return fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var exchangeRates ExchangeRates
	if err := json.Unmarshal(body, &exchangeRates); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	c.rates = exchangeRates.Rates
	c.lastFetch = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if err := c.fetchRates(); err != nil {
		return 0, err
	}

	if from == to {
		return amount, nil
	}

	fromRate, fromExists := c.rates[from]
	toRate, toExists := c.rates[to]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("unsupported currency: %s or %s", from, to)
	}

	if from == "USD" {
		return amount * toRate, nil
	}

	if to == "USD" {
		return amount / fromRate, nil
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
}

func main() {
	converter := NewCurrencyConverter(30 * time.Minute)

	amount := 100.0
	from := "EUR"
	to := "JPY"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
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
			{"USD", "JPY", 147.50},
			{"JPY", "USD", 0.0068},
			{"GBP", "USD", 1.27},
			{"USD", "GBP", 0.79},
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
	from := "USD"
	to := "EUR"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)

	converter.AddRate("EUR", "JPY", 160.32)
	eurToJpy, _ := converter.Convert(50.0, "EUR", "JPY")
	fmt.Printf("50.00 EUR = %.2f JPY\n", eurToJpy)
}
package main

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
	apiKey    string
}

func NewCurrencyConverter(apiKey string) *CurrencyConverter {
	return &CurrencyConverter{
		rates:  make(map[string]float64),
		apiKey: apiKey,
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastFetch) < 30*time.Minute && len(c.rates) > 0 {
		return nil
	}

	url := fmt.Sprintf("https://api.exchangerate.host/latest?base=USD&access_key=%s", c.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var exchangeRates ExchangeRates
	if err := json.Unmarshal(body, &exchangeRates); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	c.rates = exchangeRates.Rates
	c.rates["USD"] = 1.0
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
		return 0, fmt.Errorf("invalid currency code: %s or %s", from, to)
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
}

func main() {
	converter := NewCurrencyConverter("your_api_key_here")

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