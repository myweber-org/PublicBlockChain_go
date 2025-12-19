
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    backupCount = 5
)

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
    if rl.current != nil {
        rl.current.Close()
    }

    for i := backupCount - 1; i >= 0; i-- {
        oldName := rl.backupName(i)
        newName := rl.backupName(i + 1)
        if _, err := os.Stat(oldName); err == nil {
            os.Rename(oldName, newName)
        }
    }

    if err := os.Rename(rl.filename, rl.backupName(0)); err != nil && !os.IsNotExist(err) {
        return err
    }

    return rl.openCurrent()
}

func (rl *RotatingLogger) backupName(i int) string {
    if i == 0 {
        return rl.filename + ".1"
    }
    return fmt.Sprintf("%s.%d", rl.filename, i+1)
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

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}