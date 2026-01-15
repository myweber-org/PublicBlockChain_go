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
    pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
    re := regexp.MustCompile(pattern)
    matches := re.FindStringSubmatch(line)

    if len(matches) == 4 {
        return &LogEntry{
            Timestamp: matches[1],
            Level:     matches[2],
            Message:   matches[3],
        }
    }
    return nil
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

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: log_parser <logfile>")
        return
    }

    filename := os.Args[1]
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer file.Close()

    var entries []LogEntry
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if entry := parseLogLine(scanner.Text()); entry != nil {
            entries = append(entries, *entry)
        }
    }

    errorEntries := filterErrors(entries)
    fmt.Printf("Found %d error entries:\n", len(errorEntries))
    for _, entry := range errorEntries {
        fmt.Printf("[%s] %s\n", entry.Timestamp, entry.Message)
    }
}package main

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

func parseLogLine(line string) (LogEntry, bool) {
	pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 4 {
		return LogEntry{}, false
	}

	return LogEntry{
		Timestamp: matches[1],
		Level:     strings.ToUpper(matches[2]),
		Message:   matches[3],
	}, true
}

func filterErrors(entries []LogEntry) []LogEntry {
	var errors []LogEntry
	for _, entry := range entries {
		if entry.Level == "ERROR" || entry.Level == "FATAL" {
			errors = append(errors, entry)
		}
	}
	return errors
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if entry, ok := parseLogLine(scanner.Text()); ok {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	errorEntries := filterErrors(entries)
	fmt.Printf("Found %d error entries:\n", len(errorEntries))
	for _, entry := range errorEntries {
		fmt.Printf("[%s] %s: %s\n", entry.Timestamp, entry.Level, entry.Message)
	}
}