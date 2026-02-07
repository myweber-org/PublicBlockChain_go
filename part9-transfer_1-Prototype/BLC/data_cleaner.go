package datautils

import "sort"

// Deduplicate removes duplicate values from a slice of comparable types
func Deduplicate[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(input))

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// DeduplicateSorted removes duplicates from a sorted slice more efficiently
func DeduplicateSorted[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	result := make([]T, 0, len(input))
	result = append(result, input[0])

	for i := 1; i < len(input); i++ {
		if input[i] != input[i-1] {
			result = append(result, input[i])
		}
	}

	return result
}

// DeduplicateWithSort sorts and removes duplicates from a slice
func DeduplicateWithSort[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	// Create a copy to avoid modifying original
	sorted := make([]T, len(input))
	copy(sorted, input)

	// Sort the slice
	sort.Slice(sorted, func(i, j int) bool {
		// Use type assertion for comparison
		// This is a simplified comparison - in practice you'd need
		// a more sophisticated approach for generic types
		return false // Placeholder - actual implementation depends on type
	})

	return DeduplicateSorted(sorted)
}package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func deduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func validateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func processRecords(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	var cleaned []DataRecord

	for _, rec := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(rec.Email))
		if validateEmail(cleanEmail) && !emailSet[cleanEmail] {
			emailSet[cleanEmail] = true
			rec.Email = cleanEmail
			rec.Valid = true
			cleaned = append(cleaned, rec)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain", false},
		{5, "user@example.com", false},
	}

	cleaned := processRecords(records)
	fmt.Printf("Processed %d records, %d valid unique records found\n", len(records), len(cleaned))
	for _, rec := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", rec.ID, rec.Email, rec.Valid)
	}
}