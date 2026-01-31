
package main

import (
	"fmt"
	"math"
)

// FilterAndTransform processes a slice of integers, filters out values below threshold,
// and applies a transformation (square root of absolute value).
func FilterAndTransform(numbers []int, threshold int) []float64 {
	var result []float64
	for _, num := range numbers {
		if num > threshold {
			transformed := math.Sqrt(math.Abs(float64(num)))
			result = append(result, transformed)
		}
	}
	return result
}

func main() {
	input := []int{-10, 5, 3, 15, 8, -2, 25}
	threshold := 5
	output := FilterAndTransform(input, threshold)
	fmt.Printf("Processed slice: %v\n", output)
}
package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Comments string
}

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(trimmed, "")
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ProcessUserData(data UserData) (UserData, error) {
	sanitizedData := UserData{
		Username: SanitizeInput(data.Username),
		Email:    SanitizeInput(data.Email),
		Comments: SanitizeInput(data.Comments),
	}

	if !ValidateEmail(sanitizedData.Email) {
		return sanitizedData, &InvalidEmailError{Email: sanitizedData.Email}
	}

	if len(sanitizedData.Username) < 3 {
		return sanitizedData, &InvalidUsernameError{Username: sanitizedData.Username}
	}

	return sanitizedData, nil
}

type InvalidEmailError struct {
	Email string
}

func (e *InvalidEmailError) Error() string {
	return "Invalid email format: " + e.Email
}

type InvalidUsernameError struct {
	Username string
}

func (e *InvalidUsernameError) Error() string {
	return "Username must be at least 3 characters long: " + e.Username
}