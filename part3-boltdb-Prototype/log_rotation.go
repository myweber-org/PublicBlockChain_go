
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "sync"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    mu         sync.Mutex
    file       *os.File
    size       int64
    basePath   string
    currentDay string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: path,
    }
    if err := rl.rotateIfNeeded(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if err := rl.rotateIfNeeded(); err != nil {
        return 0, err
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
    now := time.Now()
    today := now.Format("2006-01-02")

    if rl.file == nil || rl.currentDay != today || rl.size >= maxFileSize {
        if rl.file != nil {
            rl.file.Close()
            if err := rl.compressOldLog(); err != nil {
                log.Printf("Failed to compress log: %v", err)
            }
        }

        newPath := fmt.Sprintf("%s.%s.log", rl.basePath, today)
        f, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
            return err
        }

        info, err := f.Stat()
        if err != nil {
            f.Close()
            return err
        }

        rl.file = f
        rl.size = info.Size()
        rl.currentDay = today
        rl.cleanupOldBackups()
    }
    return nil
}

func (rl *RotatingLogger) compressOldLog() error {
    oldPath := fmt.Sprintf("%s.%s.log", rl.basePath, rl.currentDay)
    compressedPath := oldPath + ".gz"

    src, err := os.Open(oldPath)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    return os.Remove(oldPath)
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := rl.basePath + ".*.log.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackups {
        toDelete := matches[:len(matches)-maxBackups]
        for _, path := range toDelete {
            os.Remove(path)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(logger)

    for i := 0; i < 100; i++ {
        log.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
        time.Sleep(100 * time.Millisecond)
    }
}