package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	fullPath := filepath.Join(basePath, logFileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
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
			return 0, fmt.Errorf("rotation failed: %w", err)
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
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
	oldPath := filepath.Join(lr.basePath, logFileName)
	newPath := filepath.Join(lr.basePath, backupName)

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
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

	if len(matches) <= maxBackupCount {
		return
	}

	sort.Strings(matches)
	filesToRemove := matches[:len(matches)-maxBackupCount]

	for _, file := range filesToRemove {
		os.Remove(file)
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
		fmt.Fprintf(os.Stderr, "Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
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
    sequence    int
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
        sequence: 0,
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
        if err := lr.compressCurrentFile(); err != nil {
            return err
        }
        lr.cleanupOldFiles()
    }

    lr.sequence++
    return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
    filename := fmt.Sprintf("%s.%d.log", lr.basePath, lr.sequence)
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (lr *LogRotator) compressCurrentFile() error {
    oldName := fmt.Sprintf("%s.%d.log", lr.basePath, lr.sequence-1)
    newName := oldName + ".gz"

    src, err := os.Open(oldName)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(newName)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    os.Remove(oldName)
    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    pattern := fmt.Sprintf("%s.*.log.gz", lr.basePath)
    files, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(files) > maxBackups {
        filesToDelete := files[:len(files)-maxBackups]
        for _, file := range filesToDelete {
            os.Remove(file)
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
    rotator, err := NewLogRotator("application")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: This is a sample log message\n",
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
    }

    fmt.Println("Log rotation completed")
}