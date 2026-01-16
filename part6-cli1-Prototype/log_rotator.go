package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	basePath    string
	maxSize     int64
	currentSize int64
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	logger := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
	}

	if err := logger.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.currentFile == nil || rl.currentSize >= rl.maxSize {
		return rl.rotate()
	}
	return nil
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	newPath := fmt.Sprintf("%s_%s.log", rl.basePath, timestamp)

	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0

	// Clean old files (keep last 10)
	go rl.cleanOldFiles()

	return nil
}

func (rl *RotatingLogger) cleanOldFiles() {
	files, err := filepath.Glob(rl.basePath + "_*.log")
	if err != nil {
		return
	}

	if len(files) > 10 {
		oldestFiles := files[:len(files)-10]
		for _, file := range oldestFiles {
			os.Remove(file)
		}
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}