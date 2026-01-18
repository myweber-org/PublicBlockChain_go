
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

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	fileCount    int
	maxFiles     int
	compressOld  bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int, compressOld bool) (*RotatingLogger, error) {
	if maxSizeMB <= 0 {
		return nil, fmt.Errorf("maxSizeMB must be positive")
	}
	if maxFiles <= 0 {
		return nil, fmt.Errorf("maxFiles must be positive")
	}

	rl := &RotatingLogger{
		basePath:    basePath,
		maxSize:     int64(maxSizeMB) * 1024 * 1024,
		maxFiles:    maxFiles,
		compressOld: compressOld,
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

	file, err := os.OpenFile(rl.basePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
		rl.currentFile = nil
	}

	dir := filepath.Dir(rl.basePath)
	baseName := filepath.Base(rl.basePath)
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)

	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := filepath.Join(dir, fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext))

	if err := os.Rename(rl.basePath, rotatedPath); err != nil {
		return err
	}

	rl.fileCount++

	if rl.compressOld {
		go rl.compressFile(rotatedPath)
	}

	if rl.fileCount > rl.maxFiles {
		go rl.cleanupOldFiles(dir, nameWithoutExt, ext)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(path string) {
	compressedPath := path + ".gz"
	log.Printf("Compressing %s to %s", path, compressedPath)
}

func (rl *RotatingLogger) cleanupOldFiles(dir, baseName, ext string) {
	pattern := filepath.Join(dir, baseName+"_*"+ext)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("Error finding old log files: %v", err)
		return
	}

	if len(matches) > rl.maxFiles {
		filesToDelete := matches[:len(matches)-rl.maxFiles]
		for _, file := range filesToDelete {
			if err := os.Remove(file); err != nil {
				log.Printf("Error deleting old log file %s: %v", file, err)
			} else {
				log.Printf("Deleted old log file: %s", file)
			}
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

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logger))

	for i := 0; i < 100; i++ {
		log.Printf("Log entry number %d: %s", i, strings.Repeat("x", 1024))
		time.Sleep(100 * time.Millisecond)
	}
}