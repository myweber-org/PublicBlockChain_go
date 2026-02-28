package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
)

type RotatingLogger struct {
    filename   string
    current   *os.File
    size      int64
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    stat, err := f.Stat()
    if err != nil {
        f.Close()
        return nil, err
    }

    return &RotatingLogger{
        filename: filename,
        current: f,
        size:    stat.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.size+int64(len(p)) > maxLogSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.current.Write(p)
    rl.size += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.current.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s.gz", rl.filename, timestamp)
    
    if err := compressFile(rl.filename, backupName); err != nil {
        return err
    }

    if err := os.Remove(rl.filename); err != nil {
        return err
    }

    f, err := os.OpenFile(rl.filename, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.current = f
    rl.size = 0

    return cleanupOldBackups(rl.filename)
}

func compressFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    gz := gzip.NewWriter(out)
    defer gz.Close()

    _, err = io.Copy(gz, in)
    return err
}

func cleanupOldBackups(baseFilename string) error {
    pattern := baseFilename + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackupFiles {
        return nil
    }

    for i := 0; i < len(matches)-maxBackupFiles; i++ {
        if err := os.Remove(matches[i]); err != nil {
            return err
        }
    }

    return nil
}

func (rl *RotatingLogger) Close() error {
    return rl.current.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logger.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
}package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize  = 1024 * 1024 // 1MB
    maxBackups   = 5
    logDirectory = "./logs"
)

type RotatingLogger struct {
    currentFile *os.File
    fileSize    int64
    baseName    string
    sequence    int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
    if err := os.MkdirAll(logDirectory, 0755); err != nil {
        return nil, err
    }

    rl := &RotatingLogger{
        baseName: baseName,
        sequence: 0,
    }

    if err := rl.openNextFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openNextFile() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    rl.sequence++
    if rl.sequence > maxBackups {
        rl.sequence = 1
    }

    filename := filepath.Join(logDirectory, fmt.Sprintf("%s.%d.log", rl.baseName, rl.sequence))
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.fileSize = 0

    rl.cleanOldFiles()
    return nil
}

func (rl *RotatingLogger) cleanOldFiles() {
    files, err := filepath.Glob(filepath.Join(logDirectory, rl.baseName+".*.log"))
    if err != nil {
        return
    }

    if len(files) > maxBackups {
        for _, file := range files[:len(files)-maxBackups] {
            os.Remove(file)
        }
    }
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.fileSize+int64(len(p)) > maxFileSize {
        if err := rl.openNextFile(); err != nil {
            return 0, err
        }
    }

    n, err = rl.currentFile.Write(p)
    rl.fileSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    multiWriter := io.MultiWriter(os.Stdout, logger)
    log.SetOutput(multiWriter)

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
        time.Sleep(10 * time.Millisecond)
    }
}