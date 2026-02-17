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
	rates map[string]map[string]float64
	mu    sync.RWMutex
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]map[string]float64),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.rates[from]; !exists {
		c.rates[from] = make(map[string]float64)
	}
	c.rates[from][to] = rate

	// Add inverse rate
	if _, exists := c.rates[to]; !exists {
		c.rates[to] = make(map[string]float64)
	}
	c.rates[to][from] = 1 / rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if from == to {
		return amount, nil
	}

	if _, exists := c.rates[from]; !exists {
		return 0, fmt.Errorf("no rates found for currency: %s", from)
	}

	rate, exists := c.rates[from][to]
	if !exists {
		return 0, fmt.Errorf("no conversion rate from %s to %s", from, to)
	}

	return amount * rate, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter()

	// Add some sample exchange rates
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("USD", "GBP", 0.73)
	converter.AddRate("EUR", "JPY", 130.0)

	// Perform conversions
	amount := 100.0

	usdToEur, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("%.2f USD = %.2f EUR\n", amount, usdToEur)
	}

	eurToJpy, err := converter.Convert(amount, "EUR", "JPY")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("%.2f EUR = %.2f JPY\n", amount, eurToJpy)
	}

	// Try circular conversion
	usdToGbp, err := converter.Convert(amount, "USD", "GBP")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("%.2f USD = %.2f GBP\n", amount, usdToGbp)
	}

	// List supported currencies
	fmt.Println("\nSupported currencies:", converter.GetSupportedCurrencies())
}package main

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
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
	fmt.Printf("Supported pairs: %v\n", converter.GetSupportedPairs())
}