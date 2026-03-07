
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentDay string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	now := time.Now()
	currentDay := now.Format("2006-01-02")

	if rl.file == nil || rl.size >= maxFileSize || rl.currentDay != currentDay {
		if rl.file != nil {
			rl.file.Close()
			if err := rl.compressOldLog(); err != nil {
				log.Printf("Failed to compress old log: %v", err)
			}
			if err := rl.cleanupOldBackups(); err != nil {
				log.Printf("Failed to cleanup old backups: %v", err)
			}
		}

		newPath := rl.getLogPath(now)
		if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
			return fmt.Errorf("create directory failed: %w", err)
		}

		file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("open log file failed: %w", err)
		}

		info, err := file.Stat()
		if err != nil {
			file.Close()
			return fmt.Errorf("stat log file failed: %w", err)
		}

		rl.file = file
		rl.size = info.Size()
		rl.currentDay = currentDay
	}
	return nil
}

func (rl *RotatingLogger) getLogPath(t time.Time) string {
	base := filepath.Base(rl.basePath)
	dir := filepath.Dir(rl.basePath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", name, t.Format("2006-01-02"), ext))
}

func (rl *RotatingLogger) compressOldLog() error {
	oldPath := rl.getLogPath(time.Now().AddDate(0, 0, -1))
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return nil
	}

	compressedPath := oldPath + ".gz"
	inFile, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, inFile); err != nil {
		return err
	}

	return os.Remove(oldPath)
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	dir := filepath.Dir(rl.basePath)
	base := filepath.Base(rl.basePath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	pattern := filepath.Join(dir, name+"-*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > backupCount {
		toDelete := matches[:len(matches)-backupCount]
		for _, file := range toDelete {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}
	return nil
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
	logger, err := NewRotatingLogger("/var/log/myapp/application.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: Application is running smoothly", i)
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
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
    }

    if err := rotator.openCurrentFile(); err != nil {
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

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    lr.cleanupOldBackups()

    return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(source string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(source + ".gz")
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    if err != nil {
        return err
    }

    os.Remove(source)
    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    backupTimes := make([]time.Time, 0, len(matches))
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        timestamp := parts[len(parts)-2]
        t, err := time.Parse("20060102150405", timestamp)
        if err != nil {
            continue
        }
        backupTimes = append(backupTimes, t)
    }

    for i := 0; i < len(backupTimes)-maxBackups; i++ {
        oldestIdx := 0
        for j := 1; j < len(backupTimes); j++ {
            if backupTimes[j].Before(backupTimes[oldestIdx]) {
                oldestIdx = j
            }
        }
        os.Remove(matches[oldestIdx])
        backupTimes = append(backupTimes[:oldestIdx], backupTimes[oldestIdx+1:]...)
        matches = append(matches[:oldestIdx], matches[oldestIdx+1:]...)
    }
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        _, err := rotator.Write([]byte(logEntry))
        if err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
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
    currentFile *os.File
    currentSize int64
    maxFiles    int
}

func NewRotatingLog(basePath string, maxSizeMB int, maxFiles int) (*RotatingLog, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    if maxSize <= 0 {
        return nil, fmt.Errorf("maxSize must be positive")
    }

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
    rotatedPath := rl.basePath + "." + timestamp

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
    destFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (rl *RotatingLog) cleanupOldFiles() error {
    if rl.maxFiles <= 0 {
        return nil
    }

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
            compressedFiles = append(compressedFiles, name)
        }
    }

    if len(compressedFiles) <= rl.maxFiles {
        return nil
    }

    sortByTimestamp := func(files []string) []string {
        var result []string
        for _, f := range files {
            parts := strings.Split(f, ".")
            if len(parts) >= 3 {
                timestamp := parts[len(parts)-2]
                if _, err := time.Parse("20060102_150405", timestamp); err == nil {
                    result = append(result, f)
                }
            }
        }
        return result
    }

    sortedFiles := sortByTimestamp(compressedFiles)
    filesToDelete := sortedFiles[:len(sortedFiles)-rl.maxFiles]

    for _, file := range filesToDelete {
        os.Remove(filepath.Join(dir, file))
    }

    return nil
}

func (rl *RotatingLog) extractTimestamp(filename string) (time.Time, error) {
    parts := strings.Split(filename, ".")
    if len(parts) < 3 {
        return time.Time{}, fmt.Errorf("invalid filename format")
    }

    timestampStr := parts[len(parts)-2]
    return time.Parse("20060102_150405", timestampStr)
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
    logFile := "application.log"
    rotator, err := NewRotatingLog(logFile, 10, 5)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}