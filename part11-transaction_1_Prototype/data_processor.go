
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
		if len(row) != 3 {
			return nil, errors.New("invalid CSV format on line " + strconv.Itoa(i+1))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, errors.New("invalid ID on line " + strconv.Itoa(i+1))
		}

		name := row[1]
		if name == "" {
			return nil, errors.New("empty name on line " + strconv.Itoa(i+1))
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, errors.New("invalid value on line " + strconv.Itoa(i+1))
		}

		data = append(data, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return data, nil
}

func ValidateData(records []DataRecord) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return errors.New("invalid ID: " + strconv.Itoa(record.ID))
		}
		if seenIDs[record.ID] {
			return errors.New("duplicate ID: " + strconv.Itoa(record.ID))
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return errors.New("negative value for ID: " + strconv.Itoa(record.ID))
		}
	}
	return nil
}

func ProcessCSVData(filePath string) ([]DataRecord, error) {
	data, err := ParseCSVFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := ValidateData(data); err != nil {
		return nil, err
	}

	return data, nil
}