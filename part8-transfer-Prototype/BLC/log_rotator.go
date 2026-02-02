
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
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logExtension = ".log"
    gzipExt      = ".gz"
)

type LogRotator struct {
    currentSize int64
    file        *os.File
    basePath    string
    sequence    int
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: strings.TrimSuffix(basePath, logExtension),
    }

    if err := rotator.openCurrent(); err != nil {
        return nil, err
    }

    rotator.cleanOldBackups()
    return rotator, nil
}

func (lr *LogRotator) openCurrent() error {
    path := lr.basePath + logExtension
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.file = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s%s", lr.basePath, timestamp, logExtension)

    if err := os.Rename(lr.basePath+logExtension, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    return lr.openCurrent()
}

func (lr *LogRotator) compressFile(path string) error {
    src, err := os.Open(path)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(path + gzipExt)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    if err := os.Remove(path); err != nil {
        return err
    }

    return nil
}

func (lr *LogRotator) cleanOldBackups() {
    pattern := lr.basePath + ".*" + logExtension + gzipExt
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    for i := 0; i < len(matches)-maxBackups; i++ {
        os.Remove(matches[i])
    }
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
        logEntry := fmt.Sprintf("[%s] Iteration %d: Sample log entry\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    maxAge      time.Duration
    currentFile *os.File
    currentSize int64
    createdAt   time.Time
}

func NewRotator(filePath string, maxSize int64, maxAge time.Duration) (*Rotator, error) {
    r := &Rotator{
        filePath: filePath,
        maxSize:  maxSize,
        maxAge:   maxAge,
    }
    if err := r.openOrCreate(); err != nil {
        return nil, err
    }
    return r, nil
}

func (r *Rotator) openOrCreate() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentFile != nil {
        r.currentFile.Close()
    }

    info, err := os.Stat(r.filePath)
    if os.IsNotExist(err) {
        file, err := os.Create(r.filePath)
        if err != nil {
            return err
        }
        r.currentFile = file
        r.currentSize = 0
        r.createdAt = time.Now()
        return nil
    }
    if err != nil {
        return err
    }

    file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    r.currentFile = file
    r.currentSize = info.Size()
    r.createdAt = info.ModTime()
    return nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.shouldRotateLocked(len(p)) {
        if err := r.rotateLocked(); err != nil {
            return 0, err
        }
    }

    n, err := r.currentFile.Write(p)
    if err != nil {
        return n, err
    }
    r.currentSize += int64(n)
    return n, nil
}

func (r *Rotator) shouldRotateLocked(appendSize int) bool {
    if r.currentSize+int64(appendSize) > r.maxSize {
        return true
    }
    if time.Since(r.createdAt) > r.maxAge {
        return true
    }
    return false
}

func (r *Rotator) rotateLocked() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    ext := filepath.Ext(r.filePath)
    base := r.filePath[:len(r.filePath)-len(ext)]
    archivePath := fmt.Sprintf("%s_%s%s", base, timestamp, ext)

    if err := os.Rename(r.filePath, archivePath); err != nil {
        return err
    }

    file, err := os.Create(r.filePath)
    if err != nil {
        return err
    }
    r.currentFile = file
    r.currentSize = 0
    r.createdAt = time.Now()
    return nil
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("app.log", 1024*1024, 24*time.Hour)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(msg)); err != nil {
            panic(err)
        }
        time.Sleep(100 * time.Millisecond)
    }
}