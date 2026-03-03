package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "sync"
)

const (
    maxLogSize    = 1024 * 1024 // 1MB
    maxBackupLogs = 5
)

type RotatingLogger struct {
    mu       sync.Mutex
    filename string
    file     *os.File
    size     int64
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    rl := &RotatingLogger{filename: filename}
    if err := rl.openFile(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) openFile() error {
    file, err := os.OpenFile(rl.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

    for i := maxBackupLogs - 1; i >= 0; i-- {
        oldName := rl.backupName(i)
        newName := rl.backupName(i + 1)
        if _, err := os.Stat(oldName); err == nil {
            os.Rename(oldName, newName)
        }
    }

    os.Rename(rl.filename, rl.backupName(0))
    return rl.openFile()
}

func (rl *RotatingLogger) backupName(index int) string {
    if index == 0 {
        return rl.filename + ".1"
    }
    return fmt.Sprintf("%s.%d", rl.filename, index+1)
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.size+int64(len(p)) > maxLogSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
        rl.size = 0
    }

    n, err = rl.file.Write(p)
    rl.size += int64(n)
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
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: This is a test log message for rotation testing.", i)
    }

    fmt.Println("Log rotation test completed. Check app.log and backup files.")
}