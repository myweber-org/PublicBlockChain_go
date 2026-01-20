package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	baseName    string
}

func NewLogRotator(name string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, name)
	lr := &LogRotator{baseName: basePath}

	if err := lr.openCurrent(); err != nil {
		return nil, err
	}

	return lr, nil
}

func (lr *LogRotator) openCurrent() error {
	path := lr.baseName + ".log"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.currentFile = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
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
	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := fmt.Sprintf("%s_%s.log", lr.baseName, timestamp)
	if err := os.Rename(lr.baseName+".log", oldPath); err != nil {
		return err
	}

	if err := lr.openCurrent(); err != nil {
		return err
	}

	lr.cleanupOld()
	return nil
}

func (lr *LogRotator) cleanupOld() {
	pattern := lr.baseName + "_*.log"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, path := range toDelete {
			os.Remove(path)
		}
	}
}

func (lr *LogRotator) Close() error {
	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app")
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
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

type LogRotator struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	backupCount int
}

func NewLogRotator(filePath string, maxSizeMB int, backupCount int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rotator := &LogRotator{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	stat, err := lr.currentFile.Stat()
	if err != nil {
		return 0, err
	}

	if stat.Size()+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	return lr.currentFile.Write(p)
}

func (lr *LogRotator) rotate() error {
	if err := lr.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.gz", lr.filePath, timestamp)

	if err := compressFile(lr.filePath, backupPath); err != nil {
		return err
	}

	if err := os.Remove(lr.filePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := lr.cleanOldBackups(); err != nil {
		return err
	}

	return lr.openCurrentFile()
}

func compressFile(source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func (lr *LogRotator) cleanOldBackups() error {
	pattern := fmt.Sprintf("%s.*.gz", lr.filePath)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= lr.backupCount {
		return nil
	}

	backups := make([]string, len(matches))
	copy(backups, matches)

	for i := 0; i < len(backups)-lr.backupCount; i++ {
		if err := os.Remove(backups[i]); err != nil {
			return err
		}
	}

	return nil
}

func (lr *LogRotator) openCurrentFile() error {
	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	lr.currentFile = file
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