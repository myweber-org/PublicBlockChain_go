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
    if err := r.openCurrent(); err != nil {
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

func (r *Rotator) rotate() error {
    if r.current != nil {
        r.current.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedName := fmt.Sprintf("%s.%s%s", strings.TrimSuffix(r.filename, logExtension), timestamp, logExtension)
    if err := os.Rename(r.filename, rotatedName); err != nil {
        return err
    }

    if err := r.openCurrent(); err != nil {
        return err
    }

    go r.cleanup()
    return nil
}

func (r *Rotator) openCurrent() error {
    f, err := os.OpenFile(r.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    stat, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    r.current = f
    r.size = stat.Size()
    return nil
}

func (r *Rotator) cleanup() {
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
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        rotator.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
}package main

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
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    mu         sync.Mutex
    file       *os.File
    basePath   string
    currentSize int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        file:       file,
        basePath:   path,
        currentSize: info.Size(),
    }, nil
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

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    rl.cleanupOldBackups()

    return nil
}

func (rl *RotatingLogger) compressFile(path string) error {
    src, err := os.Open(path)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(path + ".gz")
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    os.Remove(path)
    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    var backups []string
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        timestamp := parts[len(parts)-2]
        if _, err := time.Parse("20060102_150405", timestamp); err == nil {
            backups = append(backups, match)
        }
    }

    if len(backups) > maxBackups {
        toDelete := backups[:len(backups)-maxBackups]
        for _, file := range toDelete {
            os.Remove(file)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("/var/log/myapp/app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}