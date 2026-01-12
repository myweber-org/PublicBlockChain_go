
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logExtension = ".log"
)

type LogRotator struct {
    filename    string
    currentSize int64
}

func NewLogRotator(filename string) (*LogRotator, error) {
    rotator := &LogRotator{filename: filename}
    
    if err := rotator.ensureLogFile(); err != nil {
        return nil, err
    }
    
    if err := rotator.loadCurrentSize(); err != nil {
        return nil, err
    }
    
    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    
    file, err := os.OpenFile(lr.filename, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer file.Close()
    
    n, err := file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.cleanupOldBackups(); err != nil {
        return err
    }
    
    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s%s", strings.TrimSuffix(lr.filename, logExtension), timestamp, logExtension)
    
    if err := os.Rename(lr.filename, backupName); err != nil {
        return err
    }
    
    lr.currentSize = 0
    return lr.ensureLogFile()
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := fmt.Sprintf("%s.*%s", strings.TrimSuffix(lr.filename, logExtension), logExtension)
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }
    
    if len(matches) <= maxBackups {
        return nil
    }
    
    sort.Strings(matches)
    filesToRemove := matches[:len(matches)-maxBackups]
    
    for _, file := range filesToRemove {
        if err := os.Remove(file); err != nil {
            return err
        }
    }
    
    return nil
}

func (lr *LogRotator) ensureLogFile() error {
    file, err := os.OpenFile(lr.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    return file.Close()
}

func (lr *LogRotator) loadCurrentSize() error {
    info, err := os.Stat(lr.filename)
    if err != nil {
        if os.IsNotExist(err) {
            lr.currentSize = 0
            return nil
        }
        return err
    }
    lr.currentSize = info.Size()
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    
    for i := 1; i <= 100; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }
        
        if i%10 == 0 {
            fmt.Printf("Written %d log entries\n", i)
        }
        
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type RotatingWriter struct {
    mu          sync.Mutex
    file        *os.File
    currentSize int64
    maxSize     int64
    basePath    string
    fileIndex   int
}

func NewRotatingWriter(basePath string, maxSize int64) (*RotatingWriter, error) {
    w := &RotatingWriter{
        maxSize:  maxSize,
        basePath: basePath,
        fileIndex: 0,
    }
    
    if err := w.openFile(); err != nil {
        return nil, err
    }
    
    return w, nil
}

func (w *RotatingWriter) openFile() error {
    filename := fmt.Sprintf("%s.%d", w.basePath, w.fileIndex)
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    if w.file != nil {
        w.file.Close()
    }
    
    w.file = file
    w.currentSize = info.Size()
    return nil
}

func (w *RotatingWriter) rotateIfNeeded() error {
    if w.currentSize < w.maxSize {
        return nil
    }
    
    w.fileIndex++
    return w.openFile()
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    if err := w.rotateIfNeeded(); err != nil {
        return 0, err
    }
    
    n, err := w.file.Write(p)
    if err == nil {
        w.currentSize += int64(n)
    }
    return n, err
}

func (w *RotatingWriter) Close() error {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    if w.file != nil {
        return w.file.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("app.log", 1024*1024) // 1MB max size
    if err != nil {
        fmt.Printf("Failed to create rotating writer: %v\n", err)
        return
    }
    defer writer.Close()
    
    for i := 0; i < 100; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        writer.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
    logDir      = "./logs"
)

type RotatingWriter struct {
    currentFile *os.File
    currentSize int64
    baseName    string
    sequence    int
}

func NewRotatingWriter(baseName string) (*RotatingWriter, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    w := &RotatingWriter{
        baseName: baseName,
        sequence: 0,
    }

    if err := w.openNewFile(); err != nil {
        return nil, err
    }

    return w, nil
}

func (w *RotatingWriter) openNewFile() error {
    if w.currentFile != nil {
        w.currentFile.Close()
    }

    filename := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", w.baseName, w.sequence))
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    w.currentFile = file
    w.currentSize = info.Size()
    return nil
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
    if w.currentSize+int64(len(p)) > maxFileSize {
        w.sequence++
        if w.sequence > maxBackups {
            w.sequence = 0
        }
        if err := w.openNewFile(); err != nil {
            return 0, err
        }
    }

    n, err = w.currentFile.Write(p)
    w.currentSize += int64(n)
    return n, err
}

func (w *RotatingWriter) Close() error {
    if w.currentFile != nil {
        return w.currentFile.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("app")
    if err != nil {
        fmt.Printf("Failed to create writer: %v\n", err)
        return
    }
    defer writer.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := writer.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}