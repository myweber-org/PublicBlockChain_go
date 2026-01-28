package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	baseName    string
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, name)
	logger := &RotatingLogger{baseName: basePath}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	path := rl.baseName + ".log"
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
	if rl.currentSize+int64(len(p)) > maxFileSize {
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

	timestamp := time.Now().Format("20060102-150405")
	oldPath := rl.baseName + ".log"
	newPath := fmt.Sprintf("%s-%s.log", rl.baseName, timestamp)

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.openCurrentFile(); err != nil {
		return err
	}

	go rl.cleanupOldFiles()
	return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
	pattern := rl.baseName + "-*.log"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, file := range toDelete {
			os.Remove(file)
		}
	}
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}
}