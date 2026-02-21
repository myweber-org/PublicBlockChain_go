
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
    filename   string
    current    *os.File
    currentSize int64
}

func NewLogRotator(filename string) (*LogRotator, error) {
    rotator := &LogRotator{filename: filename}
    if err := rotator.openCurrent(); err != nil {
        return nil, err
    }
    return rotator, nil
}

func (lr *LogRotator) openCurrent() error {
    file, err := os.OpenFile(lr.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    lr.current = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := lr.current.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.current.Close(); err != nil {
        return err
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

    return lr.openCurrent()
}

func (lr *LogRotator) compressBackup(filename string) error {
    src, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(filename + ".gz")
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    if err := os.Remove(filename); err != nil {
        return err
    }

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
        oldestIdx := 0
        for j := 1; j < len(backupFiles); j++ {
            if backupFiles[j].time.Before(backupFiles[oldestIdx].time) {
                oldestIdx = j
            }
        }
        if err := os.Remove(backupFiles[oldestIdx].path); err != nil {
            return err
        }
        backupFiles = append(backupFiles[:oldestIdx], backupFiles[oldestIdx+1:]...)
    }

    return nil
}

func (lr *LogRotator) Close() error {
    if lr.current != nil {
        return lr.current.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }
}