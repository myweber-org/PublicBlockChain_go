package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logExtension = ".log"
)

type Rotator struct {
    filename   string
    current    *os.File
    size       int64
}

func NewRotator(filename string) (*Rotator, error) {
    r := &Rotator{filename: filename}
    if err := r.openFile(); err != nil {
        return nil, err
    }
    return r, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    if r.size+int64(len(p)) > maxFileSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := r.current.Write(p)
    r.size += int64(n)
    return n, err
}

func (r *Rotator) openFile() error {
    info, err := os.Stat(r.filename)
    if os.IsNotExist(err) {
        file, err := os.Create(r.filename)
        if err != nil {
            return err
        }
        r.current = file
        r.size = 0
        return nil
    }
    if err != nil {
        return err
    }

    file, err := os.OpenFile(r.filename, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    r.current = file
    r.size = info.Size()
    return nil
}

func (r *Rotator) rotate() error {
    if err := r.current.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s%s", strings.TrimSuffix(r.filename, logExtension), timestamp, logExtension)
    if err := os.Rename(r.filename, backupName); err != nil {
        return err
    }

    if err := r.openFile(); err != nil {
        return err
    }

    go r.cleanupOldLogs()
    return nil
}

func (r *Rotator) cleanupOldLogs() {
    dir := filepath.Dir(r.filename)
    base := strings.TrimSuffix(filepath.Base(r.filename), logExtension)

    files, err := filepath.Glob(filepath.Join(dir, base+".*"+logExtension))
    if err != nil {
        return
    }

    sort.Sort(sort.Reverse(sort.StringSlice(files)))
    for i := maxBackups; i < len(files); i++ {
        os.Remove(files[i])
    }
}

func (r *Rotator) Close() error {
    if r.current != nil {
        return r.current.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("app.log")
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }
}