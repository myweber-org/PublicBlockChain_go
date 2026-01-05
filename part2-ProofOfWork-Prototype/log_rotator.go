
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logExtension = ".log"
)

type LogRotator struct {
    filename    string
    currentSize int64
}

func NewLogRotator(filename string) (*LogRotator, error) {
    rotator := &LogRotator{filename: filename}
    
    if err := rotator.ensureLogFile(); err != nil {
        return nil, err
    }
    
    if err := rotator.loadCurrentSize(); err != nil {
        return nil, err
    }
    
    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    
    file, err := os.OpenFile(lr.filename, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer file.Close()
    
    n, err := file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.cleanupOldBackups(); err != nil {
        return err
    }
    
    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s%s", strings.TrimSuffix(lr.filename, logExtension), timestamp, logExtension)
    
    if err := os.Rename(lr.filename, backupName); err != nil {
        return err
    }
    
    lr.currentSize = 0
    return lr.ensureLogFile()
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := fmt.Sprintf("%s.*%s", strings.TrimSuffix(lr.filename, logExtension), logExtension)
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }
    
    if len(matches) <= maxBackups {
        return nil
    }
    
    sort.Strings(matches)
    filesToRemove := matches[:len(matches)-maxBackups]
    
    for _, file := range filesToRemove {
        if err := os.Remove(file); err != nil {
            return err
        }
    }
    
    return nil
}

func (lr *LogRotator) ensureLogFile() error {
    file, err := os.OpenFile(lr.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    return file.Close()
}

func (lr *LogRotator) loadCurrentSize() error {
    info, err := os.Stat(lr.filename)
    if err != nil {
        if os.IsNotExist(err) {
            lr.currentSize = 0
            return nil
        }
        return err
    }
    lr.currentSize = info.Size()
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    
    for i := 1; i <= 100; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }
        
        if i%10 == 0 {
            fmt.Printf("Written %d log entries\n", i)
        }
        
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}