
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
	Severity  string
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
		Severity:  parts[1],
		Message:   parts[2],
	}, nil
}

func filterLogsBySeverity(entries []LogEntry, severity string) []LogEntry {
	var filtered []LogEntry
	for _, entry := range entries {
		if entry.Severity == severity {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func filterLogsByTimeRange(entries []LogEntry, start, end time.Time) []LogEntry {
	var filtered []LogEntry
	for _, entry := range entries {
		if !entry.Timestamp.Before(start) && !entry.Timestamp.After(end) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func readLogFile(filename string) ([]LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}
	return entries, scanner.Err()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile>")
		os.Exit(1)
	}

	entries, err := readLogFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading log file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total log entries: %d\n", len(entries))

	errorLogs := filterLogsBySeverity(entries, "ERROR")
	fmt.Printf("Error entries: %d\n", len(errorLogs))

	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()
	recentLogs := filterLogsByTimeRange(entries, startTime, endTime)
	fmt.Printf("Last 24 hour entries: %d\n", len(recentLogs))
}