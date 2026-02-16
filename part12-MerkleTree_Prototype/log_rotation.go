package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    currentFile *os.File
    maxSize     int64
    maxFiles    int
    currentSize int64
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := logger.openOrCreate(); err != nil {
        return nil, err
    }
    return logger, nil
}

func (l *RotatingLogger) openOrCreate() error {
    dir := filepath.Dir(l.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    file, err := os.OpenFile(l.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    l.currentFile = file
    l.currentSize = info.Size()
    return nil
}

func (l *RotatingLogger) Write(p []byte) (n int, err error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentSize+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = l.currentFile.Write(p)
    if err == nil {
        l.currentSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) rotate() error {
    if l.currentFile != nil {
        l.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := l.basePath + "." + timestamp

    if err := os.Rename(l.basePath, rotatedPath); err != nil {
        return err
    }

    if err := l.openOrCreate(); err != nil {
        return err
    }

    go l.cleanupOldFiles()
    return nil
}

func (l *RotatingLogger) cleanupOldFiles() {
    dir := filepath.Dir(l.basePath)
    baseName := filepath.Base(l.basePath)

    files, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var logFiles []string
    for _, file := range files {
        name := file.Name()
        if strings.HasPrefix(name, baseName+".") && !file.IsDir() {
            logFiles = append(logFiles, filepath.Join(dir, name))
        }
    }

    if len(logFiles) <= l.maxFiles {
        return
    }

    sort.Strings(logFiles)
    filesToDelete := logFiles[:len(logFiles)-l.maxFiles]

    for _, file := range filesToDelete {
        os.Remove(file)
    }
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(message))
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}