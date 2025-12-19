package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	Name  string  `json:"name"`
	Age   int     `json:"age"`
	Score float64 `json:"score"`
}

func readCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	if _, err := reader.Read(); err != nil {
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

		age, _ := strconv.Atoi(row[1])
		score, _ := strconv.ParseFloat(row[2], 64)

		records = append(records, Record{
			Name:  row[0],
			Age:   age,
			Score: score,
		})
	}

	return records, nil
}

func convertToJSON(records []Record) (string, error) {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func processData(filename string) error {
	records, err := readCSV(filename)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("no records found")
	}

	jsonData, err := convertToJSON(records)
	if err != nil {
		return err
	}

	fmt.Println("Processed", len(records), "records")
	fmt.Println(jsonData)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	if err := processData(os.Args[1]); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}