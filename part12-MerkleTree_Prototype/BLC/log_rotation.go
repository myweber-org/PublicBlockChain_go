
package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxLogSize    = 5 * 1024 * 1024 // 5MB
    maxBackupFiles = 10
    logFileName   = "app.log"
)

type RotatingLogger struct {
    currentFile *os.File
    filePath    string
    baseName    string
    dir         string
    sequence    int
}

func NewRotatingLogger(dir string) (*RotatingLogger, error) {
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

    filePath := filepath.Join(dir, logFileName)
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }

    rl := &RotatingLogger{
        currentFile: file,
        filePath:    filePath,
        baseName:    strings.TrimSuffix(logFileName, filepath.Ext(logFileName)),
        dir:         dir,
        sequence:    0,
    }

    rl.initializeSequence()
    return rl, nil
}

func (rl *RotatingLogger) initializeSequence() {
    files, err := os.ReadDir(rl.dir)
    if err != nil {
        return
    }

    for _, file := range files {
        if strings.HasPrefix(file.Name(), rl.baseName) && strings.HasSuffix(file.Name(), ".log") {
            name := file.Name()
            if name == logFileName {
                continue
            }

            parts := strings.Split(name, ".")
            if len(parts) >= 3 {
                if seq, err := strconv.Atoi(parts[1]); err == nil && seq > rl.sequence {
                    rl.sequence = seq
                }
            }
        }
    }
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.shouldRotate() {
        if err := rl.rotate(); err != nil {
            log.Printf("Failed to rotate log file: %v", err)
        }
    }
    return rl.currentFile.Write(p)
}

func (rl *RotatingLogger) shouldRotate() bool {
    info, err := rl.currentFile.Stat()
    if err != nil {
        return false
    }
    return info.Size() >= maxLogSize
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.currentFile.Close(); err != nil {
        return fmt.Errorf("failed to close current log file: %w", err)
    }

    rl.sequence++
    backupName := fmt.Sprintf("%s.%d.log", rl.baseName, rl.sequence)
    backupPath := filepath.Join(rl.dir, backupName)

    if err := os.Rename(rl.filePath, backupPath); err != nil {
        return fmt.Errorf("failed to rename log file: %w", err)
    }

    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to create new log file: %w", err)
    }

    rl.currentFile = file
    rl.cleanupOldFiles()
    return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
    files, err := os.ReadDir(rl.dir)
    if err != nil {
        return
    }

    var logFiles []string
    for _, file := range files {
        if strings.HasPrefix(file.Name(), rl.baseName) && strings.HasSuffix(file.Name(), ".log") && file.Name() != logFileName {
            logFiles = append(logFiles, filepath.Join(rl.dir, file.Name()))
        }
    }

    if len(logFiles) <= maxBackupFiles {
        return
    }

    for i := 0; i < len(logFiles)-maxBackupFiles; i++ {
        os.Remove(logFiles[i])
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger("./logs")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
        time.Sleep(100 * time.Millisecond)
    }
}