
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

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewLogRotator(name string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, name)
	file, err := os.OpenFile(basePath+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &LogRotator{
		currentFile: file,
		currentSize: info.Size(),
		baseName:    name,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

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
	if err := lr.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102-150405")
	oldPath := filepath.Join(logDir, lr.baseName+".log")
	newPath := filepath.Join(logDir, fmt.Sprintf("%s-%s.log", lr.baseName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := lr.compressFile(newPath); err != nil {
		return err
	}

	if err := lr.cleanupOldFiles(); err != nil {
		return err
	}

	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.currentFile = file
	lr.currentSize = 0
	return nil
}

func (lr *LogRotator) compressFile(path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path + ".gz")
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

func (lr *LogRotator) cleanupOldFiles() error {
	pattern := filepath.Join(logDir, lr.baseName+"-*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > maxBackups {
		filesToRemove := matches[:len(matches)-maxBackups]
		for _, file := range filesToRemove {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}
	return nil
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator("application")
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}