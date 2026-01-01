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
			"EUR": 0.85,
			"GBP": 0.73,
			"JPY": 110.0,
			"CAD": 1.25,
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

	if amount < 0 {
		return 0, fmt.Errorf("amount cannot be negative")
	}

	usdAmount := amount / fromRate
	convertedAmount := usdAmount * toRate

	return convertedAmount, nil
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

	result, err := converter.Convert(100.0, "USD", "EUR")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("100 USD = %.2f EUR\n", result)
	fmt.Printf("Available currencies: %v\n", converter.ListCurrencies())
}