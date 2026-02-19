
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
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type LogRotator struct {
	mu         sync.Mutex
	filePath   string
	maxSize    int64
	currentSize int64
	file       *os.File
}

func NewLogRotator(filePath string, maxSize int64) (*LogRotator, error) {
	lr := &LogRotator{
		filePath: filePath,
		maxSize:  maxSize,
	}

	if err := lr.openFile(); err != nil {
		return nil, err
	}

	return lr, nil
}

func (lr *LogRotator) openFile() error {
	info, err := os.Stat(lr.filePath)
	if err == nil {
		lr.currentSize = info.Size()
	} else if os.IsNotExist(err) {
		lr.currentSize = 0
	} else {
		return err
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	lr.file = file
	return nil
}

func (lr *LogRotator) rotate() error {
	if lr.file != nil {
		lr.file.Close()
	}

	dir := filepath.Dir(lr.filePath)
	base := filepath.Base(lr.filePath)
	backupPath := filepath.Join(dir, fmt.Sprintf("%s.1", base))

	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return err
	}

	if err := lr.openFile(); err != nil {
		return err
	}

	lr.currentSize = 0
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.file.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app.log", 1024*1024) // 1MB max size
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: This is a sample log message.\n", i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
	}

	fmt.Println("Log rotation test completed")
}