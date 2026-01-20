package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
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
    Errors       []string
}

func parseLogLine(line string) (LogEntry, error) {
    re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)\s+(\w+)\s+(.+)$`)
    matches := re.FindStringSubmatch(line)
    
    if len(matches) != 4 {
        return LogEntry{}, fmt.Errorf("invalid log format")
    }
    
    timestamp, err := time.Parse(time.RFC3339, matches[1])
    if err != nil {
        return LogEntry{}, err
    }
    
    return LogEntry{
        Timestamp: timestamp,
        Level:     matches[2],
        Message:   matches[3],
    }, nil
}

func analyzeLogs(filePath string) (LogSummary, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return LogSummary{}, err
    }
    defer file.Close()

    summary := LogSummary{
        LevelCounts: make(map[string]int),
        Errors:      []string{},
    }

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        entry, err := parseLogLine(scanner.Text())
        if err != nil {
            continue
        }

        summary.TotalEntries++
        summary.LevelCounts[entry.Level]++

        if entry.Level == "ERROR" {
            summary.Errors = append(summary.Errors, entry.Message)
        }
    }

    return summary, scanner.Err()
}

func printSummary(summary LogSummary) {
    fmt.Printf("Log Analysis Report\n")
    fmt.Printf("===================\n")
    fmt.Printf("Total entries: %d\n", summary.TotalEntries)
    
    fmt.Printf("\nLevel distribution:\n")
    for level, count := range summary.LevelCounts {
        percentage := float64(count) / float64(summary.TotalEntries) * 100
        fmt.Printf("  %-6s: %d (%.1f%%)\n", level, count, percentage)
    }
    
    if len(summary.Errors) > 0 {
        fmt.Printf("\nError messages found (%d):\n", len(summary.Errors))
        for i, err := range summary.Errors {
            fmt.Printf("  %d. %s\n", i+1, err)
        }
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: log_analyzer <log_file>")
        os.Exit(1)
    }

    summary, err := analyzeLogs(os.Args[1])
    if err != nil {
        fmt.Printf("Error analyzing logs: %v\n", err)
        os.Exit(1)
    }

    printSummary(summary)
}