
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
	rates := map[Currency]map[Currency]float64{
		USD: {EUR: 0.92, GBP: 0.79, JPY: 149.50},
		EUR: {USD: 1.09, GBP: 0.86, JPY: 162.50},
		GBP: {USD: 1.27, EUR: 1.16, JPY: 189.20},
		JPY: {USD: 0.0067, EUR: 0.0062, GBP: 0.0053},
	}
	return &ExchangeRates{rates: rates}
}

func (er *ExchangeRates) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	rateMap, exists := er.rates[from]
	if !exists {
		return 0, fmt.Errorf("unsupported source currency: %s", from)
	}

	rate, exists := rateMap[to]
	if !exists {
		return 0, fmt.Errorf("conversion from %s to %s not supported", from, to)
	}

	converted := amount * rate
	return math.Round(converted*100) / 100, nil
}

func (er *ExchangeRates) AddRate(from, to Currency, rate float64) {
	if er.rates[from] == nil {
		er.rates[from] = make(map[Currency]float64)
	}
	er.rates[from][to] = rate
}

func main() {
	converter := NewExchangeRates()

	amount := 100.0
	result, err := converter.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)

	converter.AddRate(USD, CAD, 1.35)
	cadResult, _ := converter.Convert(amount, USD, CAD)
	fmt.Printf("%.2f %s = %.2f CAD\n", amount, USD, cadResult)
}