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

type LogSummary struct {
	TotalEntries int
	LevelCounts  map[string]int
	StartTime    time.Time
	EndTime      time.Time
}

func parseLogLine(line string) (LogEntry, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return LogEntry{}, fmt.Errorf("invalid log format")
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05Z", parts[0])
	if err != nil {
		return LogEntry{}, err
	}

	return LogEntry{
		Timestamp: timestamp,
		Level:     parts[1],
		Message:   parts[2],
	}, nil
}

func analyzeLogs(filepath string) (LogSummary, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return LogSummary{}, err
	}
	defer file.Close()

	summary := LogSummary{
		LevelCounts: make(map[string]int),
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			continue
		}

		summary.TotalEntries++
		summary.LevelCounts[entry.Level]++

		if summary.StartTime.IsZero() || entry.Timestamp.Before(summary.StartTime) {
			summary.StartTime = entry.Timestamp
		}
		if entry.Timestamp.After(summary.EndTime) {
			summary.EndTime = entry.Timestamp
		}
	}

	return summary, scanner.Err()
}

func printSummary(summary LogSummary) {
	fmt.Printf("Log Analysis Summary:\n")
	fmt.Printf("Total entries: %d\n", summary.TotalEntries)
	fmt.Printf("Time range: %s to %s\n", summary.StartTime.Format(time.RFC3339), summary.EndTime.Format(time.RFC3339))
	fmt.Println("Level distribution:")
	for level, count := range summary.LevelCounts {
		percentage := float64(count) / float64(summary.TotalEntries) * 100
		fmt.Printf("  %s: %d (%.1f%%)\n", level, count, percentage)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_analyzer <logfile>")
		os.Exit(1)
	}

	summary, err := analyzeLogs(os.Args[1])
	if err != nil {
		fmt.Printf("Error analyzing logs: %v\n", err)
		os.Exit(1)
	}

	printSummary(summary)
}