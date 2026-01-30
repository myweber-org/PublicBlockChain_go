package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package datautil

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Record struct {
	ID    string
	Email string
	Phone string
}

type Cleaner struct {
	seenHashes map[string]bool
}

func NewCleaner() *Cleaner {
	return &Cleaner{
		seenHashes: make(map[string]bool),
	}
}

func (c *Cleaner) NormalizeEmail(email string) string {
	parts := strings.Split(strings.ToLower(email), "@")
	if len(parts) != 2 {
		return ""
	}
	local := strings.Split(parts[0], "+")[0]
	local = strings.ReplaceAll(local, ".", "")
	return local + "@" + parts[1]
}

func (c *Cleaner) GenerateHash(record Record) string {
	normalizedEmail := c.NormalizeEmail(record.Email)
	phoneDigits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, record.Phone)

	data := fmt.Sprintf("%s|%s", normalizedEmail, phoneDigits)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (c *Cleaner) IsDuplicate(record Record) bool {
	hash := c.GenerateHash(record)
	if c.seenHashes[hash] {
		return true
	}
	c.seenHashes[hash] = true
	return false
}

func (c *Cleaner) ValidateRecord(record Record) bool {
	if record.ID == "" {
		return false
	}
	if c.NormalizeEmail(record.Email) == "" {
		return false
	}
	if len(record.Phone) < 10 {
		return false
	}
	return true
}

func (c *Cleaner) ProcessRecords(records []Record) []Record {
	var cleaned []Record
	for _, rec := range records {
		if !c.ValidateRecord(rec) {
			continue
		}
		if c.IsDuplicate(rec) {
			continue
		}
		cleaned = append(cleaned, rec)
	}
	return cleaned
}

func main() {
	cleaner := NewCleaner()
	
	records := []Record{
		{ID: "1", Email: "test@example.com", Phone: "1234567890"},
		{ID: "2", Email: "TEST@example.com", Phone: "1234567890"},
		{ID: "3", Email: "test+tag@example.com", Phone: "1234567890"},
		{ID: "4", Email: "invalid-email", Phone: "1234567890"},
		{ID: "5", Email: "another@test.com", Phone: "9876543210"},
		{ID: "6", Email: "test.work@example.com", Phone: "1234567890"},
		{ID: "", Email: "emptyid@test.com", Phone: "1112223333"},
	}
	
	cleaned := cleaner.ProcessRecords(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
	
	for _, rec := range cleaned {
		fmt.Printf("ID: %s, Email: %s\n", rec.ID, rec.Email)
	}
}