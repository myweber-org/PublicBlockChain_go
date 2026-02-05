package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strconv"
)

type ExchangeRates struct {
    Rates map[string]float64 `json:"rates"`
    Base  string             `json:"base"`
    Date  string             `json:"date"`
}

func fetchExchangeRates(apiKey string) (*ExchangeRates, error) {
    url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/USD")
    if apiKey != "" {
        url = fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/USD", apiKey)
    }

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var rates ExchangeRates
    err = json.Unmarshal(body, &rates)
    if err != nil {
        return nil, err
    }

    return &rates, nil
}

func convertCurrency(amount float64, fromCurrency, toCurrency string, rates *ExchangeRates) (float64, error) {
    fromRate, okFrom := rates.Rates[fromCurrency]
    toRate, okTo := rates.Rates[toCurrency]

    if !okFrom || !okTo {
        return 0, fmt.Errorf("invalid currency code")
    }

    amountInUSD := amount / fromRate
    convertedAmount := amountInUSD * toRate

    return convertedAmount, nil
}

func main() {
    apiKey := os.Getenv("EXCHANGE_RATE_API_KEY")

    rates, err := fetchExchangeRates(apiKey)
    if err != nil {
        fmt.Printf("Error fetching exchange rates: %v\n", err)
        return
    }

    if len(os.Args) != 4 {
        fmt.Println("Usage: currency_converter <amount> <from_currency> <to_currency>")
        fmt.Println("Example: currency_converter 100 USD EUR")
        return
    }

    amount, err := strconv.ParseFloat(os.Args[1], 64)
    if err != nil {
        fmt.Printf("Invalid amount: %v\n", err)
        return
    }

    fromCurrency := os.Args[2]
    toCurrency := os.Args[3]

    result, err := convertCurrency(amount, fromCurrency, toCurrency, rates)
    if err != nil {
        fmt.Printf("Conversion error: %v\n", err)
        return
    }

    fmt.Printf("%.2f %s = %.2f %s (as of %s)\n", amount, fromCurrency, result, toCurrency, rates.Date)
}