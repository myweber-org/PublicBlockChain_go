
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRates struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	rates     map[string]float64
	lastFetch time.Time
	apiKey    string
}

func NewCurrencyConverter(apiKey string) *CurrencyConverter {
	return &CurrencyConverter{
		rates:  make(map[string]float64),
		apiKey: apiKey,
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastFetch) < 30*time.Minute && len(c.rates) > 0 {
		return nil
	}

	url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/USD")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var exchangeRates ExchangeRates
	err = json.Unmarshal(body, &exchangeRates)
	if err != nil {
		return err
	}

	c.rates = exchangeRates.Rates
	c.rates["USD"] = 1.0
	c.lastFetch = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	err := c.fetchRates()
	if err != nil {
		return 0, err
	}

	fromRate, fromExists := c.rates[from]
	toRate, toExists := c.rates[to]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("invalid currency code")
	}

	usdAmount := amount / fromRate
	convertedAmount := usdAmount * toRate

	return convertedAmount, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	err := c.fetchRates()
	if err != nil {
		return []string{}
	}

	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter("")

	amount := 100.0
	from := "USD"
	to := "EUR"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)

	fmt.Println("Supported currencies:")
	currencies := converter.GetSupportedCurrencies()
	for i, currency := range currencies {
		if i < 10 {
			fmt.Printf("%s ", currency)
		}
	}
	fmt.Println()
}package main

import (
	"fmt"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
}

type CurrencyConverter struct {
	rates []ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: []ExchangeRate{
			{"USD", "EUR", 0.85},
			{"USD", "GBP", 0.73},
			{"USD", "JPY", 110.5},
			{"EUR", "USD", 1.18},
			{"GBP", "USD", 1.37},
			{"JPY", "USD", 0.00905},
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	for _, rate := range c.rates {
		if rate.FromCurrency == fromCurrency && rate.ToCurrency == toCurrency {
			return amount * rate.Rate, nil
		}
	}

	return 0, fmt.Errorf("conversion rate not found for %s to %s", fromCurrency, toCurrency)
}

func (c *CurrencyConverter) AddRate(fromCurrency, toCurrency string, rate float64) {
	c.rates = append(c.rates, ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
	})
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)

	converter.AddRate("EUR", "GBP", 0.86)
	result2, err := converter.Convert(50.0, "EUR", "GBP")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f EUR = %.2f GBP\n", 50.0, result2)
}