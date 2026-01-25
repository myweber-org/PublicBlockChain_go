
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
	maxFileSize   = 10 * 1024 * 1024 // 10MB
	backupCount   = 5
	checkInterval = 30 * time.Second
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	filePath   string
	currentPos int64
	closed     bool
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath: path,
	}
	if err := rl.openFile(); err != nil {
		return nil, err
	}
	go rl.monitor()
	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.currentPos = stat.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.closed {
		return 0, io.ErrClosedPipe
	}

	n, err = rl.file.Write(p)
	if err != nil {
		return n, err
	}

	rl.currentPos += int64(n)
	if rl.currentPos >= maxFileSize {
		go rl.rotate()
	}
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.closed {
		return nil
	}

	if err := rl.file.Close(); err != nil {
		return err
	}

	// Rename current log file
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)
	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}

	// Reopen log file
	if err := rl.openFile(); err != nil {
		return err
	}

	// Clean old backups
	go rl.cleanOldBackups()
	return nil
}

func (rl *RotatingLogger) cleanOldBackups() {
	dir := filepath.Dir(rl.filePath)
	baseName := filepath.Base(rl.filePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	var backups []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, baseName+".") && entry.Type().IsRegular() {
			backups = append(backups, filepath.Join(dir, name))
		}
	}

	if len(backups) <= backupCount {
		return
	}

	// Sort by modification time (oldest first)
	for i := 0; i < len(backups)-backupCount; i++ {
		os.Remove(backups[i])
	}
}

func (rl *RotatingLogger) monitor() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		closed := rl.closed
		rl.mu.Unlock()

		if closed {
			return
		}

		stat, err := os.Stat(rl.filePath)
		if err == nil && stat.Size() >= maxFileSize {
			rl.rotate()
		}
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.closed = true
	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("logs/app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}
}