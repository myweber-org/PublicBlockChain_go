package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    filePath    string
    maxSize     int64
    maxAge      time.Duration
    currentFile *os.File
    currentSize int64
    mu          sync.Mutex
}

func NewRotator(filePath string, maxSize int64, maxAge time.Duration) (*Rotator, error) {
    r := &Rotator{
        filePath: filePath,
        maxSize:  maxSize,
        maxAge:   maxAge,
    }
    if err := r.openCurrent(); err != nil {
        return nil, err
    }
    go r.ageMonitor()
    return r, nil
}

func (r *Rotator) openCurrent() error {
    dir := filepath.Dir(r.filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    f, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    stat, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    r.currentFile = f
    r.currentSize = stat.Size()
    return nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := r.currentFile.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }
    timestamp := time.Now().Format("20060102_150405")
    backupPath := r.filePath + "." + timestamp
    if err := os.Rename(r.filePath, backupPath); err != nil {
        return err
    }
    return r.openCurrent()
}

func (r *Rotator) ageMonitor() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()
    for range ticker.C {
        r.mu.Lock()
        if r.currentFile != nil {
            stat, err := r.currentFile.Stat()
            if err == nil && time.Since(stat.ModTime()) > r.maxAge {
                r.rotate()
            }
        }
        r.mu.Unlock()
    }
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
    rotator, err := NewRotator("/var/log/myapp/app.log", 10*1024*1024, 24*time.Hour)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(msg)); err != nil {
            fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }
    fmt.Println("Log rotation demo completed")
}