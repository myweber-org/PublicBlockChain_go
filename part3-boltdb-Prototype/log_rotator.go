package logrotator

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Rotator struct {
	mu            sync.Mutex
	currentFile   *os.File
	currentSize   int64
	maxSize       int64
	maxBackups    int
	compress      bool
	basePath      string
	rotationCount int
}

func NewRotator(basePath string, maxSize int64, maxBackups int, compress bool) (*Rotator, error) {
	if maxSize <= 0 {
		return nil, fmt.Errorf("maxSize must be positive")
	}

	r := &Rotator{
		maxSize:    maxSize,
		maxBackups: maxBackups,
		compress:   compress,
		basePath:   basePath,
	}

	if err := r.openCurrentFile(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentSize+int64(len(p)) > r.maxSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := r.currentFile.Write(p)
	if err == nil {
		r.currentSize += int64(n)
	}
	return n, err
}

func (r *Rotator) rotate() error {
	if r.currentFile != nil {
		r.currentFile.Close()
	}

	r.rotationCount++
	timestamp := time.Now().Format("20060102_150405")
	oldPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)

	if err := os.Rename(r.basePath, oldPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	if r.compress {
		if err := r.compressFile(oldPath); err != nil {
			return fmt.Errorf("failed to compress log file: %w", err)
		}
		oldPath = oldPath + ".gz"
	}

	if err := r.cleanupOldBackups(); err != nil {
		return fmt.Errorf("failed to cleanup old backups: %w", err)
	}

	return r.openCurrentFile()
}

func (r *Rotator) compressFile(src string) error {
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

	return os.Remove(src)
}

func (r *Rotator) cleanupOldBackups() error {
	if r.maxBackups <= 0 {
		return nil
	}

	pattern := r.basePath + ".*"
	if r.compress {
		pattern += ".gz"
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= r.maxBackups {
		return nil
	}

	toDelete := matches[:len(matches)-r.maxBackups]
	for _, file := range toDelete {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func (r *Rotator) openCurrentFile() error {
	file, err := os.OpenFile(r.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	r.currentFile = file
	r.currentSize = info.Size()
	return nil
}

func (r *Rotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentFile != nil {
		return r.currentFile.Close()
	}
	return nil
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackupCount = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    if err := os.MkdirAll(basePath, 0755); err != nil {
        return nil, err
    }

    logPath := filepath.Join(basePath, logFileName)
    file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &LogRotator{
        currentFile: file,
        currentSize: info.Size(),
        basePath:    basePath,
    }, nil
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
    if err := lr.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    oldPath := filepath.Join(lr.basePath, logFileName)
    newPath := filepath.Join(lr.basePath, backupName)

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.currentFile = file
    lr.currentSize = 0

    go lr.cleanupOldBackups()

    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    pattern := filepath.Join(lr.basePath, logFileName+".*")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    sort.Sort(sort.Reverse(sort.StringSlice(matches)))

    for i, match := range matches {
        if i >= maxBackupCount {
            os.Remove(match)
        }
    }
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("./logs")
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
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

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	backupCount int
}

func NewRotatingLogger(filePath string, maxSizeMB int, backupCount int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingLogger{
		currentFile: file,
		filePath:    filePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
		backupCount: backupCount,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	backupDir := filepath.Dir(rl.filePath)
	baseName := filepath.Base(rl.filePath)
	timestamp := time.Now().Format("20060102_150405")

	for i := rl.backupCount - 1; i > 0; i-- {
		oldName := fmt.Sprintf("%s.%d.gz", baseName, i)
		newName := fmt.Sprintf("%s.%d.gz", baseName, i+1)
		oldPath := filepath.Join(backupDir, oldName)
		newPath := filepath.Join(backupDir, newName)

		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	rotatedFile := filepath.Join(backupDir, fmt.Sprintf("%s.%s", baseName, timestamp))
	if err := os.Rename(rl.filePath, rotatedFile); err != nil {
		return err
	}

	if err := rl.compressFile(rotatedFile); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0
	return nil
}

func (rl *RotatingLogger) compressFile(source string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	compressed := source + ".gz"
	writer, err := os.Create(compressed)
	if err != nil {
		return err
	}
	defer writer.Close()

	gzWriter := gzip.NewWriter(writer)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, reader); err != nil {
		return err
	}

	os.Remove(source)
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
}