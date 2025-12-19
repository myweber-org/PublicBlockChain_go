package main

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

type Record struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func ParseCSVData(reader io.Reader) ([]Record, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var result []Record
	for i, row := range records {
		if len(row) != 4 {
			return nil, errors.New("invalid column count at row " + strconv.Itoa(i))
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, errors.New("invalid ID at row " + strconv.Itoa(i))
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			return nil, errors.New("empty name at row " + strconv.Itoa(i))
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, errors.New("invalid value at row " + strconv.Itoa(i))
		}

		active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
		if err != nil {
			return nil, errors.New("invalid active flag at row " + strconv.Itoa(i))
		}

		result = append(result, Record{
			ID:     id,
			Name:   name,
			Value:  value,
			Active: active,
		})
	}

	return result, nil
}

func ValidateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	for _, rec := range records {
		if rec.ID <= 0 {
			return errors.New("invalid ID: " + strconv.Itoa(rec.ID))
		}
		if seenIDs[rec.ID] {
			return errors.New("duplicate ID: " + strconv.Itoa(rec.ID))
		}
		seenIDs[rec.ID] = true

		if rec.Value < 0 {
			return errors.New("negative value for ID: " + strconv.Itoa(rec.ID))
		}
	}
	return nil
}

func CalculateTotal(records []Record) float64 {
	var total float64
	for _, rec := range records {
		if rec.Active {
			total += rec.Value
		}
	}
	return total
}