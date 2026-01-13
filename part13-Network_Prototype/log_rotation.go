
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
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
)

type RotatingLogger struct {
    filename   string
    current    *os.File
    fileSize   int64
    backupTime time.Time
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
    file, err := os.OpenFile(rl.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    rl.current = file
    rl.fileSize = info.Size()
    rl.backupTime = time.Now()
    
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.shouldRotate(len(p)) {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err = rl.current.Write(p)
    if err == nil {
        rl.fileSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) shouldRotate(addSize int) bool {
    if rl.fileSize+int64(addSize) > maxFileSize {
        return true
    }
    
    if time.Since(rl.backupTime) > 24*time.Hour {
        return true
    }
    
    return false
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.current.Close(); err != nil {
        return err
    }
    
    backupName := rl.generateBackupName()
    if err := os.Rename(rl.filename, backupName); err != nil {
        return err
    }
    
    rl.cleanOldBackups()
    
    return rl.openFile()
}

func (rl *RotatingLogger) generateBackupName() string {
    timestamp := time.Now().Format("20060102_150405")
    ext := filepath.Ext(rl.filename)
    base := strings.TrimSuffix(rl.filename, ext)
    return fmt.Sprintf("%s_%s%s", base, timestamp, ext)
}

func (rl *RotatingLogger) cleanOldBackups() {
    pattern := strings.TrimSuffix(rl.filename, filepath.Ext(rl.filename)) + "_*" + filepath.Ext(rl.filename)
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    if len(matches) <= maxBackups {
        return
    }
    
    backups := make([]backupInfo, 0, len(matches))
    for _, match := range matches {
        if info, err := parseBackupInfo(match); err == nil {
            backups = append(backups, info)
        }
    }
    
    sortBackups(backups)
    
    for i := maxBackups; i < len(backups); i++ {
        os.Remove(backups[i].path)
    }
}

type backupInfo struct {
    path string
    time time.Time
}

func parseBackupInfo(path string) (backupInfo, error) {
    base := filepath.Base(path)
    ext := filepath.Ext(path)
    name := strings.TrimSuffix(base, ext)
    
    parts := strings.Split(name, "_")
    if len(parts) < 3 {
        return backupInfo{}, fmt.Errorf("invalid backup name")
    }
    
    timestamp := parts[len(parts)-2] + "_" + parts[len(parts)-1]
    t, err := time.Parse("20060102_150405", timestamp)
    if err != nil {
        return backupInfo{}, err
    }
    
    return backupInfo{path: path, time: t}, nil
}

func sortBackups(backups []backupInfo) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            if backups[i].time.After(backups[j].time) {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func (rl *RotatingLogger) Close() error {
    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()
    
    log.SetOutput(logger)
    
    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: %s", i, strings.Repeat("x", 1000))
        time.Sleep(100 * time.Millisecond)
    }
}