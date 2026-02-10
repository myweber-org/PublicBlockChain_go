
package main

import (
	"fmt"
)

const usdToEurRate = 0.92

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	amount := 100.0
	converted := ConvertUSDToEUR(amount)
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
}