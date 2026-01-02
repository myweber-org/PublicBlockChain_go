
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

type RotatingLog struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    maxFiles    int
    currentFile *os.File
    currentSize int64
}

func NewRotatingLog(basePath string, maxSizeMB int, maxFiles int) (*RotatingLog, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    rl := &RotatingLog{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLog) openCurrentFile() error {
    dir := filepath.Dir(rl.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLog) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    if err := rl.cleanupOldFiles(); err != nil {
        return err
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLog) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    compressedPath := sourcePath + ".gz"
    compressedFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer compressedFile.Close()

    gzWriter := gzip.NewWriter(compressedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (rl *RotatingLog) cleanupOldFiles() error {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var compressedFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            compressedFiles = append(compressedFiles, filepath.Join(dir, name))
        }
    }

    if len(compressedFiles) <= rl.maxFiles {
        return nil
    }

    sortByTimestamp := func(files []string) []string {
        type fileInfo struct {
            path      string
            timestamp string
        }

        var fileInfos []fileInfo
        for _, file := range files {
            parts := strings.Split(strings.TrimSuffix(filepath.Base(file), ".gz"), ".")
            if len(parts) >= 2 {
                fileInfos = append(fileInfos, fileInfo{
                    path:      file,
                    timestamp: parts[len(parts)-1],
                })
            }
        }

        for i := 0; i < len(fileInfos); i++ {
            for j := i + 1; j < len(fileInfos); j++ {
                if fileInfos[i].timestamp > fileInfos[j].timestamp {
                    fileInfos[i], fileInfos[j] = fileInfos[j], fileInfos[i]
                }
            }
        }

        sorted := make([]string, len(fileInfos))
        for i, info := range fileInfos {
            sorted[i] = info.path
        }
        return sorted
    }

    sortedFiles := sortByTimestamp(compressedFiles)
    filesToDelete := sortedFiles[:len(sortedFiles)-rl.maxFiles]

    for _, file := range filesToDelete {
        if err := os.Remove(file); err != nil {
            return err
        }
    }

    return nil
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("/var/log/myapp/app.log", 10, 5)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := log.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}