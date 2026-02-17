
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
	logDir        = "./logs"
	currentLog    = "app.log"
	compressOld   = true
)

type LogRotator struct {
	mu        sync.Mutex
	file      *os.File
	size      int64
	basePath  string
}

func NewLogRotator() (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	filePath := filepath.Join(logDir, currentLog)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &LogRotator{
		file:     file,
		size:     info.Size(),
		basePath: filePath,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	n, err := lr.file.Write(p)
	if err != nil {
		return n, err
	}

	lr.size += int64(n)
	if lr.size >= maxFileSize {
		if err := lr.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
	}

	return n, nil
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

	if err := os.Rename(lr.basePath, backupPath); err != nil {
		return err
	}

	file, err := os.OpenFile(lr.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.file = file
	lr.size = 0

	go lr.manageBackups(backupPath)

	return nil
}

func (lr *LogRotator) manageBackups(backupPath string) {
	dir := filepath.Dir(backupPath)
	baseName := filepath.Base(lr.basePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	var backups []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, baseName+".") && name != currentLog {
			backups = append(backups, filepath.Join(dir, name))
		}
	}

	if len(backups) > backupCount {
		toRemove := backups[:len(backups)-backupCount]
		for _, path := range toRemove {
			os.Remove(path)
		}
		backups = backups[len(backups)-backupCount:]
	}

	if compressOld {
		for _, path := range backups {
			if !strings.HasSuffix(path, ".gz") {
				lr.compressFile(path)
			}
		}
	}
}

func (lr *LogRotator) compressFile(path string) {
	// Compression implementation placeholder
	// In production, use compress/gzip
	fmt.Printf("Would compress: %s\n", path)
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.file.Close()
}

func main() {
	rotator, err := NewLogRotator()
	if err != nil {
		log.Fatal(err)
	}
	defer rotator.Close()

	log.SetOutput(rotator)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, strings.Repeat("X", 1024))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation demo completed")
}