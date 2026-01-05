package main

import (
	"fmt"
	"sync"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
}

type CurrencyConverter struct {
	rates map[string]float64
	mu    sync.RWMutex
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]float64),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := fmt.Sprintf("%s:%s", from, to)
	c.rates[key] = rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", from, to)
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
	}

	return amount * rate, nil
}

func (c *CurrencyConverter) GetSupportedPairs() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	pairs := make([]string, 0, len(c.rates))
	for key := range c.rates {
		pairs = append(pairs, key)
	}
	return pairs
}

func main() {
	converter := NewCurrencyConverter()

	converter.AddRate("USD", "EUR", 0.92)
	converter.AddRate("EUR", "USD", 1.09)
	converter.AddRate("USD", "JPY", 148.50)

	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
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
	apiEndpoint string
	client      *http.Client
	cache       map[string]float64
	lastUpdated time.Time
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		apiEndpoint: "https://api.exchangerate-api.com/v4/latest/USD",
		client:      &http.Client{Timeout: 10 * time.Second},
		cache:       make(map[string]float64),
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastUpdated) < 30*time.Minute && len(c.cache) > 0 {
		return nil
	}

	resp, err := c.client.Get(c.apiEndpoint)
	if err != nil {
		return fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var rateResponse ExchangeRateResponse
	if err := json.Unmarshal(body, &rateResponse); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	c.cache = rateResponse.Rates
	c.lastUpdated = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if err := c.fetchRates(); err != nil {
		return 0, err
	}

	fromRate, fromExists := c.cache[fromCurrency]
	toRate, toExists := c.cache[toCurrency]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("unsupported currency pair: %s to %s", fromCurrency, toCurrency)
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
    fromRate, okFrom := rates.Rates[fromCurrency]
    toRate, okTo := rates.Rates[toCurrency]

    if !okFrom || !okTo {
        return 0, fmt.Errorf("invalid currency code")
    }

    amountInUSD := amount / fromRate
    convertedAmount := amountInUSD * toRate
    return convertedAmount, nil
}

func main() {
    if len(os.Args) != 4 {
        fmt.Println("Usage: go run currency_converter.go <amount> <from_currency> <to_currency>")
        fmt.Println("Example: go run currency_converter.go 100 USD EUR")
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
        fmt.Printf("Error fetching exchange rates: %v\n", err)
        os.Exit(1)
    }

    converted, err := convertCurrency(amount, fromCurrency, toCurrency, rates)
    if err != nil {
        fmt.Printf("Conversion error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, converted, toCurrency)
}