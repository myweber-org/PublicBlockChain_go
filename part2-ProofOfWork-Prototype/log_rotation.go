package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type RotatingLogger struct {
	filePath    string
	maxSize     int64
	maxAge      time.Duration
	currentFile *os.File
	logger      *log.Logger
	written     int64
}

func NewRotatingLogger(filePath string, maxSize int64, maxAge time.Duration) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath: filePath,
		maxSize:  maxSize,
		maxAge:   maxAge,
	}
	if err := rl.openFile(); err != nil {
		return nil, err
	}
	go rl.cleanupOldFiles()
	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.written = info.Size()
	rl.logger = log.New(file, "", log.LstdFlags)
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	if rl.written+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.currentFile.Write(p)
	if err == nil {
		rl.written += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	timestamp := time.Now().Format("20060102_150405")
	backupPath := rl.filePath + "." + timestamp

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}

	return rl.openFile()
}

func (rl *RotatingLogger) cleanupOldFiles() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		files, err := filepath.Glob(rl.filePath + ".*")
		if err != nil {
			continue
		}

		cutoff := time.Now().Add(-rl.maxAge)
		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				os.Remove(file)
			}
		}
	}
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}