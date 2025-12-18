
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
    maxFileSize = 10 * 1024 * 1024
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
    rl.size += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
    now := time.Now()
    dateStr := now.Format("2006-01-02")

    if rl.file == nil || rl.currentDay != dateStr || rl.size >= maxFileSize {
        if rl.file != nil {
            rl.file.Close()
            if err := rl.compressOldLog(); err != nil {
                log.Printf("Failed to compress log: %v", err)
            }
        }

        newPath := fmt.Sprintf("%s.%s.log", rl.basePath, dateStr)
        file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
            return err
        }

        stat, err := file.Stat()
        if err != nil {
            file.Close()
            return err
        }

        rl.file = file
        rl.size = stat.Size()
        rl.currentDay = dateStr
        rl.cleanupOldBackups()
    }
    return nil
}

func (rl *RotatingLogger) compressOldLog() error {
    files, err := filepath.Glob(rl.basePath + ".*.log")
    if err != nil {
        return err
    }

    for _, f := range files {
        if filepath.Ext(f) == ".gz" {
            continue
        }

        gzPath := f + ".gz"
        if _, err := os.Stat(gzPath); err == nil {
            continue
        }

        src, err := os.Open(f)
        if err != nil {
            return err
        }

        dst, err := os.Create(gzPath)
        if err != nil {
            src.Close()
            return err
        }

        gz := gzip.NewWriter(dst)
        if _, err := io.Copy(gz, src); err != nil {
            src.Close()
            gz.Close()
            dst.Close()
            return err
        }

        src.Close()
        gz.Close()
        dst.Close()
        os.Remove(f)
    }
    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
    files, err := filepath.Glob(rl.basePath + ".*.log.gz")
    if err != nil {
        return
    }

    if len(files) > maxBackups {
        for i := 0; i < len(files)-maxBackups; i++ {
            os.Remove(files[i])
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

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
        time.Sleep(100 * time.Millisecond)
    }
}