package main

import (
	"fmt"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
)

var exchangeRates = map[Currency]float64{
	USD: 1.0,
	EUR: 0.85,
	GBP: 0.73,
}

func Convert(amount float64, from Currency, to Currency) (float64, error) {
	fromRate, okFrom := exchangeRates[from]
	toRate, okTo := exchangeRates[to]

	if !okFrom || !okTo {
		return 0, fmt.Errorf("unsupported currency")
	}

	usdAmount := amount / fromRate
	convertedAmount := usdAmount * toRate
	return convertedAmount, nil
}

func main() {
	amount := 100.0
	result, err := Convert(amount, USD, EUR)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)

	result, err = Convert(amount, EUR, GBP)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, EUR, result, GBP)
}
package main

import (
	"fmt"
)

func convertUSDToEUR(amount float64) float64 {
	const exchangeRate = 0.85
	return amount * exchangeRate
}

func main() {
	var usdAmount float64
	fmt.Print("Enter amount in USD: ")
	fmt.Scan(&usdAmount)

	if usdAmount < 0 {
		fmt.Println("Amount cannot be negative")
		return
	}

	eurAmount := convertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}package main

import (
	"fmt"
	"sync"
)

type ExchangeRate struct {
	BaseCurrency    string
	TargetCurrency  string
	Rate            float64
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

func (c *CurrencyConverter) AddRate(base, target string, rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.rates[base] == nil {
		c.rates[base] = make(map[string]float64)
	}
	c.rates[base][target] = rate

	// Add inverse rate
	if c.rates[target] == nil {
		c.rates[target] = make(map[string]float64)
	}
	c.rates[target][base] = 1.0 / rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if from == to {
		return amount, nil
	}

	if targetRates, ok := c.rates[from]; ok {
		if rate, ok := targetRates[to]; ok {
			return amount * rate, nil
		}
	}

	return 0, fmt.Errorf("no conversion rate available from %s to %s", from, to)
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
	amounts := []float64{100.0, 250.0, 500.0}
	conversions := []struct {
		from string
		to   string
	}{
		{"USD", "EUR"},
		{"EUR", "USD"},
		{"USD", "GBP"},
		{"GBP", "JPY"},
	}

	for _, amount := range amounts {
		fmt.Printf("Converting %.2f:\n", amount)
		for _, conv := range conversions {
			result, err := converter.Convert(amount, conv.from, conv.to)
			if err != nil {
				fmt.Printf("  %s to %s: %v\n", conv.from, conv.to, err)
			} else {
				fmt.Printf("  %s to %s: %.2f\n", conv.from, conv.to, result)
			}
		}
		fmt.Println()
	}

	// Show supported currencies
	fmt.Println("Supported currencies:", converter.GetSupportedCurrencies())
}