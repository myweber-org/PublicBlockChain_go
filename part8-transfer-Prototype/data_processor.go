package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

func ValidateJSON(data []byte) (*User, error) {
	var user User
	err := json.Unmarshal(data, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if user.Username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if user.ID <= 0 {
		return nil, fmt.Errorf("ID must be a positive integer")
	}

	return &user, nil
}

func main() {
	jsonData := []byte(`{"id": 123, "username": "johndoe", "email": "john@example.com", "active": true}`)

	user, err := ValidateJSON(jsonData)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	fmt.Printf("Validated user: %+v\n", user)
}