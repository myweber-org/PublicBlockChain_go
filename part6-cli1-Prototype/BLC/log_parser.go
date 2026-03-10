
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

func parseLogLine(line string) (*LogEntry, error) {
	pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if matches == nil {
		return nil, fmt.Errorf("invalid log format")
	}

	return &LogEntry{
		Timestamp: matches[1],
		Level:     strings.ToUpper(matches[2]),
		Message:   matches[3],
	}, nil
}

func extractErrors(logPath string) ([]LogEntry, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var errors []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			continue
		}
		if entry.Level == "ERROR" {
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
		fmt.Printf("Error parsing log: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d error(s):\n", len(errors))
	for _, entry := range errors {
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

func parseLogLine(line string) *LogEntry {
    re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`)
    matches := re.FindStringSubmatch(line)
    
    if matches == nil {
        return nil
    }
    
    return &LogEntry{
        Timestamp: matches[1],
        Level:     matches[2],
        Message:   matches[3],
    }
}

func filterErrors(entries []LogEntry) []LogEntry {
    var errorEntries []LogEntry
    for _, entry := range entries {
        if strings.ToUpper(entry.Level) == "ERROR" {
            errorEntries = append(errorEntries, entry)
        }
    }
    return errorEntries
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
        if entry := parseLogLine(scanner.Text()); entry != nil {
            entries = append(entries, *entry)
        }
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
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }

    errorEntries := filterErrors(entries)
    
    fmt.Printf("Total log entries: %d\n", len(entries))
    fmt.Printf("Error entries: %d\n", len(errorEntries))
    
    for _, entry := range errorEntries {
        fmt.Printf("[%s] %s\n", entry.Timestamp, entry.Message)
    }
}