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
	rates map[string]map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]map[string]float64),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	if c.rates[from] == nil {
		c.rates[from] = make(map[string]float64)
	}
	c.rates[from][to] = rate
	
	if c.rates[to] == nil {
		c.rates[to] = make(map[string]float64)
	}
	c.rates[to][from] = 1.0 / rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}

	if rate, exists := c.rates[from][to]; exists {
		return amount * rate, nil
	}

	return 0, fmt.Errorf("conversion rate not available from %s to %s", from, to)
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("USD", "GBP", 0.73)
	converter.AddRate("EUR", "JPY", 130.0)
	
	amount := 100.0
	
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
	
	currencies := converter.GetSupportedCurrencies()
	fmt.Println("Supported currencies:", currencies)
}package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ExchangeRates struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	client     *http.Client
	apiBaseURL string
	cache      map[string]ExchangeRates
	lastUpdate time.Time
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		client:     &http.Client{Timeout: 10 * time.Second},
		apiBaseURL: "https://api.exchangerate.host/latest",
		cache:      make(map[string]ExchangeRates),
	}
}

func (c *CurrencyConverter) fetchRates(baseCurrency string) (*ExchangeRates, error) {
	if cached, exists := c.cache[baseCurrency]; exists {
		if time.Since(c.lastUpdate) < 30*time.Minute {
			return &cached, nil
		}
	}

	url := fmt.Sprintf("%s?base=%s", c.apiBaseURL, baseCurrency)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rates ExchangeRates
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return nil, err
	}

	c.cache[baseCurrency] = rates
	c.lastUpdate = time.Now()
	return &rates, nil
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	rates, err := c.fetchRates(from)
	if err != nil {
		return 0, err
	}

	rate, exists := rates.Rates[to]
	if !exists {
		return 0, fmt.Errorf("currency %s not found in exchange rates", to)
	}

	return amount * rate, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() ([]string, error) {
	rates, err := c.fetchRates("USD")
	if err != nil {
		return nil, err
	}

	currencies := make([]string, 0, len(rates.Rates)+1)
	currencies = append(currencies, rates.Base)
	for currency := range rates.Rates {
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)

	currencies, err := converter.GetSupportedCurrencies()
	if err != nil {
		fmt.Printf("Error fetching currencies: %v\n", err)
		return
	}

	fmt.Printf("Supported currencies: %v\n", currencies[:10])
}