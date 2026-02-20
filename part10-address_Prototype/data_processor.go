
package data

import (
	"errors"
	"regexp"
	"strings"
)

type Record struct {
	ID    string
	Email string
	Score int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateRecord(r Record) error {
	if r.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}
	if r.Score < 0 || r.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func TransformRecords(records []Record) ([]Record, error) {
	var processed []Record
	for _, r := range records {
		if err := ValidateRecord(r); err != nil {
			return nil, err
		}
		r.Email = NormalizeEmail(r.Email)
		processed = append(processed, r)
	}
	return processed, nil
}

func CalculateAverage(records []Record) float64 {
	if len(records) == 0 {
		return 0.0
	}
	total := 0
	for _, r := range records {
		total += r.Score
	}
	return float64(total) / float64(len(records))
}