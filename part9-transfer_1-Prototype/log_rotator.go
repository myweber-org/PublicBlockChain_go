
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	maxFileSize    = 10 * 1024 * 1024 // 10MB
	backupCount    = 5
	checkInterval  = 30 * time.Second
	compressBackup = true
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	filePath   string
	currentSize int64
	baseName   string
	dir        string
}

func NewRotatingLogger(filePath string) (*RotatingLogger, error) {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
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

	rl := &RotatingLogger{
		file:       file,
		filePath:   filePath,
		currentSize: info.Size(),
		baseName:   base,
		dir:        dir,
	}

	go rl.monitor()
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}

	rl.currentSize += int64(n)
	if rl.currentSize >= maxFileSize {
		if err := rl.rotate(); err != nil {
			log.Printf("rotation failed: %v", err)
		}
	}
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(rl.dir, fmt.Sprintf("%s.%s", rl.baseName, timestamp))

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	if compressBackup {
		go compressFile(backupPath)
	}

	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	rl.file = file
	rl.currentSize = 0
	rl.cleanupOldBackups()
	return nil
}

func compressFile(path string) {
	// Compression implementation would go here
	// For simplicity, just log the intention
	log.Printf("Compressing backup file: %s", path)
	// In real implementation, use compress/gzip or similar
}

func (rl *RotatingLogger) cleanupOldBackups() {
	pattern := filepath.Join(rl.dir, fmt.Sprintf("%s.*", rl.baseName))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("failed to list backup files: %v", err)
		return
	}

	if len(matches) <= backupCount {
		return
	}

	// Sort by modification time (oldest first)
	// For simplicity, just remove excess files
	for i := 0; i < len(matches)-backupCount; i++ {
		if err := os.Remove(matches[i]); err != nil {
			log.Printf("failed to remove old backup %s: %v", matches[i], err)
		} else {
			log.Printf("removed old backup: %s", matches[i])
		}
	}
}

func (rl *RotatingLogger) monitor() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		if rl.currentSize >= maxFileSize {
			if err := rl.rotate(); err != nil {
				log.Printf("scheduled rotation failed: %v", err)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("./logs/application.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	// Redirect standard log output
	log.SetOutput(io.MultiWriter(os.Stdout, logger))

	// Simulate log writing
	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, strings.Repeat("x", 1024))
		time.Sleep(100 * time.Millisecond)
	}
}