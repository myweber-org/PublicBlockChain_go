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

func parseLogLine(line string) *LogEntry {
	pattern := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+(\w+)\s+(.+)$`)
	matches := pattern.FindStringSubmatch(line)

	if matches == nil {
		return nil
	}

	return &LogEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Message:   matches[3],
	}
}

func extractErrors(logFile string) ([]LogEntry, error) {
	file, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var errors []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		entry := parseLogLine(scanner.Text())
		if entry != nil && strings.ToUpper(entry.Level) == "ERROR" {
			errors = append(errors, *entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return errors, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile>")
		os.Exit(1)
	}

	errors, err := extractErrors(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing log file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d error entries:\n", len(errors))
	for i, entry := range errors {
		fmt.Printf("%d. [%s] %s\n", i+1, entry.Timestamp, entry.Message)
	}
}