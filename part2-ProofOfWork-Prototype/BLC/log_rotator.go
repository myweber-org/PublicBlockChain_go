package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

const (
    maxFileSize   = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
)

type RotatingLogger struct {
    filename    string
    currentSize int64
    file        *os.File
    mu          sync.Mutex
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        filename: filename,
    }

    if err := rl.openFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openFile() error {
    file, err := os.OpenFile(rl.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
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

    if rl.currentSize+int64(len(p)) > maxFileSize {
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
    if err := rl.file.Close(); err != nil {
        return err
    }

    if err := rl.compressOldFiles(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedName := fmt.Sprintf("%s.%s", rl.filename, timestamp)
    if err := os.Rename(rl.filename, rotatedName); err != nil {
        return err
    }

    if err := rl.openFile(); err != nil {
        return err
    }

    rl.currentSize = 0
    return nil
}

func (rl *RotatingLogger) compressOldFiles() error {
    dir := filepath.Dir(rl.filename)
    base := filepath.Base(rl.filename)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var backupFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, base+".") && !strings.HasSuffix(name, ".gz") {
            backupFiles = append(backupFiles, name)
        }
    }

    if len(backupFiles) <= maxBackupFiles {
        return nil
    }

    sortBackupFiles(backupFiles)

    for i := 0; i < len(backupFiles)-maxBackupFiles; i++ {
        oldPath := filepath.Join(dir, backupFiles[i])
        compressedPath := oldPath + ".gz"

        if err := compressFile(oldPath, compressedPath); err != nil {
            fmt.Printf("Failed to compress %s: %v\n", oldPath, err)
            continue
        }

        if err := os.Remove(oldPath); err != nil {
            fmt.Printf("Failed to remove %s: %v\n", oldPath, err)
        }
    }

    return nil
}

func compressFile(src, dst string) error {
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
    return err
}

func sortBackupFiles(files []string) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            timeI := extractTimestamp(files[i])
            timeJ := extractTimestamp(files[j])
            if timeI > timeJ {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func extractTimestamp(filename string) int64 {
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return 0
    }

    timestampStr := parts[len(parts)-1]
    timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
    if err != nil {
        return 0
    }
    return timestamp
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
    logger, err := NewRotatingLogger("application.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := logger.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
}