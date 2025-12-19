
package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize    = 1024 * 1024 // 1MB
    maxBackupFiles = 5
    logFileName   = "app.log"
)

type RotatingWriter struct {
    currentSize int64
    file        *os.File
    basePath    string
}

func NewRotatingWriter(path string) (*RotatingWriter, error) {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingWriter{
        currentSize: stat.Size(),
        file:        file,
        basePath:    path,
    }, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    if w.currentSize+int64(len(p)) > maxLogSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := w.file.Write(p)
    if err == nil {
        w.currentSize += int64(n)
    }
    return n, err
}

func (w *RotatingWriter) rotate() error {
    if err := w.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", w.basePath, timestamp)
    if err := os.Rename(w.basePath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(w.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    w.file = file
    w.currentSize = 0

    go w.cleanupOldLogs()
    return nil
}

func (w *RotatingWriter) cleanupOldLogs() {
    dir := filepath.Dir(w.basePath)
    baseName := filepath.Base(w.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        matched, _ := filepath.Match(baseName+".*", name)
        if matched && entry.Type().IsRegular() {
            backups = append(backups, filepath.Join(dir, name))
        }
    }

    if len(backups) <= maxBackupFiles {
        return
    }

    for i := 0; i < len(backups)-maxBackupFiles; i++ {
        os.Remove(backups[i])
    }
}

func (w *RotatingWriter) Close() error {
    return w.file.Close()
}

func main() {
    writer, err := NewRotatingWriter("./logs/app.log")
    if err != nil {
        log.Fatal(err)
    }
    defer writer.Close()

    logger := log.New(writer, "", log.LstdFlags)

    for i := 0; i < 1000; i++ {
        logger.Printf("Log entry %d: Application is running normally", i)
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation demonstration completed")
}