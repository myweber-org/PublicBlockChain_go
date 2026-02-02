package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

func parseLogLine(line string) (LogEntry, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return LogEntry{}, fmt.Errorf("invalid log format")
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05", parts[0])
	if err != nil {
		return LogEntry{}, err
	}

	return LogEntry{
		Timestamp: timestamp,
		Level:     parts[1],
		Message:   parts[2],
	}, nil
}

func analyzeLogs(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	levelCount := make(map[string]int)
	var earliest, latest time.Time
	errorMessages := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			continue
		}

		levelCount[entry.Level]++

		if earliest.IsZero() || entry.Timestamp.Before(earliest) {
			earliest = entry.Timestamp
		}
		if latest.IsZero() || entry.Timestamp.After(latest) {
			latest = entry.Timestamp
		}

		if entry.Level == "ERROR" {
			errorMessages = append(errorMessages, entry.Message)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("Log Analysis Report\n")
	fmt.Printf("===================\n")
	fmt.Printf("Time Range: %s to %s\n", earliest.Format(time.RFC3339), latest.Format(time.RFC3339))
	fmt.Printf("\nLevel Distribution:\n")
	for level, count := range levelCount {
		fmt.Printf("  %s: %d\n", level, count)
	}

	if len(errorMessages) > 0 {
		fmt.Printf("\nRecent Errors:\n")
		for i, msg := range errorMessages {
			if i >= 5 {
				break
			}
			fmt.Printf("  - %s\n", msg)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: log_analyzer <logfile>")
		os.Exit(1)
	}

	if err := analyzeLogs(os.Args[1]); err != nil {
		fmt.Printf("Error analyzing logs: %v\n", err)
		os.Exit(1)
	}
}