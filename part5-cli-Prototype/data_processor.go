package main

import (
	"errors"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	for _, r := range data.Username {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' {
			return errors.New("username can only contain letters, numbers, and underscores")
		}
	}

	if !strings.Contains(data.Email, "@") || !strings.Contains(data.Email, ".") {
		return errors.New("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(rawUsername string, rawEmail string, rawAge int) (UserData, error) {
	normalizedUsername := NormalizeUsername(rawUsername)

	userData := UserData{
		Username: normalizedUsername,
		Email:    strings.TrimSpace(rawEmail),
		Age:      rawAge,
	}

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	return userData, nil
}
package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ParseCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataRecord
	for i, row := range records {
		if len(row) < 3 {
			return nil, errors.New("invalid row length at line " + strconv.Itoa(i+1))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, errors.New("invalid ID format at line " + strconv.Itoa(i+1))
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, errors.New("invalid value format at line " + strconv.Itoa(i+1))
		}

		data = append(data, DataRecord{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
	}

	return data, nil
}

func ValidateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return errors.New("invalid ID: " + strconv.Itoa(record.ID))
		}
		if record.Name == "" {
			return errors.New("empty name for ID: " + strconv.Itoa(record.ID))
		}
		if record.Value < 0 {
			return errors.New("negative value for ID: " + strconv.Itoa(record.ID))
		}
		if seenIDs[record.ID] {
			return errors.New("duplicate ID found: " + strconv.Itoa(record.ID))
		}
		seenIDs[record.ID] = true
	}
	return nil
}

func ProcessCSVData(filePath string) ([]DataRecord, error) {
	records, err := ParseCSVFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := ValidateRecords(records); err != nil {
		return nil, err
	}

	return records, nil
}