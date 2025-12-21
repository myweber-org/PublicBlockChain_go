package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
)

const maxFileSize = 1024 * 1024 // 1MB
const backupCount = 5

type RotatingLogger struct {
    filename   string
    current   *os.File
    size      int64
    mu        sync.Mutex
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    rl := &RotatingLogger{filename: filename}
    if err := rl.openCurrent(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
    file, err := os.OpenFile(rl.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    rl.current = file
    rl.size = info.Size()
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.current.Write(p)
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.current.Close(); err != nil {
        return err
    }

    for i := backupCount - 1; i >= 0; i-- {
        oldName := rl.backupName(i)
        newName := rl.backupName(i + 1)
        if _, err := os.Stat(oldName); err == nil {
            os.Rename(oldName, newName)
        }
    }

    os.Rename(rl.filename, rl.backupName(0))
    return rl.openCurrent()
}

func (rl *RotatingLogger) backupName(index int) string {
    if index == 0 {
        return rl.filename + ".1"
    }
    return fmt.Sprintf("%s.%d", rl.filename, index+1)
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 10000; i++ {
        logger.Write([]byte(fmt.Sprintf("Log entry %d: Some sample log data for testing rotation mechanism\n", i)))
    }
    fmt.Println("Log rotation test completed")
}