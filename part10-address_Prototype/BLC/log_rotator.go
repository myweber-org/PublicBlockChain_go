
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

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	backupCount int
}

func NewRotatingLogger(filePath string, maxSize int64, backupCount int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory failed: %w", err)
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file failed: %w", err)
	}

	rl.currentFile = file
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	info, err := rl.currentFile.Stat()
	if err != nil {
		return 0, fmt.Errorf("stat file failed: %w", err)
	}

	if info.Size()+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, fmt.Errorf("rotate failed: %w", err)
		}
	}

	return rl.currentFile.Write(p)
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return fmt.Errorf("close current file failed: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return fmt.Errorf("rename file failed: %w", err)
	}

	if err := rl.compressBackup(backupPath); err != nil {
		return fmt.Errorf("compress backup failed: %w", err)
	}

	if err := rl.cleanOldBackups(); err != nil {
		return fmt.Errorf("clean old backups failed: %w", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressBackup(backupPath string) error {
	source, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("open backup file failed: %w", err)
	}
	defer source.Close()

	compressedPath := backupPath + ".gz"
	target, err := os.Create(compressedPath)
	if err != nil {
		return fmt.Errorf("create compressed file failed: %w", err)
	}
	defer target.Close()

	gzWriter := gzip.NewWriter(target)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, source); err != nil {
		return fmt.Errorf("compress data failed: %w", err)
	}

	if err := os.Remove(backupPath); err != nil {
		return fmt.Errorf("remove original backup failed: %w", err)
	}

	return nil
}

func (rl *RotatingLogger) cleanOldBackups() error {
	pattern := rl.filePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob backup files failed: %w", err)
	}

	if len(matches) <= rl.backupCount {
		return nil
	}

	backupsToRemove := matches[:len(matches)-rl.backupCount]
	for _, backup := range backupsToRemove {
		if err := os.Remove(backup); err != nil {
			return fmt.Errorf("remove old backup %s failed: %w", backup, err)
		}
	}

	return nil
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
	logger, err := NewRotatingLogger("./logs/app.log", 1024*1024, 5)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write failed: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}