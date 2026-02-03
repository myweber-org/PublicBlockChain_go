package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
)

const maxFileSize = 1024 * 1024 // 1MB

type RotatingLogger struct {
    mu       sync.Mutex
    filePath string
    file     *os.File
    size     int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{filePath: path}
    if err := rl.openFile(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) openFile() error {
    file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    rl.file = file
    rl.size = info.Size()
    return nil
}

func (rl *RotatingLogger) rotate() error {
    rl.file.Close()
    for i := 1; ; i++ {
        backupPath := rl.filePath + "." + strconv.Itoa(i)
        if _, err := os.Stat(backupPath); os.IsNotExist(err) {
            if err := os.Rename(rl.filePath, backupPath); err != nil {
                return err
            }
            break
        }
    }
    return rl.openFile()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
        rl.size = 0
    }
    n, err := rl.file.Write(p)
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
        os.Exit(1)
    }
    defer logger.Close()
    for i := 0; i < 10000; i++ {
        msg := fmt.Sprintf("Log entry %d: Some sample data for testing rotation\n", i)
        if _, err := logger.Write([]byte(msg)); err != nil {
            fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
            break
        }
    }
    fmt.Println("Log rotation test completed")
}