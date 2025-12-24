package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    filePath     string
    maxSize      int64
    rotationTime time.Duration
    currentSize  int64
    lastRotation time.Time
    mu           sync.Mutex
    file         *os.File
}

func NewRotator(filePath string, maxSize int64, rotationTime time.Duration) (*Rotator, error) {
    r := &Rotator{
        filePath:     filePath,
        maxSize:      maxSize,
        rotationTime: rotationTime,
        lastRotation: time.Now(),
    }

    if err := r.openFile(); err != nil {
        return nil, err
    }

    return r, nil
}

func (r *Rotator) openFile() error {
    dir := filepath.Dir(r.filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    r.file = file
    r.currentSize = info.Size()
    return nil
}

func (r *Rotator) rotate() error {
    if r.file != nil {
        r.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)

    if err := os.Rename(r.filePath, backupPath); err != nil {
        return err
    }

    r.lastRotation = time.Now()
    return r.openFile()
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    now := time.Now()
    if r.currentSize+int64(len(p)) > r.maxSize || now.Sub(r.lastRotation) > r.rotationTime {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.file.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.file != nil {
        return r.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("logs/app.log", 1024*1024, 24*time.Hour)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        rotator.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
}
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    currentFile *os.File
    basePath    string
    maxSize     int64
    fileSize    int64
    fileCount   int
    maxFiles    int
}

func NewRotatingLogger(basePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
    if maxSize <= 0 {
        return nil, fmt.Errorf("maxSize must be positive")
    }
    if maxFiles <= 0 {
        return nil, fmt.Errorf("maxFiles must be positive")
    }

    rl := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.fileSize = info.Size()
    return nil
}

func (rl *RotatingLogger) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

    source, err := os.Open(rl.basePath)
    if err != nil {
        return err
    }
    defer source.Close()

    dest, err := os.Create(archivedPath)
    if err != nil {
        return err
    }
    defer dest.Close()

    gzWriter := gzip.NewWriter(dest)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, source); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    rl.fileCount++
    if rl.fileCount > rl.maxFiles {
        rl.cleanupOldFiles()
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLogger) cleanupOldFiles() {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > rl.maxFiles {
        filesToDelete := len(matches) - rl.maxFiles
        for i := 0; i < filesToDelete; i++ {
            os.Remove(matches[i])
        }
    }
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.fileSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
        rl.fileSize = 0
    }

    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.fileSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
}