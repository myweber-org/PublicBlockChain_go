
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        if len(row) != 3 {
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

        records = append(records, DataRecord{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    for _, record := range records {
        sum += record.Value
    }

    average := sum / float64(len(records))

    var variance float64
    for _, record := range records {
        diff := record.Value - average
        variance += diff * diff
    }
    variance = variance / float64(len(records))

    return average, variance
}

func FilterRecords(records []DataRecord, minValue float64) []DataRecord {
    filtered := make([]DataRecord, 0)
    for _, record := range records {
        if record.Value >= minValue {
            filtered = append(filtered, record)
        }
    }
    return filtered
}