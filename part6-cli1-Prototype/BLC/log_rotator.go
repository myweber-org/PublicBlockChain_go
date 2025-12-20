package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    filePath     string
    maxSize      int64
    rotationTime time.Duration
    currentSize  int64
    lastRotation time.Time
    mu           sync.Mutex
    file         *os.File
}

func NewRotator(filePath string, maxSize int64, rotationTime time.Duration) (*Rotator, error) {
    r := &Rotator{
        filePath:     filePath,
        maxSize:      maxSize,
        rotationTime: rotationTime,
        lastRotation: time.Now(),
    }

    if err := r.openFile(); err != nil {
        return nil, err
    }

    return r, nil
}

func (r *Rotator) openFile() error {
    dir := filepath.Dir(r.filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    r.file = file
    r.currentSize = info.Size()
    return nil
}

func (r *Rotator) rotate() error {
    if r.file != nil {
        r.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)

    if err := os.Rename(r.filePath, backupPath); err != nil {
        return err
    }

    r.lastRotation = time.Now()
    return r.openFile()
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    now := time.Now()
    if r.currentSize+int64(len(p)) > r.maxSize || now.Sub(r.lastRotation) > r.rotationTime {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.file.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.file != nil {
        return r.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("logs/app.log", 1024*1024, 24*time.Hour)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        rotator.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
}