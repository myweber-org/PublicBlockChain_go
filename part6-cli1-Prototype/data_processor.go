
package main

import (
    "fmt"
    "strings"
)

type UserData struct {
    Username string
    Email    string
}

func normalizeUsername(username string) string {
    return strings.ToLower(strings.TrimSpace(username))
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processUserData(username, email string) (*UserData, error) {
    normalizedUsername := normalizeUsername(username)
    
    if !validateEmail(email) {
        return nil, fmt.Errorf("invalid email format")
    }
    
    return &UserData{
        Username: normalizedUsername,
        Email:    strings.ToLower(strings.TrimSpace(email)),
    }, nil
}

func main() {
    user, err := processUserData("  JohnDoe  ", "john@example.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Printf("Processed user: %+v\n", user)
}