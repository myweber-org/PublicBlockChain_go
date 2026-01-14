
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

type RotatingLogger struct {
    mu         sync.Mutex
    file       *os.File
    filePath   string
    maxSize    int64
    maxBackups int
    currentSize int64
}

func NewRotatingLogger(filePath string, maxSize int64, maxBackups int) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        filePath:   filePath,
        maxSize:    maxSize,
        maxBackups: maxBackups,
    }

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
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
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
    if rl.file != nil {
        rl.file.Close()
    }

    for i := rl.maxBackups - 1; i >= 0; i-- {
        oldPath := rl.backupPath(i)
        newPath := rl.backupPath(i + 1)

        if _, err := os.Stat(oldPath); err == nil {
            if i == rl.maxBackups-1 {
                os.Remove(oldPath)
            } else {
                if err := rl.compressAndMove(oldPath, newPath); err != nil {
                    return err
                }
            }
        }
    }

    if err := rl.compressAndMove(rl.filePath, rl.backupPath(0)); err != nil {
        return err
    }

    return rl.openFile()
}

func (rl *RotatingLogger) backupPath(index int) string {
    if index == 0 {
        return rl.filePath + ".1"
    }
    return fmt.Sprintf("%s.%d.gz", rl.filePath, index)
}

func (rl *RotatingLogger) compressAndMove(src, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    if err != nil {
        return err
    }

    return os.Remove(src)
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
    logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(logger)

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
        time.Sleep(10 * time.Millisecond)
    }
}