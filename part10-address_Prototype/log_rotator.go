
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

const (
    maxFileSize   = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logExtension  = ".log"
    gzipExtension = ".gz"
)

type RotatingLogger struct {
    filename    string
    currentSize int64
    file        *os.File
    mu          sync.Mutex
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        filename: filename,
    }

    if err := rl.openOrCreate(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openOrCreate() error {
    file, err := os.OpenFile(rl.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.file = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s%s", rl.filename, timestamp, logExtension)

    if err := os.Rename(rl.filename, backupName); err != nil {
        return err
    }

    if err := rl.openOrCreate(); err != nil {
        return err
    }

    go rl.compressAndCleanup(backupName)
    return nil
}

func (rl *RotatingLogger) compressAndCleanup(backupName string) {
    compressedName := backupName + gzipExtension
    if err := compressFile(backupName, compressedName); err != nil {
        fmt.Printf("Failed to compress %s: %v\n", backupName, err)
        return
    }

    if err := os.Remove(backupName); err != nil {
        fmt.Printf("Failed to remove %s: %v\n", backupName, err)
    }

    rl.cleanupOldBackups()
}

func compressFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()

    gz := gzip.NewWriter(destination)
    defer gz.Close()

    _, err = io.Copy(gz, source)
    return err
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := rl.filename + ".*" + logExtension + gzipExtension
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackupFiles {
        return
    }

    sortByTimestamp(matches)
    for i := 0; i < len(matches)-maxBackupFiles; i++ {
        os.Remove(matches[i])
    }
}

func sortByTimestamp(files []string) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            if extractTimestamp(files[i]) > extractTimestamp(files[j]) {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func extractTimestamp(filename string) string {
    base := filepath.Base(filename)
    parts := strings.Split(base, ".")
    if len(parts) < 3 {
        return ""
    }
    return parts[1]
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("application.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 1; i <= 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := logger.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}