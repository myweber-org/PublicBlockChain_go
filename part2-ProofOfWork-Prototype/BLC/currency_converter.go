
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

type ExchangeRate struct {
	From Currency
	To   Currency
	Rate float64
}

type CurrencyConverter struct {
	rates []ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: []ExchangeRate{
			{USD, EUR, 0.92},
			{USD, GBP, 0.79},
			{USD, JPY, 149.50},
			{EUR, USD, 1.09},
			{EUR, GBP, 0.86},
			{EUR, JPY, 162.50},
			{GBP, USD, 1.27},
			{GBP, EUR, 1.16},
			{GBP, JPY, 189.24},
			{JPY, USD, 0.0067},
			{JPY, EUR, 0.0062},
			{JPY, GBP, 0.0053},
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from Currency, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	for _, rate := range c.rates {
		if rate.From == from && rate.To == to {
			result := amount * rate.Rate
			return math.Round(result*100) / 100, nil
		}
	}

	return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
}

func (c *CurrencyConverter) AddRate(from Currency, to Currency, rate float64) {
	c.rates = append(c.rates, ExchangeRate{from, to, rate})
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	converted, err := converter.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, converted, EUR)

	converted, err = converter.Convert(amount, EUR, JPY)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, EUR, converted, JPY)

	converter.AddRate(USD, CAD, 1.36)
	converted, err = converter.Convert(amount, USD, CAD)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, converted, CAD)
}