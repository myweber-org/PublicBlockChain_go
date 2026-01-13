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
			"JPY": 149.32,
			"CAD": 1.36,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	fromRate, fromExists := c.rates[fromCurrency]
	toRate, toExists := c.rates[toCurrency]

	if !fromExists {
		return 0, fmt.Errorf("unsupported source currency: %s", fromCurrency)
	}
	if !toExists {
		return 0, fmt.Errorf("unsupported target currency: %s", toCurrency)
	}

	if fromRate == 0 {
		return 0, fmt.Errorf("invalid exchange rate for currency: %s", fromCurrency)
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
}

func (c *CurrencyConverter) AddRate(currency string, rate float64) error {
	if rate <= 0 {
		return fmt.Errorf("exchange rate must be positive")
	}
	c.rates[currency] = rate
	return nil
}

func (c *CurrencyConverter) ListCurrencies() []string {
	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter()

	err := converter.AddRate("AUD", 1.54)
	if err != nil {
		fmt.Printf("Failed to add rate: %v\n", err)
		os.Exit(1)
	}

	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Available currencies: %v\n", converter.ListCurrencies())
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)

	result, err = converter.Convert(50.0, "GBP", "JPY")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("50.00 GBP = %.2f JPY\n", result)
}package main

import (
	"fmt"
)

const usdToEurRate = 0.92

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	usdAmount := 100.0
	eurAmount := ConvertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}