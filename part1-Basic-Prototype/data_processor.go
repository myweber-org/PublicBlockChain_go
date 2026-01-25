
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

type UserData struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
    IsActive bool   `json:"is_active"`
}

func ValidateAndTransform(rawData string) (*UserData, error) {
    var data map[string]interface{}
    err := json.Unmarshal([]byte(rawData), &data)
    if err != nil {
        return nil, fmt.Errorf("invalid JSON format: %v", err)
    }

    processed := &UserData{}

    if name, ok := data["name"].(string); ok {
        processed.Name = strings.TrimSpace(name)
        if processed.Name == "" {
            return nil, fmt.Errorf("name cannot be empty")
        }
    } else {
        return nil, fmt.Errorf("name field missing or invalid")
    }

    if email, ok := data["email"].(string); ok {
        processed.Email = strings.ToLower(strings.TrimSpace(email))
        if !strings.Contains(processed.Email, "@") {
            return nil, fmt.Errorf("invalid email format")
        }
    } else {
        return nil, fmt.Errorf("email field missing or invalid")
    }

    if age, ok := data["age"].(float64); ok {
        if age < 0 || age > 150 {
            return nil, fmt.Errorf("age must be between 0 and 150")
        }
        processed.Age = int(age)
    } else {
        return nil, fmt.Errorf("age field missing or invalid")
    }

    if isActive, ok := data["is_active"].(bool); ok {
        processed.IsActive = isActive
    } else {
        processed.IsActive = true
    }

    return processed, nil
}

func main() {
    testData := `{"name": "John Doe", "email": "JOHN@EXAMPLE.COM", "age": 30, "is_active": true}`
    result, err := ValidateAndTransform(testData)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Processed data: %+v\n", result)
}