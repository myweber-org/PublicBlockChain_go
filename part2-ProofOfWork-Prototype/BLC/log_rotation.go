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
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    maxFiles    int
    currentFile *os.File
    currentSize int64
}

func NewRotatingLogger(basePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }

    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(l.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (l *RotatingLogger) rotate() error {
    l.currentFile.Close()

    for i := l.maxFiles - 1; i > 0; i-- {
        oldPath := l.basePath + "." + strconv.Itoa(i)
        newPath := l.basePath + "." + strconv.Itoa(i+1)

        if _, err := os.Stat(oldPath); err == nil {
            os.Rename(oldPath, newPath)
        }
    }

    if _, err := os.Stat(l.basePath); err == nil {
        os.Rename(l.basePath, l.basePath+".1")
    }

    return l.openCurrentFile()
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentSize+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := l.currentFile.Write(p)
    if err == nil {
        l.currentSize += int64(n)
    }
    return n, err
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
    logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}