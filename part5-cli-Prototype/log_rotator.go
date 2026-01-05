
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

	logger := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(l.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	l.currentFile = file
	l.currentSize = info.Size()
	return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.currentSize+int64(len(p)) > l.maxSize {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := l.currentFile.Write(p)
	if err == nil {
		l.currentSize += int64(n)
	}
	return n, err
}

func (l *RotatingLogger) rotate() error {
	if err := l.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", l.filePath, timestamp)

	if err := os.Rename(l.filePath, backupPath); err != nil {
		return err
	}

	if err := l.compressFile(backupPath); err != nil {
		return err
	}

	if err := l.cleanupOldBackups(); err != nil {
		return err
	}

	return l.openCurrentFile()
}

func (l *RotatingLogger) compressFile(sourcePath string) error {
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

	if err := os.Remove(sourcePath); err != nil {
		return err
	}

	return nil
}

func (l *RotatingLogger) cleanupOldBackups() error {
	pattern := l.filePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= l.backupCount {
		return nil
	}

	backups := make([]string, len(matches))
	copy(backups, matches)

	for i := 0; i < len(backups)-l.backupCount; i++ {
		if err := os.Remove(backups[i]); err != nil {
			return err
		}
	}

	return nil
}

func (l *RotatingLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.currentFile != nil {
		return l.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("logs/app.log", 10, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}