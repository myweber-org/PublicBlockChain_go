package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024
	logDir      = "./logs"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewLogRotator(baseName string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	lr := &LogRotator{
		baseName: baseName,
	}

	if err := lr.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return lr, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if err := lr.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := lr.currentFile.Write(p)
	if err != nil {
		return n, err
	}

	lr.currentSize += int64(n)
	return n, nil
}

func (lr *LogRotator) rotateIfNeeded() error {
	if lr.currentFile == nil || lr.currentSize >= maxFileSize {
		if lr.currentFile != nil {
			lr.currentFile.Close()
		}

		timestamp := time.Now().Format("20060102_150405")
		filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", lr.baseName, timestamp))

		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		lr.currentFile = file
		lr.currentSize = 0

		fmt.Printf("Rotated to new log file: %s\n", filename)
	}
	return nil
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
	rotator, err := NewLogRotator("app")
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}