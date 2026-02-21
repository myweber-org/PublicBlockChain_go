package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
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

	var data ExchangeRateResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	c.rates = data.Rates
	c.lastFetch = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if err := c.fetchRates(); err != nil {
		return 0, err
	}

	fromRate, ok1 := c.rates[from]
	toRate, ok2 := c.rates[to]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("invalid currency code")
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
}

func main() {
	converter := NewCurrencyConverter("")
	
	amount := 100.0
	from := "EUR"
	to := "JPY"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
}