
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

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentNum int
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

	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}
	rl.size += int64(n)

	if err := rl.rotateIfNeeded(); err != nil {
		return n, err
	}
	return n, nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.size < maxFileSize && rl.file != nil {
		return nil
	}

	if rl.file != nil {
		rl.file.Close()
		if err := rl.compressCurrent(); err != nil {
			return err
		}
		rl.cleanOldBackups()
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	rl.file = file
	info, err := file.Stat()
	if err != nil {
		return err
	}
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) compressCurrent() error {
	src, err := os.Open(rl.basePath)
	if err != nil {
		return err
	}
	defer src.Close()

	backupName := fmt.Sprintf("%s.%d.gz", rl.basePath, rl.currentNum)
	dst, err := os.Create(backupName)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}
	rl.currentNum = (rl.currentNum + 1) % backupCount
	return os.Remove(rl.basePath)
}

func (rl *RotatingLogger) cleanOldBackups() {
	pattern := fmt.Sprintf("%s.*.gz", rl.basePath)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= backupCount {
		return
	}

	oldest := matches[0]
	for _, match := range matches[1:] {
		info1, _ := os.Stat(oldest)
		info2, _ := os.Stat(match)
		if info2.ModTime().Before(info1.ModTime()) {
			oldest = match
		}
	}
	os.Remove(oldest)
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
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(100 * time.Millisecond)
	}
}package main

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
    file         *os.File
    mu           sync.Mutex
}

func NewRotator(filePath string, maxSize int64, rotationTime time.Duration) (*Rotator, error) {
    dir := filepath.Dir(filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create directory: %w", err)
    }

    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("failed to stat log file: %w", err)
    }

    return &Rotator{
        filePath:     filePath,
        maxSize:      maxSize,
        rotationTime: rotationTime,
        currentSize:  info.Size(),
        lastRotation: time.Now(),
        file:         file,
    }, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if err := r.rotateIfNeeded(); err != nil {
        return 0, err
    }

    n, err := r.file.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) rotateIfNeeded() error {
    now := time.Now()
    shouldRotate := false

    if r.maxSize > 0 && r.currentSize >= r.maxSize {
        shouldRotate = true
    }

    if r.rotationTime > 0 && now.Sub(r.lastRotation) >= r.rotationTime {
        shouldRotate = true
    }

    if !shouldRotate {
        return nil
    }

    if err := r.rotate(); err != nil {
        return err
    }

    r.lastRotation = now
    return nil
}

func (r *Rotator) rotate() error {
    if err := r.file.Close(); err != nil {
        return fmt.Errorf("failed to close current file: %w", err)
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)

    if err := os.Rename(r.filePath, backupPath); err != nil {
        return fmt.Errorf("failed to rename log file: %w", err)
    }

    file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to create new log file: %w", err)
    }

    r.file = file
    r.currentSize = 0
    return nil
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.file.Close()
}

func main() {
    rotator, err := NewRotator("logs/app.log", 1024*1024, 24*time.Hour)
    if err != nil {
        fmt.Printf("Failed to create rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
        }
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
    "sync"
    "time"
)

const (
    maxFileSize   = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
)

type LogRotator struct {
    filePath    string
    currentFile *os.File
    currentSize int64
    mu          sync.Mutex
}

func NewLogRotator(path string) (*LogRotator, error) {
    lr := &LogRotator{
        filePath: path,
    }

    if err := lr.openCurrentFile(); err != nil {
        return nil, err
    }

    return lr, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

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

    if err := lr.compressOldFiles(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)
    if err := os.Rename(lr.filePath, rotatedPath); err != nil {
        return err
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = stat.Size()
    return nil
}

func (lr *LogRotator) compressOldFiles() error {
    dir := filepath.Dir(lr.filePath)
    baseName := filepath.Base(lr.filePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var backupFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && !strings.HasSuffix(name, ".gz") {
            backupFiles = append(backupFiles, name)
        }
    }

    if len(backupFiles) <= maxBackupFiles {
        return nil
    }

    sortBackupFiles(backupFiles)

    for i := 0; i < len(backupFiles)-maxBackupFiles; i++ {
        oldPath := filepath.Join(dir, backupFiles[i])
        if err := compressFile(oldPath); err != nil {
            return err
        }
        os.Remove(oldPath)
    }

    return nil
}

func sortBackupFiles(files []string) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            if extractTimestamp(files[i]) > extractTimestamp(files[j]) {
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
    timestampStr = strings.ReplaceAll(timestampStr, "_", "")
    ts, err := strconv.ParseInt(timestampStr, 10, 64)
    if err != nil {
        return 0
    }
    return ts
}

func compressFile(srcPath string) error {
    srcFile, err := os.Open(srcPath)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstPath := srcPath + ".gz"
    dstFile, err := os.Create(dstPath)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
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
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024
	backupCount = 5
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentSize int64
	basePath    string
	file        *os.File
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{basePath: path}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
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

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

	if err := os.Rename(rl.basePath, backupPath); err != nil {
		return err
	}

	if err := rl.compressBackup(backupPath); err != nil {
		return err
	}

	if err := rl.cleanOldBackups(); err != nil {
		return err
	}

	return rl.openCurrent()
}

func (rl *RotatingLogger) compressBackup(src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(src + ".gz")
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	os.Remove(src)
	return nil
}

func (rl *RotatingLogger) cleanOldBackups() error {
	pattern := rl.basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= backupCount {
		return nil
	}

	for i := 0; i < len(matches)-backupCount; i++ {
		os.Remove(matches[i])
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
	logger, err := NewRotatingLogger("app.log")
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