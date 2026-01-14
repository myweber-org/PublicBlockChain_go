package main

import (
	"fmt"
	"time"
)

type ExchangeRate struct {
	BaseCurrency    string
	TargetCurrency  string
	Rate            float64
	LastUpdated     time.Time
}

type CurrencyConverter struct {
	rates map[string]ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]ExchangeRate),
	}
}

func (c *CurrencyConverter) AddRate(base, target string, rate float64) {
	key := base + "_" + target
	c.rates[key] = ExchangeRate{
		BaseCurrency:   base,
		TargetCurrency: target,
		Rate:           rate,
		LastUpdated:    time.Now(),
	}
}

func (c *CurrencyConverter) Convert(amount float64, base, target string) (float64, error) {
	if base == target {
		return amount, nil
	}

	key := base + "_" + target
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", base, target)
	}

	return amount * rate.Rate, nil
}

func (c *CurrencyConverter) GetSupportedPairs() []string {
	var pairs []string
	for key := range c.rates {
		pairs = append(pairs, key)
	}
	return pairs
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("EUR", "USD", 1.18)
	converter.AddRate("USD", "JPY", 110.5)
	
	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	fmt.Printf("Supported pairs: %v\n", converter.GetSupportedPairs())
}package main

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
		apiURL:   fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/USD"),
		rates:    make(map[string]float64),
		cacheTTL: 30 * time.Minute,
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastUpdate) < c.cacheTTL && len(c.rates) > 0 {
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

	fromRate, ok1 := c.rates[fromCurrency]
	toRate, ok2 := c.rates[toCurrency]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("invalid currency code: %s or %s", fromCurrency, toCurrency)
	}

	usdAmount := amount / fromRate
	convertedAmount := usdAmount * toRate

	return convertedAmount, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() ([]string, error) {
	if err := c.fetchRates(); err != nil {
		return nil, err
	}

	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func main() {
	converter := NewCurrencyConverter("")

	currencies, err := converter.GetSupportedCurrencies()
	if err != nil {
		fmt.Printf("Error getting currencies: %v\n", err)
		return
	}
	fmt.Printf("Supported currencies: %v\n", currencies[:5])

	amount := 100.0
	from := "USD"
	to := "EUR"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
}