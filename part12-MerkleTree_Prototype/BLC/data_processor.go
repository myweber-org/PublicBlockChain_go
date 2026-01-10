package main

import (
	"errors"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Email string
}

func ValidateUserData(data UserData) error {
	if data.ID <= 0 {
		return errors.New("invalid user ID")
	}
	if strings.TrimSpace(data.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func TransformUserName(data UserData) UserData {
	data.Name = strings.ToUpper(strings.TrimSpace(data.Name))
	return data
}

func ProcessUserInput(rawName, rawEmail string, id int) (UserData, error) {
	user := UserData{
		ID:    id,
		Name:  rawName,
		Email: rawEmail,
	}

	if err := ValidateUserData(user); err != nil {
		return UserData{}, err
	}

	user = TransformUserName(user)
	return user, nil
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ReadCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func CalculateAverage(records []Record) float64 {
	if len(records) == 0 {
		return 0
	}

	var sum float64
	for _, record := range records {
		sum += record.Value
	}

	return sum / float64(len(records))
}

func FilterByThreshold(records []Record, threshold float64) []Record {
	var filtered []Record
	for _, record := range records {
		if record.Value >= threshold {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func main() {
	records, err := ReadCSVFile("data.csv")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Average value: %.2f\n", CalculateAverage(records))

	threshold := 50.0
	filtered := FilterByThreshold(records, threshold)
	fmt.Printf("Records above %.1f: %d\n", threshold, len(filtered))
}