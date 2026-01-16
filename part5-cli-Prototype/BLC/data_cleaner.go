
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(input)
	fmt.Printf("Original: %v\n", input)
	fmt.Printf("Cleaned: %v\n", cleaned)
}package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    ID    int
    Name  string
    Score float64
    Valid bool
}

func cleanCSV(inputPath, outputPath string) error {
    inFile, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer inFile.Close()

    outFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer outFile.Close()

    reader := csv.NewReader(inFile)
    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    header, err := reader.Read()
    if err != nil {
        return err
    }
    header = append(header, "Valid")
    writer.Write(header)

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            continue
        }

        record := parseRecord(row)
        outputRow := []string{
            strconv.Itoa(record.ID),
            strings.TrimSpace(record.Name),
            fmt.Sprintf("%.2f", record.Score),
            strconv.FormatBool(record.Valid),
        }
        writer.Write(outputRow)
    }

    return nil
}

func parseRecord(row []string) Record {
    if len(row) < 3 {
        return Record{Valid: false}
    }

    id, err := strconv.Atoi(row[0])
    if err != nil {
        return Record{Valid: false}
    }

    name := row[1]
    if name == "" {
        return Record{Valid: false}
    }

    score, err := strconv.ParseFloat(row[2], 64)
    if err != nil || score < 0 || score > 100 {
        return Record{Valid: false}
    }

    return Record{
        ID:    id,
        Name:  name,
        Score: score,
        Valid: true,
    }
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    err := cleanCSV(os.Args[1], os.Args[2])
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}