package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
}

func parseLogLine(line string) (LogEntry, error) {
	pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 4 {
		return LogEntry{}, fmt.Errorf("invalid log format")
	}

	return LogEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Message:   matches[3],
	}, nil
}

func filterErrors(entries []LogEntry) []LogEntry {
	var errors []LogEntry
	for _, entry := range entries {
		if strings.ToUpper(entry.Level) == "ERROR" {
			errors = append(errors, entry)
		}
	}
	return errors
}

func readLogFile(filename string) ([]LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			fmt.Printf("Warning: Failed to parse line %d: %v\n", lineNumber, err)
		} else {
			entries = append(entries, entry)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func displayErrorReport(errors []LogEntry) {
	fmt.Println("=== ERROR LOG REPORT ===")
	fmt.Printf("Total errors found: %d\n\n", len(errors))
	for i, err := range errors {
		fmt.Printf("%d. [%s] %s\n", i+1, err.Timestamp, err.Message)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile>")
		os.Exit(1)
	}

	filename := os.Args[1]
	entries, err := readLogFile(filename)
	if err != nil {
		fmt.Printf("Error reading log file: %v\n", err)
		os.Exit(1)
	}

	errorEntries := filterErrors(entries)
	displayErrorReport(errorEntries)
}