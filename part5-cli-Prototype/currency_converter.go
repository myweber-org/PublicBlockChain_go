
package main

import (
	"fmt"
)

const usdToEurRate = 0.85

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	usdAmount := 100.0
	eurAmount := ConvertUSDToEUR(usdAmount)
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
	
	if _, exists := c.rates[base]; !exists {
		c.rates[base] = make(map[string]float64)
	}
	c.rates[base][target] = rate
	
	if _, exists := c.rates[target]; !exists {
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
	
	if targetRates, exists := c.rates[from]; exists {
		if rate, found := targetRates[to]; found {
			return amount * rate, nil
		}
	}
	
	return 0, fmt.Errorf("conversion rate not available from %s to %s", from, to)
}

func (c *CurrencyConverter) GetAllRates() []ExchangeRate {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	var rates []ExchangeRate
	for base, targetRates := range c.rates {
		for target, rate := range targetRates {
			rates = append(rates, ExchangeRate{
				BaseCurrency:   base,
				TargetCurrency: target,
				Rate:           rate,
			})
		}
	}
	return rates
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("USD", "GBP", 0.73)
	converter.AddRate("EUR", "JPY", 130.0)
	
	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
	
	allRates := converter.GetAllRates()
	fmt.Println("\nAvailable exchange rates:")
	for _, rate := range allRates {
		fmt.Printf("%s -> %s: %.4f\n", rate.BaseCurrency, rate.TargetCurrency, rate.Rate)
	}
}