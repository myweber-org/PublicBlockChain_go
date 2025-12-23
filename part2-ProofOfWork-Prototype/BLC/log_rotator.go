package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    sequence    int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: path,
    }
    err := rl.openCurrentFile()
    return rl, err
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := rl.currentFile.Write(p)
    rl.currentSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
        if err := rl.compressCurrent(); err != nil {
            return err
        }
        rl.cleanOldBackups()
    }
    return rl.openCurrentFile()
}

func (rl *RotatingLogger) openCurrentFile() error {
    f, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    stat, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    rl.currentFile = f
    rl.currentSize = stat.Size()
    return nil
}

func (rl *RotatingLogger) compressCurrent() error {
    timestamp := time.Now().Format("20060102_150405")
    compressedName := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)
    source, err := os.Open(rl.basePath)
    if err != nil {
        return err
    }
    defer source.Close()
    dest, err := os.Create(compressedName)
    if err != nil {
        return err
    }
    defer dest.Close()
    gz := gzip.NewWriter(dest)
    defer gz.Close()
    _, err = io.Copy(gz, source)
    return err
}

func (rl *RotatingLogger) cleanOldBackups() {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    if len(matches) > maxBackups {
        toDelete := matches[:len(matches)-maxBackups]
        for _, f := range toDelete {
            os.Remove(f)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}