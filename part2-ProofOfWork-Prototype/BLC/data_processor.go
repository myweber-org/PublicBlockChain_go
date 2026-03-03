
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Email    string
	Username string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(p UserProfile) error {
	if !emailRegex.MatchString(p.Email) {
		return errors.New("invalid email format")
	}
	if strings.TrimSpace(p.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if p.Age < 0 || p.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(p *UserProfile) {
	p.Username = strings.ToLower(strings.TrimSpace(p.Username))
}

func ProcessUserProfile(p UserProfile) (UserProfile, error) {
	if err := ValidateProfile(p); err != nil {
		return p, err
	}
	TransformUsername(&p)
	return p, nil
}
package main

import "fmt"

func FilterAndDouble(nums []int, threshold int) []int {
    var result []int
    for _, num := range nums {
        if num > threshold {
            result = append(result, num*2)
        }
    }
    return result
}

func main() {
    input := []int{1, 5, 10, 15, 20}
    filtered := FilterAndDouble(input, 8)
    fmt.Println("Original:", input)
    fmt.Println("Filtered and doubled:", filtered)
}
package main

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	trimmed := strings.TrimSpace(input)

	pattern := `^[a-zA-Z0-9\s\-_\.@]+$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	if !re.MatchString(trimmed) {
		return "", nil
	}

	return trimmed, nil
}