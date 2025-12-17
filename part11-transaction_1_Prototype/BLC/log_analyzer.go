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

func analyzeLogs(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    levelCount := make(map[string]int)
    var entries []LogEntry

    for scanner.Scan() {
        entry, err := parseLogLine(scanner.Text())
        if err != nil {
            continue
        }
        entries = append(entries, entry)
        levelCount[entry.Level]++
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    fmt.Println("Log Analysis Report")
    fmt.Println("===================")
    fmt.Printf("Total entries: %d\n", len(entries))

    if len(entries) > 0 {
        fmt.Printf("Time range: %s to %s\n",
            entries[0].Timestamp.Format("2006-01-02 15:04:05"),
            entries[len(entries)-1].Timestamp.Format("2006-01-02 15:04:05"))
    }

    fmt.Println("\nLog Level Distribution:")
    for level, count := range levelCount {
        percentage := float64(count) / float64(len(entries)) * 100
        fmt.Printf("%-8s: %d (%.1f%%)\n", level, count, percentage)
    }

    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: log_analyzer <logfile>")
        os.Exit(1)
    }

    if err := analyzeLogs(os.Args[1]); err != nil {
        fmt.Printf("Error analyzing logs: %v\n", err)
        os.Exit(1)
    }
}