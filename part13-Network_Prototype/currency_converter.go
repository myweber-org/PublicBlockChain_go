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
	key := from + "->" + to
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

	key := from + "->" + to
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
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
	rates map[Currency]float64
}

func NewExchangeRates() *ExchangeRates {
	return &ExchangeRates{
		rates: map[Currency]float64{
			USD: 1.0,
			EUR: 0.85,
			GBP: 0.73,
			JPY: 110.0,
		},
	}
}

func (er *ExchangeRates) Convert(amount float64, from, to Currency) (float64, error) {
	fromRate, fromOk := er.rates[from]
	toRate, toOk := er.rates[to]

	if !fromOk || !toOk {
		return 0, fmt.Errorf("unsupported currency")
	}

	if amount < 0 {
		return 0, fmt.Errorf("amount cannot be negative")
	}

	converted := (amount / fromRate) * toRate
	return math.Round(converted*100) / 100, nil
}

func (er *ExchangeRates) AddRate(currency Currency, rate float64) {
	if rate > 0 {
		er.rates[currency] = rate
	}
}

func main() {
	rates := NewExchangeRates()

	amount := 100.0
	result, err := rates.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)

	rates.AddRate("CAD", 1.25)
	cadResult, _ := rates.Convert(50.0, USD, "CAD")
	fmt.Printf("50.00 %s = %.2f CAD\n", USD, cadResult)
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