package main

import (
	"fmt"
	"time"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
	LastUpdated  time.Time
}

type CurrencyConverter struct {
	rates map[string]ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]ExchangeRate),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	key := from + "_" + to
	c.rates[key] = ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		LastUpdated:  time.Now(),
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}

	key := from + "_" + to
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
	}

	return amount * rate.Rate, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	currencies := make(map[string]bool)
	for _, rate := range c.rates {
		currencies[rate.FromCurrency] = true
		currencies[rate.ToCurrency] = true
	}

	result := make([]string, 0, len(currencies))
	for currency := range currencies {
		result = append(result, currency)
	}
	return result
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("EUR", "USD", 1.18)
	converter.AddRate("USD", "JPY", 110.5)
	converter.AddRate("GBP", "USD", 1.38)
	
	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	fmt.Printf("Supported currencies: %v\n", converter.GetSupportedCurrencies())
}