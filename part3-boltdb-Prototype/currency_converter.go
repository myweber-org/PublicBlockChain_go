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