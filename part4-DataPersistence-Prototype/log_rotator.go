package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type RotatingLogger struct {
    currentFile   *os.File
    basePath      string
    maxSize       int64
    currentSize   int64
    rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    
    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }
    
    return &RotatingLogger{
        currentFile:   file,
        basePath:      basePath,
        maxSize:       maxSize,
        currentSize:   info.Size(),
        rotationCount: 0,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    rl.currentFile.Close()
    
    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s.%d", rl.basePath, timestamp, rl.rotationCount)
    
    err := os.Rename(rl.basePath, rotatedPath)
    if err != nil {
        return err
    }
    
    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    rl.currentFile = file
    rl.currentSize = 0
    rl.rotationCount++
    
    rl.cleanOldLogs()
    
    return nil
}

func (rl *RotatingLogger) cleanOldLogs() {
    pattern := rl.basePath + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    if len(matches) > 10 {
        for i := 0; i < len(matches)-10; i++ {
            os.Remove(matches[i])
        }
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    
    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation completed")
}