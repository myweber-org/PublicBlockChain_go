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
}package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

func fetchExchangeRates(apiKey string) (*ExchangeRates, error) {
	url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/USD")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rates ExchangeRates
	err = json.Unmarshal(body, &rates)
	if err != nil {
		return nil, err
	}

	return &rates, nil
}

func convertCurrency(amount float64, fromCurrency, toCurrency string, rates *ExchangeRates) (float64, error) {
	fromRate, ok1 := rates.Rates[fromCurrency]
	toRate, ok2 := rates.Rates[toCurrency]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("invalid currency code")
	}

	amountInUSD := amount / fromRate
	convertedAmount := amountInUSD * toRate
	return convertedAmount, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: currency_converter <amount> <from_currency> <to_currency>")
		fmt.Println("Example: currency_converter 100 USD EUR")
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		os.Exit(1)
	}

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	rates, err := fetchExchangeRates("")
	if err != nil {
		fmt.Printf("Failed to fetch exchange rates: %v\n", err)
		os.Exit(1)
	}

	result, err := convertCurrency(amount, fromCurrency, toCurrency, rates)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s (as of %s)\n", amount, fromCurrency, result, toCurrency, rates.Date)
}
package main

import (
	"fmt"
	"math"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
)

type ExchangeRates struct {
	rates map[Currency]map[Currency]float64
}

func NewExchangeRates() *ExchangeRates {
	er := &ExchangeRates{
		rates: make(map[Currency]map[Currency]float64),
	}

	baseRates := map[Currency]float64{
		USD: 1.0,
		EUR: 0.85,
		GBP: 0.73,
		JPY: 110.0,
	}

	for from, fromRate := range baseRates {
		er.rates[from] = make(map[Currency]float64)
		for to, toRate := range baseRates {
			er.rates[from][to] = toRate / fromRate
		}
	}

	return er
}

func (er *ExchangeRates) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	rate, exists := er.rates[from][to]
	if !exists {
		return 0, fmt.Errorf("conversion rate not available from %s to %s", from, to)
	}

	converted := amount * rate
	return math.Round(converted*100) / 100, nil
}

func (er *ExchangeRates) UpdateRate(from, to Currency, rate float64) {
	if _, exists := er.rates[from]; !exists {
		er.rates[from] = make(map[Currency]float64)
	}
	er.rates[from][to] = rate

	reciprocalRate := 1.0 / rate
	if _, exists := er.rates[to]; !exists {
		er.rates[to] = make(map[Currency]float64)
	}
	er.rates[to][from] = reciprocalRate
}

func main() {
	rates := NewExchangeRates()

	amount := 100.0
	converted, err := rates.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, converted, EUR)

	rates.UpdateRate(USD, CAD, 1.25)
	converted, err = rates.Convert(amount, USD, CAD)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, converted, CAD)
}