
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LogRotator struct {
	filePath    string
	maxSize     int64
	backupCount int
}

func NewLogRotator(filePath string, maxSize int64, backupCount int) *LogRotator {
	return &LogRotator{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if err := lr.rotateIfNeeded(); err != nil {
		return 0, err
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write(p)
}

func (lr *LogRotator) rotateIfNeeded() error {
	info, err := os.Stat(lr.filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if info.Size() < lr.maxSize {
		return nil
	}

	for i := lr.backupCount - 1; i > 0; i-- {
		oldName := fmt.Sprintf("%s.%d", lr.filePath, i)
		newName := fmt.Sprintf("%s.%d", lr.filePath, i+1)
		if _, err := os.Stat(oldName); err == nil {
			os.Rename(oldName, newName)
		}
	}

	backupName := fmt.Sprintf("%s.1", lr.filePath)
	if err := os.Rename(lr.filePath, backupName); err != nil {
		return err
	}

	return nil
}

func (lr *LogRotator) CleanOldBackups() error {
	for i := lr.backupCount + 1; ; i++ {
		backupPath := fmt.Sprintf("%s.%d", lr.filePath, i)
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			break
		}
		if err := os.Remove(backupPath); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	rotator := NewLogRotator("app.log", 1024*1024, 5)

	for i := 0; i < 100; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		rotator.Write([]byte(logEntry))
	}

	rotator.CleanOldBackups()
	fmt.Println("Log rotation completed")
}