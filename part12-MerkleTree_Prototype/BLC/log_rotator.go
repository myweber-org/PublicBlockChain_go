
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    filename    string
    currentSize int64
    file        *os.File
}

func NewLogRotator(filename string) (*LogRotator, error) {
    rotator := &LogRotator{filename: filename}
    if err := rotator.openFile(); err != nil {
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

    n, err := lr.file.Write(p)
    if err != nil {
        return n, err
    }
    lr.currentSize += int64(n)
    return n, nil
}

func (lr *LogRotator) openFile() error {
    file, err := os.OpenFile(lr.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.file = file
    lr.currentSize = stat.Size()
    return nil
}

func (lr *LogRotator) rotate() error {
    if lr.file != nil {
        lr.file.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    backupName := fmt.Sprintf("%s.%s", lr.filename, timestamp)

    if err := os.Rename(lr.filename, backupName); err != nil {
        return err
    }

    if err := lr.compressBackup(backupName); err != nil {
        return err
    }

    if err := lr.cleanupOldBackups(); err != nil {
        return err
    }

    return lr.openFile()
}

func (lr *LogRotator) compressBackup(filename string) error {
    source, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer source.Close()

    compressedName := filename + ".gz"
    target, err := os.Create(compressedName)
    if err != nil {
        return err
    }
    defer target.Close()

    gzWriter := gzip.NewWriter(target)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, source); err != nil {
        return err
    }

    os.Remove(filename)
    return nil
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.filename + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    var backupFiles []struct {
        path string
        time time.Time
    }

    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }

        timestamp := parts[len(parts)-2]
        t, err := time.Parse("20060102150405", timestamp)
        if err != nil {
            continue
        }

        backupFiles = append(backupFiles, struct {
            path string
            time time.Time
        }{match, t})
    }

    for i := 0; i < len(backupFiles)-maxBackups; i++ {
        oldestIndex := 0
        for j := 1; j < len(backupFiles); j++ {
            if backupFiles[j].time.Before(backupFiles[oldestIndex].time) {
                oldestIndex = j
            }
        }

        os.Remove(backupFiles[oldestIndex].path)
        backupFiles = append(backupFiles[:oldestIndex], backupFiles[oldestIndex+1:]...)
    }

    return nil
}

func (lr *LogRotator) Close() error {
    if lr.file != nil {
        return lr.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }

        if i%100 == 0 {
            time.Sleep(100 * time.Millisecond)
        }
    }

    fmt.Println("Log rotation test completed")
}