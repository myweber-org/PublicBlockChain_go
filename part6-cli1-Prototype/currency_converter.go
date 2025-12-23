package main

import (
	"fmt"
)

func ConvertUSDToEUR(amount float64) float64 {
	const exchangeRate = 0.85
	return amount * exchangeRate
}

func main() {
	usdAmount := 100.0
	eurAmount := ConvertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD is equivalent to %.2f EUR\n", usdAmount, eurAmount)
}