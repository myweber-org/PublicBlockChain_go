
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	fileIndex    int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}
	
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}
	
	filename := fmt.Sprintf("%s.%d.log", rl.basePath, rl.fileIndex)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	
	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.currentSize < rl.maxSize {
		return nil
	}
	
	rl.fileIndex++
	return rl.openCurrentFile()
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

func (rl *RotatingLogger) Log(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, message)
	rl.Write([]byte(logEntry))
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("application", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()
	
	for i := 0; i < 1000; i++ {
		logger.Log("INFO", fmt.Sprintf("Processing item %d", i))
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Println("Log rotation example completed")
}
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
	// For simplicity, just log the compression attempt
	log.Printf("Compressing backup file: %s", path)
	// Actual compression would use gzip or similar
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

	// Redirect standard log output to our rotating logger
	log.SetOutput(logger)

	// Simulate log writing
	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, strings.Repeat("x", 1024))
		time.Sleep(100 * time.Millisecond)
	}
}package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type RotatingLogger struct {
	currentFile   *os.File
	currentSize   int64
	maxSize       int64
	logDir        string
	baseName      string
	rotationCount int
}

func NewRotatingLogger(dir, name string, maxSize int64) (*RotatingLogger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		maxSize:  maxSize,
		logDir:   dir,
		baseName: name,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	path := filepath.Join(rl.logDir, rl.baseName+".log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(rl.logDir, rl.baseName+".log")
	newPath := filepath.Join(rl.logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.compressFile(newPath); err != nil {
		return err
	}

	rl.rotationCount++
	return rl.openCurrentFile()
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

	if err := os.Remove(source); err != nil {
		return err
	}

	return nil
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("./logs", "app", 1024*1024)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 10000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
}